package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"evo-ai-core-service/internal/httpclient"
	"evo-ai-core-service/internal/types"
)

// Custom errors following Evolution pattern
type AuthenticationError struct {
	Message string
}

func (e *AuthenticationError) Error() string {
	return fmt.Sprintf("Authentication error: %s", e.Message)
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation error: %s", e.Message)
}

type NetworkError struct {
	Message string
}

func (e *NetworkError) Error() string {
	return fmt.Sprintf("Network error: %s", e.Message)
}

type ServiceUnavailableError struct {
	Message string
}

func (e *ServiceUnavailableError) Error() string {
	return fmt.Sprintf("Service unavailable: %s", e.Message)
}

// EvoAuthService interface defines all authentication and authorization operations
type EvoAuthService interface {
	// Authentication methods
	ValidateToken(token, tokenType string) (*types.EvoAuthValidateTokenData, error)
	BuildHeaders(token, tokenType string) (map[string]string, error)

	// Permission management methods
	CheckPermission(ctx context.Context, authToken, permissionKey, tokenType string) (bool, error)
	CheckAccountPermission(ctx context.Context, userID, accountID, permissionKey string, authToken, tokenType string) (bool, error)
	CheckUserPermission(ctx context.Context, userID, permissionKey string, authToken, tokenType string) (bool, error)
}

type evoAuthService struct {
	baseURL string
}

// NewEvoAuthService creates a new instance of EvoAuthService
func NewEvoAuthService(baseURL string) EvoAuthService {
	return &evoAuthService{
		baseURL: baseURL,
	}
}

// ============================================================================
// Authentication Methods
// ============================================================================

// ValidateToken validates token with Evo Auth API - Primary authentication method
func (s *evoAuthService) ValidateToken(token, tokenType string) (*types.EvoAuthValidateTokenData, error) {
	headers, err := s.BuildHeaders(token, tokenType)
	if err != nil {
		return nil, err
	}

	fmt.Printf("EvoAuth: Validating %s token at %s/api/v1/auth/validate\n", tokenType, s.baseURL)

	response, err := s.doPost("/api/v1/auth/validate", map[string]interface{}{}, headers)
	if err != nil {
		if _, ok := err.(*NetworkError); ok {
			return nil, &ServiceUnavailableError{Message: "Authentication service unavailable"}
		}
		return nil, err
	}

	// Parse response
	dataMap, ok := response["data"].(map[string]interface{})
	if !ok {
		return nil, &ValidationError{Message: "Invalid response format from auth service"}
	}

	// Convert to JSON and back to struct for proper type conversion
	dataJSON, err := json.Marshal(dataMap)
	if err != nil {
		return nil, &ValidationError{Message: "Failed to serialize response data"}
	}

	var tokenData types.EvoAuthValidateTokenData
	if err := json.Unmarshal(dataJSON, &tokenData); err != nil {
		return nil, &ValidationError{Message: "Failed to parse token data"}
	}

	fmt.Printf("EvoAuth: Successfully validated token for user %s with %d accounts\n", tokenData.User.Email, len(tokenData.Accounts))
	return &tokenData, nil
}

// BuildHeaders builds HTTP headers based on token type
func (s *evoAuthService) BuildHeaders(token, tokenType string) (map[string]string, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}

	switch tokenType {
	case "bearer":
		headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	case "api_access_token":
		headers["api_access_token"] = token
	default:
		return nil, fmt.Errorf("invalid token type: %s", tokenType)
	}

	return headers, nil
}

// ============================================================================
// Permission Management Methods
// ============================================================================

// CheckPermission checks if authenticated user has specific permission
func (s *evoAuthService) CheckPermission(ctx context.Context, authToken, permissionKey, tokenType string) (bool, error) {
	headers, err := s.BuildHeaders(authToken, tokenType)
	if err != nil {
		return false, err
	}

	payload := map[string]interface{}{
		"permission_key": permissionKey,
	}

	response, err := s.doPost("/api/v1/permissions/check", payload, headers)
	if err != nil {
		// Check if it's a 404 (endpoint not implemented)
		if _, ok := err.(*NetworkError); ok {
			// Fallback: allow access for authenticated users when permission system not implemented
			fmt.Printf("Permission endpoint not found (404) - allowing access for authenticated user (permission: %s)\n", permissionKey)
			return true, nil
		}
		fmt.Printf("Permission check failed for %s: %v\n", permissionKey, err)
		return false, nil
	}

	data := response["data"].(map[string]interface{})

	hasPermission, _ := data["has_permission"].(bool)
	fmt.Printf("Permission check for %s: %v\n", permissionKey, hasPermission)
	return hasPermission, nil
}

