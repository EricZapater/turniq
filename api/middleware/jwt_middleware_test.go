package middleware

import (
	"api/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestJWTMiddleware(t *testing.T) {
	// Setup Config
	cfg := config.Config{}
	cfg.Auth.Secret = "test_secret"
	cfg.Auth.TTL = time.Hour

	// Setup Middleware
	mw, err := SetupJWT(cfg)
	if err != nil {
		t.Fatalf("Failed to setup JWT middleware: %v", err)
	}

	// Create a user
	user := &AuthUser{
		ID:         "user123",
		Username:   "testuser",
		Email:      "test@example.com",
		CustomerID: "cust456",
		IsAdmin:    true,
	}

	// Generate Token
	tokenString, _, err := mw.TokenGenerator(user)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Setup Gin Router
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(func(c *gin.Context) {
		// Mock Authorization header
		c.Request.Header.Set("Authorization", "Bearer "+tokenString)
		c.Next()
	})
	r.Use(mw.MiddlewareFunc())
	
	// Test Endpoint
	r.GET("/test", func(c *gin.Context) {
		// Verify Helpers
		gotUser := GetUser(c)
		if gotUser == nil {
			t.Error("GetUser returned nil")
			return
		}
		if gotUser.ID != user.ID {
			t.Errorf("Expected ID %s, got %s", user.ID, gotUser.ID)
		}
		if gotUser.Username != user.Username {
			t.Errorf("Expected Username %s, got %s", user.Username, gotUser.Username)
		}
		if gotUser.Email != user.Email {
			t.Errorf("Expected Email %s, got %s", user.Email, gotUser.Email)
		}
		if gotUser.CustomerID != user.CustomerID {
			t.Errorf("Expected CustomerID %s, got %s", user.CustomerID, gotUser.CustomerID)
		}
		if gotUser.IsAdmin != user.IsAdmin {
			t.Errorf("Expected IsAdmin %v, got %v", user.IsAdmin, gotUser.IsAdmin)
		}

		if GetUserID(c) != user.ID {
			t.Errorf("GetUserID failed")
		}
		if GetCustomerID(c) != user.CustomerID {
			t.Errorf("GetCustomerID failed")
		}
		if IsAdmin(c) != user.IsAdmin {
			t.Errorf("IsAdmin failed")
		}

		c.Status(http.StatusOK)
	})

	// Perform Request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
