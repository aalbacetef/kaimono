package kaimono

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func setCookie(req *http.Request, v string) {
	cookie := http.Cookie{
		Name:     testCookieName,
		Value:    v,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	req.AddCookie(&cookie)
}

func TestStandardGet(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mock := newMockBackend()

	svc, err := NewService(mock, mock, mock, logger)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	tests := []struct {
		label        string
		sessionToken string
		wantCode     int
	}{
		{
			label:        "should return 400 for non-existent sessions",
			sessionToken: "non-existent-session",
			wantCode:     400,
		},
	}

	for _, c := range tests {
		t.Run(c.label, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "", nil)
			if err != nil {
				t.Fatalf("could not make request: %v", err)
			}

			setCookie(req, c.sessionToken)

			w := httptest.NewRecorder()
			svc.Get(w, req)

			result := w.Result()
			defer result.Body.Close()

			m := make(map[string]any)

			if err := json.NewDecoder(result.Body).Decode(&m); err != nil {
				t.Fatalf("error decoding body: %v", err)
			}

			t.Logf("body: %+v", m)

			if result.StatusCode != c.wantCode {
				t.Fatalf("got code %d, want %d", result.StatusCode, c.wantCode)
			}
		})
	}
}
