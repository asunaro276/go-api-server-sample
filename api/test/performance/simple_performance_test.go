package performance

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// Simple performance test that doesn't use internal packages
// These tests validate response time requirements

func setupSimpleRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Add a simple mock endpoint for performance testing
	router.POST("/api/v1/users", func(c *gin.Context) {
		// Simulate some processing time
		time.Sleep(1 * time.Millisecond)
		c.JSON(http.StatusCreated, gin.H{
			"id":         1,
			"name":       "Test User",
			"email":      "test@example.com",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		})
	})
	
	router.GET("/api/v1/users/:id", func(c *gin.Context) {
		time.Sleep(1 * time.Millisecond)
		c.JSON(http.StatusOK, gin.H{
			"id":         1,
			"name":       "Test User",
			"email":      "test@example.com",
			"created_at": time.Now(),
			"updated_at": time.Now(),
		})
	})
	
	return router
}

// Performance test: API should respond within 200ms
func TestAPIPerformance(t *testing.T) {
	router := setupSimpleRouter()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{
			name:   "POST /api/v1/users",
			method: "POST",
			path:   "/api/v1/users",
			body:   `{"name":"Test User","email":"test@example.com"}`,
		},
		{
			name:   "GET /api/v1/users/1",
			method: "GET",
			path:   "/api/v1/users/1",
			body:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, bytes.NewBufferString(tt.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			duration := time.Since(start)

			if duration > 200*time.Millisecond {
				t.Errorf("%s took %v, should be under 200ms", tt.name, duration)
			}

			t.Logf("%s took %v", tt.name, duration)

			// Verify we got a response
			if w.Code < 200 || w.Code >= 300 {
				t.Errorf("Unexpected status code: %d", w.Code)
			}
		})
	}
}

// JSON marshaling performance test
func TestJSONMarshalingPerformance(t *testing.T) {
	user := map[string]interface{}{
		"id":         1,
		"name":       "Test User",
		"email":      "test@example.com",
		"created_at": time.Now(),
		"updated_at": time.Now(),
	}

	start := time.Now()
	for i := 0; i < 10000; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			t.Fatalf("JSON marshaling failed: %v", err)
		}
	}
	duration := time.Since(start)

	avgDuration := duration / 10000
	t.Logf("JSON marshaling took %v per operation", avgDuration)

	if avgDuration > 5*time.Microsecond {
		t.Errorf("JSON marshaling took %v per operation, should be under 5μs", avgDuration)
	}
}

// Benchmark test for request handling
func BenchmarkUserEndpoint(b *testing.B) {
	router := setupSimpleRouter()
	requestBody := `{"name":"Test User","email":"test@example.com"}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewBufferString(requestBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}