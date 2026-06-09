package handler

import (
	"evo-ai-core-service/internal/httpclient/errors"
	"evo-ai-core-service/internal/httpclient/response"
	"evo-ai-core-service/internal/middleware"
	folderService "evo-ai-core-service/pkg/folder/service"
	"evo-ai-core-service/pkg/folder_share/model"
	"evo-ai-core-service/pkg/folder_share/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FolderShareHandler interface defines the contract for folder share handlers
type FolderShareHandler interface {
	RegisterRoutesMiddleware(router gin.IRouter)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetSharedFolder(c *gin.Context)
	GetSharedFolders(c *gin.Context)
	ListSharedFolders(c *gin.Context)
}

// folderShareHandler implements the FolderShareHandler interface
type folderShareHandler struct {
	folderService      folderService.FolderService
	folderShareService service.FolderShareService
}

// NewFolderShareHandler creates a new folder share handler
func NewFolderShareHandler(folderShareService service.FolderShareService, folderService folderService.FolderService) FolderShareHandler {
	return &folderShareHandler{
		folderService:      folderService,
		folderShareService: folderShareService,
	}
}

// RegisterRoutesMiddleware registers the routes for the folder share handler with middleware
func (h *folderShareHandler) RegisterRoutesMiddleware(router gin.IRouter) {
	// Usar middleware global de permissão
	permissionMiddleware := middleware.GetGlobalPermissionMiddleware()

	foldersShare := router.Group("/agents/folders")
	{
		// Share permissions
		foldersShare.POST("/:id/share",
			permissionMiddleware.RequirePermission("ai_folders", "share"),
			h.Create)

		// Access shared permissions
		foldersShare.GET("/:id/shared",
			permissionMiddleware.RequirePermission("ai_folders", "access_shared"),
			h.GetSharedFolder)
		foldersShare.GET("/shared",
			permissionMiddleware.RequirePermission("ai_folders", "access_shared"),
			h.GetSharedFolders)
		foldersShare.GET("/accessible",
			permissionMiddleware.RequirePermission("ai_folders", "access_shared"),
			h.ListSharedFolders)

		// Update shared permissions
		foldersShare.PUT("/shared/:id",
			permissionMiddleware.RequirePermission("ai_folders", "share"),
			h.Update)

		// Delete shared permissions
		foldersShare.DELETE("/shared/:id",
			permissionMiddleware.RequirePermission("ai_folders", "share"),
			h.Delete)
	}
}

// ShareFolder handles the share folder request
func (h *folderShareHandler) Create(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	var req model.FolderShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	folderShare := model.FolderShare{
		FolderID:        id,
		SharedWithEmail: req.SharedWithEmail,
		PermissionLevel: req.PermissionLevel,
	}

	sharedFolder, err := h.folderShareService.Create(c.Request.Context(), folderShare)

	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, sharedFolder.ToResponse(), "Folder shared successfully", http.StatusOK)
}

// Update handles the update shared folder request
func (h *folderShareHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	var req model.FolderShareUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationErrorResponse(c, err)
		return
	}

	folderShare := &model.FolderShare{
		PermissionLevel: req.PermissionLevel,
		IsActive:        *req.IsActive,
	}

	updatedFolderShare, err := h.folderShareService.Update(c.Request.Context(), id, folderShare)

	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, updatedFolderShare.ToResponse(), "Folder share updated successfully", http.StatusOK)
}

// Delete handles the delete folder share request
func (h *folderShareHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	_, err = h.folderShareService.Delete(c.Request.Context(), id)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.SuccessResponse(c, nil, "Folder share deleted successfully", http.StatusNoContent)
}

// GetSharedFolder handles the get folder shared request
func (h *folderShareHandler) GetSharedFolder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	folderShares, err := h.folderShareService.GetSharedFolder(c.Request.Context(), id, page, pageSize)

	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.PaginatedResponse(c, folderShares.Items, folderShares.Page, folderShares.PageSize, int(folderShares.TotalItems), "Folder shares retrieved successfully", http.StatusOK)
}

// GetSharedFolders handles the get shared folders request
func (h *folderShareHandler) GetSharedFolders(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	folderShares, err := h.folderShareService.GetSharedFoldersWithEmail(c.Request.Context(), page, pageSize)

	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.PaginatedResponse(c, folderShares.Items, folderShares.Page, folderShares.PageSize, int(folderShares.TotalItems), "Shared folders retrieved successfully", http.StatusOK)
	}

// ListSharedFolders handles the list shared folders request
func (h *folderShareHandler) ListSharedFolders(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	if skip := c.Query("skip"); skip != "" {
		skipVal, err := strconv.Atoi(skip)
		if err != nil {
			code, message, httpCode := errors.HandleError(err)
			response.ErrorResponse(c, code, message, nil, httpCode)
			return
		}

		limitVal, err := strconv.Atoi(c.DefaultQuery("limit", "100"))
		if err != nil {
			code, message, httpCode := errors.HandleError(err)
			response.ErrorResponse(c, code, message, nil, httpCode)
			return
		}

		page = (skipVal / limitVal) + 1
		pageSize = limitVal
	}

	sharedFolders, err := h.folderShareService.ListSharedFolders(c.Request.Context(), page, pageSize)
	if err != nil {
		code, message, httpCode := errors.HandleError(err)
		response.ErrorResponse(c, code, message, nil, httpCode)
		return
	}

	response.PaginatedResponse(c, sharedFolders.Items, sharedFolders.Page, sharedFolders.PageSize, int(sharedFolders.TotalItems), "Accessible folders retrieved successfully", http.StatusOK)
}
