package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ulshv/online-store-app/backend-go/internal/modules/user"
)

func initModule() *AuthModule {
	userModule := user.NewUserModule()
	return NewAuthModule(userModule)
}

func TestLogin(t *testing.T) {
	module := initModule()
	ts := httptest.NewServer(http.HandlerFunc(module.AuthController.loginHandler))
	defer ts.Close()

	// Test cases
	tests := []struct {
		name       string
		payload    loginDto
		wantStatus int
	}{
		{
			name: "valid login",
			payload: loginDto{
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "invalid credentials",
			payload: loginDto{
				Email:    "wrong@example.com",
				Password: "wrongpass",
			},
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(tt.payload)
			if err != nil {
				t.Fatal(err)
			}

			fmt.Println("ts.URL:", ts.URL)
			resp, err := http.Post(ts.URL, "application/json", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				body, _ := io.ReadAll(resp.Body)
				fmt.Println("resp.Body:", string(body))
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
