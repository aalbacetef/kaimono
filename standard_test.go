package kaimono

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStandardGet(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mock := newMockBackend()

	svc, err := NewService(mock, mock, mock, logger)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	// add an empty cart to the first session
	mock.carts = append(mock.carts, mkEmptyTestCart())
	mock.data[mock.sessions[0]] = len(mock.carts) - 1

	tests := []struct {
		label        string
		sessionToken string
		wantCode     int
	}{
		{
			label:        "should return 200 on cart existing for session",
			sessionToken: mock.sessions[0],
			wantCode:     http.StatusOK,
		},
		{
			label:        "should return 400 for non-existent sessions",
			sessionToken: "non-existent-session",
			wantCode:     http.StatusBadRequest,
		},
		{
			label:        "should return 404 for non-existent carts",
			sessionToken: mock.sessions[1],
			wantCode:     http.StatusNotFound,
		},
		{
			label:        "should return 500 on the rest of the errors",
			sessionToken: "", // our test func won't add a session on empty tokens
			wantCode:     http.StatusInternalServerError,
		},
	}

	for _, c := range tests {
		t.Run(c.label, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, "", nil)
			if err != nil {
				t.Fatalf("could not make request: %v", err)
			}

			if c.sessionToken != "" {
				setTestCookie(req, c.sessionToken)
			}

			w := httptest.NewRecorder()
			svc.Get(w, req)

			result := w.Result()
			defer result.Body.Close()

			if result.StatusCode != c.wantCode {
				t.Fatalf("got code %d, want %d", result.StatusCode, c.wantCode)
			}
		})
	}
}

func setTestCookie(req *http.Request, v string) {
	cookie := http.Cookie{
		Name:     testCookieName,
		Value:    v,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	req.AddCookie(&cookie)
}