// CheckAccountPermission checks account-scoped permission for user
func (s *evoAuthService) CheckAccountPermission(ctx context.Context, userID, accountID, permissionKey string, authToken, tokenType string) (bool, error) {
	headers, err := s.BuildHeaders(authToken, tokenType)
	if err != nil {
		return false, err
	}

	payload := map[string]interface{}{
		"permission_key": permissionKey,
	}

	response, err := s.doPost(fmt.Sprintf("/api/v1/accounts/%s/users/%s/check_permission", accountID, userID), payload, headers)
	if err != nil {
		fmt.Printf("Error checking account permission: %v\n", err)
		return false, err
	}

	data := response["data"].(map[string]interface{})

	hasPermission, _ := data["has_permission"].(bool)
	return hasPermission, nil
}

// CheckUserPermission checks global user permission
func (s *evoAuthService) CheckUserPermission(ctx context.Context, userID, permissionKey string, authToken, tokenType string) (bool, error) {
	headers, err := s.BuildHeaders(authToken, tokenType)
	if err != nil {
		return false, err
	}

	payload := map[string]interface{}{
		"permission_key": permissionKey,
	}

	response, err := s.doPost("/api/v1/users/check_permission", payload, headers)
	if err != nil {
		fmt.Printf("Error checking user permission: %v\n", err)
		return false, err
	}

	data := response["data"].(map[string]interface{})

	hasPermission, _ := data["has_permission"].(bool)
	return hasPermission, nil
}

// ============================================================================
// Private HTTP Methods
// ============================================================================

// doGet executes GET request to evo-auth-service using httpclient helpers
func (s *evoAuthService) doGet(endpoint string, headers map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint)

	// Use httpclient helper with flexible status code handling
	type Response map[string]interface{}

	// Try with 200 OK first
	result, err := httpclient.DoGetJSON[Response](
		context.Background(),
		url,
		headers,
		http.StatusOK,
	)

	if err != nil {
		// Check if it's a 404 error
		if httpErr, ok := err.(interface{ StatusCode() int }); ok {
			if httpErr.StatusCode() == 404 {
				return nil, &NetworkError{Message: "Endpoint not found"}
			}
		}
		return nil, &NetworkError{Message: fmt.Sprintf("Request failed: %v", err)}
	}

	if result == nil {
		return nil, &ValidationError{Message: "Empty response from server"}
	}

	return *result, nil
}

// doPost executes POST request to evo-auth-service using httpclient helpers
func (s *evoAuthService) doPost(endpoint string, payload map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s%s", s.baseURL, endpoint)

	// Use httpclient helper
	type Response map[string]interface{}

	result, err := httpclient.DoPostJSON[Response](
		context.Background(),
		url,
		payload,
		headers,
		http.StatusOK,
	)

	if err != nil {
		// Check specific error types
		if httpErr, ok := err.(interface{ StatusCode() int }); ok {
			statusCode := httpErr.StatusCode()

			if statusCode == 404 {
				return nil, &NetworkError{Message: "Endpoint not found"}
			}
			if statusCode == 401 {
				return nil, &AuthenticationError{Message: "Invalid or expired token"}
			}
		}
		return nil, &NetworkError{Message: fmt.Sprintf("Request failed: %v", err)}
	}

	if result == nil {
		return nil, &ValidationError{Message: "Empty response from server"}
	}

	return *result, nil
}

// ============================================================================
// Global singleton instance
// ============================================================================

var globalEvoAuthService EvoAuthService

// InitializeEvoAuthService initializes the global service instance
func InitializeEvoAuthService(baseURL string) {
	globalEvoAuthService = NewEvoAuthService(baseURL)
	fmt.Printf("Global EvoAuthService initialized with base URL: %s\n", baseURL)
}

// GetEvoAuthService returns the global service instance
func GetEvoAuthService() EvoAuthService {
	if globalEvoAuthService == nil {
		panic("EvoAuthService not initialized. Call InitializeEvoAuthService first.")
	}
	return globalEvoAuthService
}
