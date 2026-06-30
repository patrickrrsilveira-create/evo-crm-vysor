package main

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
)

type EvoAuthValidateTokenData struct {
	User     EvoAuthUser      `json:"user"`
	Accounts []EvoAuthAccount `json:"accounts"`
}

type EvoAuthUser struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	DisplayName  *string   `json:"display_name"`
	Availability interface{}    `json:"availability"`
	MFAEnabled   bool      `json:"mfa_enabled"`
	Confirmed    bool      `json:"confirmed"`
	Type         string    `json:"type"`
	Role         *Role     `json:"role"`
}

type Role struct {
	ID   uuid.UUID `json:"id"`
	Key  string    `json:"key"`
	Name string    `json:"name"`
}

type EvoAuthFeature struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type EvoAuthPlan struct {
	ID       uuid.UUID        `json:"id"`
	PlanName string           `json:"name"`
	IsActive bool             `json:"is_active"`
	IsCustom bool             `json:"is_custom"`
	StartsAt string           `json:"starts_at"`
	EndsAt   string           `json:"ends_at"`
	Features []EvoAuthFeature `json:"features"`
}

type EvoAuthAccount struct {
	ID         uuid.UUID        `json:"id"`
	Name       string           `json:"name"`
	Status     string           `json:"status"`
	Locale     string           `json:"locale"`
	CreatedAt  string           `json:"created_at"`
	UpdatedAt  string           `json:"updated_at"`
	Features   json.RawMessage  `json:"features"`
	ActivePlan *EvoAuthPlan     `json:"active_plan,omitempty"`
}

func main() {
	jsonData := `{
		"user": {
		  "id": "ac7f0c13-33e1-4bd5-8eb3-228cf0a1c1d8",
		  "name": "Admin",
		  "email": "contato@vysortech.com.br",
		  "type": "SuperAdmin",
		  "role": {
			"id": "123e4567-e89b-12d3-a456-426614174000",
			"key": "admin",
			"name": "Admin"
		  },
		  "pubsub_token": "token123",
		  "avatar_url": null,
		  "created_at": "2026-06-29T20:00:00.000Z",
		  "updated_at": "2026-06-29T20:00:00.000Z",
		  "ui_settings": {},
		  "mfa_enabled": false,
		  "mfa_setup_incomplete": false,
		  "display_name": "Admin User",
		  "available_name": "Admin",
		  "availability": 0,
		  "confirmed": true,
		  "confirmed_at": "2026-06-29T20:00:00.000Z",
		  "custom_attributes": {},
		  "setup_survey_completed": false
		},
		"accounts": [
		  {
			"id": "b1b8b8b8-b8b8-b8b8-b8b8-b8b8b8b8b8b8",
			"name": "Evolution Community",
			"domain": "localhost",
			"support_email": "support@evolution.com",
			"locale": "en",
			"status": "active",
			"features": "",
			"settings": {},
			"custom_attributes": {},
			"role": {
			  "id": "123e4567-e89b-12d3-a456-426614174000",
			  "key": "admin",
			  "name": "Admin"
			}
		  }
		]
	  }`

	var tokenData EvoAuthValidateTokenData
	err := json.Unmarshal([]byte(jsonData), &tokenData)
	if err != nil {
		fmt.Printf("Unmarshal Error: %v\n", err)
	} else {
		fmt.Println("Success!")
	}
}
