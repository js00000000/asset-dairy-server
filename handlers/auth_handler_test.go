package handlers

import (
	"asset-dairy/models"
	"asset-dairy/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthService is a mock implementation of AuthServiceInterface
type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) SignUp(req *models.SignUpRequest) (*models.User, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockAuthService) SignIn(email, password string) (*models.AuthResponse, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AuthResponse), args.Error(1)
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.String(1), args.Error(1)
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	// Set up test environment
	gin.SetMode(gin.TestMode)
	os.Setenv("JWT_REFRESH_SECRET", "test-refresh-secret")

	tests := []struct {
		name           string
		setupMock      func(*MockAuthService)
		cookieValue    string
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful token refresh",
			setupMock: func(m *MockAuthService) {
				// Mock successful access token generation
				m.On("RefreshToken", mock.Anything).Return("new-access-token", nil)
				// Mock successful refresh token generation
				m.On("GenerateNewRefreshToken", "user123", "test@example.com").Return("new-refresh-token", nil)
			},
			cookieValue:    generateTestRefreshToken(t, "user123", "test@example.com"),
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"token": "new-access-token",
			},
		},
		{
			name:           "missing refresh token cookie",
			setupMock:      func(m *MockAuthService) {},
			cookieValue:    "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Refresh token missing",
			},
		},
		{
			name: "invalid refresh token",
			setupMock: func(m *MockAuthService) {
				m.On("RefreshToken", "invalid-token").Return("", services.ErrInvalidToken)
			},
			cookieValue:    "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid refresh token",
			},
		},
		{
			name: "failed to generate new access token",
			setupMock: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything).Return("", assert.AnError)
			},
			cookieValue:    generateTestRefreshToken(t, "user123", "test@example.com"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to generate new access token",
			},
		},
		{
			name: "failed to generate new refresh token",
			setupMock: func(m *MockAuthService) {
				m.On("RefreshToken", mock.Anything).Return("new-access-token", nil)
				m.On("GenerateNewRefreshToken", "user123", "test@example.com").Return("", assert.AnError)
			},
			cookieValue:    generateTestRefreshToken(t, "user123", "test@example.com"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to generate new refresh token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock service
			mockService := new(MockAuthService)
			tt.setupMock(mockService)

			// Create a new handler with the mock service
			handler := NewAuthHandler(mockService)

			// Create a new test router
			router := gin.New()
			router.POST("/refresh", handler.RefreshToken)

			// Create a new test request
			req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
			if tt.cookieValue != "" {
				req.AddCookie(&http.Cookie{
					Name:     "refresh_token",
					Value:    tt.cookieValue,
					Path:     "/",
					HttpOnly: true,
				})
			}

			// Create a new test response recorder
			w := httptest.NewRecorder()

			// Serve the request
			router.ServeHTTP(w, req)

			// Assert the response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			// Verify that the mock was called as expected
			mockService.AssertExpectations(t)
		})
	}
}

// generateTestRefreshToken creates a valid refresh token for testing
func generateTestRefreshToken(t *testing.T, userID, email string) string {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("test-refresh-secret"))
	assert.NoError(t, err)
	return tokenString
}
