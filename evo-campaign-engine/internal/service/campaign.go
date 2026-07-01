package service

import (
	"context"
	"errors"
	"evo-campaign-engine/internal/domain"
	"evo-campaign-engine/internal/repository"
	"time"

	"github.com/google/uuid"
)

type CampaignService struct {
	campaigns    *repository.CampaignRepo
	audiences    *repository.AudienceRepo
	jobs         *repository.SendJobRepo
	senders      *repository.SenderRepo
	suppressions *repository.SuppressionRepo
	throttles    *repository.ThrottleRepo
}

func NewCampaignService(
	campaigns *repository.CampaignRepo,
	audiences *repository.AudienceRepo,
	jobs *repository.SendJobRepo,
	senders *repository.SenderRepo,
	suppressions *repository.SuppressionRepo,
	throttles *repository.ThrottleRepo,
) *CampaignService {
	return &CampaignService{
		campaigns:    campaigns,
		audiences:    audiences,
		jobs:         jobs,
		senders:      senders,
		suppressions: suppressions,
		throttles:    throttles,
	}
}

type CreateCampaignRequest struct {
	AccountID         int        `json:"account_id" binding:"required"`
	Name              string     `json:"name" binding:"required"`
	ThrottleProfileID *uuid.UUID `json:"throttle_profile_id"`
	InboxIDs          []int      `json:"inbox_ids" binding:"required,min=1"`

	Audience AudienceInput  `json:"audience" binding:"required"`
	Variants []VariantInput `json:"variants" binding:"required,min=1"`
}

type AudienceInput struct {
	Mode          string                 `json:"mode"`
	SegmentFilter map[string]interface{} `json:"segment_filter"`
	Contacts      []ContactInput         `json:"contacts"`
}

type ContactInput struct {
	ContactID int    `json:"contact_id"`
	Recipient string `json:"recipient"`
	Timezone  string `json:"timezone"`
}

type VariantInput struct {
	Body      string `json:"body" binding:"required"`
	MediaURL  string `json:"media_url"`
	MediaType string `json:"media_type"`
	Weight    int    `json:"weight"`
}

func (s *CampaignService) Create(ctx context.Context, req CreateCampaignRequest) (*domain.Campaign, error) {
	campaign := &domain.Campaign{
		AccountID:         req.AccountID,
		Name:              req.Name,
		TriggerType:       "manual",
		ThrottleProfileID: req.ThrottleProfileID,
		Status:            domain.CampaignDraft,
		TotalRecipients:   len(req.Audience.Contacts),
	}

	if err := s.campaigns.Create(ctx, campaign); err != nil {
		return nil, err
	}

	channels := make([]domain.CampaignChannel, len(req.InboxIDs))
	for i, inboxID := range req.InboxIDs {
		channels[i] = domain.CampaignChannel{
			CampaignID: campaign.ID,
			InboxID:    inboxID,
		}
	}
	if err := s.campaigns.SaveChannels(ctx, channels); err != nil {
		return nil, err
	}

	variants := make([]domain.MessageVariant, len(req.Variants))
	for i, vi := range req.Variants {
		weight := vi.Weight
		if weight <= 0 {
			weight = 1
		}
		variants[i] = domain.MessageVariant{
			CampaignID: campaign.ID,
			Body:       vi.Body,
			MediaURL:   vi.MediaURL,
			MediaType:  vi.MediaType,
			Weight:     weight,
		}
	}
	if err := s.campaigns.SaveVariants(ctx, variants); err != nil {
		return nil, err
	}

	members := make([]domain.AudienceMember, len(req.Audience.Contacts))
	for i, c := range req.Audience.Contacts {
		members[i] = domain.AudienceMember{
			CampaignID: campaign.ID,
			ContactID:  c.ContactID,
			Recipient:  c.Recipient,
			Timezone:   c.Timezone,
			State:      domain.MemberPending,
		}
	}
	if len(members) > 0 {
		if err := s.audiences.SaveMembers(ctx, members); err != nil {
			return nil, err
		}
	}

	return s.campaigns.GetByID(ctx, campaign.ID)
}

func (s *CampaignService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Campaign, error) {
	return s.campaigns.GetByID(ctx, id)
}

func (s *CampaignService) List(ctx context.Context, accountID, page, pageSize int) ([]domain.Campaign, int64, error) {
	return s.campaigns.ListByAccount(ctx, accountID, page, pageSize)
}

func (s *CampaignService) Start(ctx context.Context, id uuid.UUID) error {
	campaign, err := s.campaigns.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if campaign.Status != domain.CampaignDraft && campaign.Status != domain.CampaignPaused {
		return errors.New("campaign must be in draft or paused state to start")
	}

	members, err := s.audiences.ListMembers(ctx, id)
	if err != nil {
		return err
	}
	if len(members) == 0 {
		return errors.New("campaign has no audience members")
	}

	campaignVariants, _ := s.campaigns.GetVariants(ctx, id)

	jobs := make([]domain.SendJob, 0, len(members))
	for _, m := range members {
		if m.State == domain.MemberDone || m.State == domain.MemberSkipped {
			continue
		}

		var variantID *uuid.UUID
		if len(campaignVariants) > 0 {
			vid := campaignVariants[0].ID
			variantID = &vid
		}

		jobs = append(jobs, domain.SendJob{
			CampaignID: id,
			ContactID:  m.ContactID,
			Recipient:  m.Recipient,
			VariantID:  variantID,
			State:      domain.JobQueued,
		})
	}

	if len(jobs) > 0 {
		if err := s.jobs.BulkCreate(ctx, jobs); err != nil {
			return err
		}
	}

	now := time.Now().UTC()
	campaign.Status = domain.CampaignRunning
	campaign.StartedAt = &now
	campaign.TotalRecipients = len(jobs)
	return s.campaigns.Update(ctx, campaign)
}

func (s *CampaignService) Pause(ctx context.Context, id uuid.UUID) error {
	campaign, err := s.campaigns.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if campaign.Status != domain.CampaignRunning {
		return errors.New("campaign must be running to pause")
	}

	campaign.Status = domain.CampaignPaused
	return s.campaigns.Update(ctx, campaign)
}

func (s *CampaignService) Cancel(ctx context.Context, id uuid.UUID) error {
	campaign, err := s.campaigns.GetByID(ctx, id)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	campaign.Status = domain.CampaignCancelled
	campaign.FinishedAt = &now
	return s.campaigns.Update(ctx, campaign)
}

func (s *CampaignService) Stats(ctx context.Context, id uuid.UUID) (map[domain.JobState]int, error) {
	return s.jobs.CountByState(ctx, id)
}

func (s *CampaignService) Delete(ctx context.Context, id uuid.UUID) error {
	campaign, err := s.campaigns.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if campaign.Status == domain.CampaignRunning {
		return errors.New("cannot delete a running campaign")
	}
	return s.campaigns.Delete(ctx, id)
}
