package controller

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

type lookup func(context.Context, sql.NullString) ([]sql.NullString, error)

func (f lookup) GetMediaByTag(ctx context.Context, tag sql.NullString) ([]sql.NullString, error) {
	return f(ctx, tag)
}

func TestMediaEndpoint(t *testing.T) {
	for _, tc := range []struct {
		query, tag string
		code       int
	}{
		{"", "", 400},
		{"?type=bad&id=1", "", 400},
		{"?id=123", "imdb://tt123", 200},
		{"?type=imdb&id=tt123", "imdb://tt123", 200},
		{"?type=tmdb&id=42", "tmdb://42", 200},
		{"?type=tvdb&id=7", "tvdb://7", 200},
	} {
		t.Run(tc.query, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()
			called := false
			e := echo.New()
			e.GET("/media", New(lookup(func(got context.Context, tag sql.NullString) ([]sql.NullString, error) {
				called = true
				if got.Err() != context.Canceled || !tag.Valid || tag.String != tc.tag {
					t.Fatalf("context %v, tag %v", got.Err(), tag)
				}
				return nil, nil
			})).GetMediaByTag)
			w := httptest.NewRecorder()
			e.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/media"+tc.query, nil).WithContext(ctx))
			if w.Code != tc.code || called != (tc.tag != "") {
				t.Fatalf("code %d called %v", w.Code, called)
			}
			if w.Code == 200 && strings.TrimSpace(w.Body.String()) != "[]" {
				t.Fatalf("unexpected body: %q", w.Body.String())
			}
		})
	}
}
