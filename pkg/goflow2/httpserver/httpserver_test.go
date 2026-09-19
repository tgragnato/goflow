package httpserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type failingResponseWriter struct{}

func (failingResponseWriter) Header() http.Header {
	return make(http.Header)
}

func (failingResponseWriter) Write([]byte) (int, error) {
	return 0, http.ErrBodyNotAllowed
}

func (failingResponseWriter) WriteHeader(int) {}

func TestHealthHandler(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		collecting bool
		wantStatus int
		wantBody   string
	}{
		{
			name:       "collecting",
			collecting: true,
			wantStatus: http.StatusOK,
			wantBody:   "OK\n",
		},
		{
			name:       "not collecting",
			collecting: false,
			wantStatus: http.StatusServiceUnavailable,
			wantBody:   "Not OK\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			HealthHandler(func() bool { return tt.collecting }).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/__health", nil))

			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if got := recorder.Body.String(); got != tt.wantBody {
				t.Fatalf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}

	// Test error handling with failing response writer (recorder not used as it's expected to fail)
	HealthHandler(func() bool { return true }).ServeHTTP(failingResponseWriter{}, httptest.NewRequest(http.MethodGet, "/__health", nil))
}

func TestStoreHandler(t *testing.T) {
	t.Parallel()
	body := []byte(`{"ok":true}`)
	recorder := httptest.NewRecorder()
	StoreHandler(func() []byte { return body }).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/store", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := recorder.Body.Bytes(); !strings.EqualFold(string(got), string(body)) {
		t.Fatalf("body = %q, want %q", got, body)
	}

	recorder = httptest.NewRecorder()
	StoreHandler(func() []byte { return nil }).ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/store", nil))
	if recorder.Code != http.StatusNotFound || recorder.Body.String() != "Not Found\n" {
		t.Fatalf("missing store response = (%d, %q), want (%d, %q)", recorder.Code, recorder.Body.String(), http.StatusNotFound, "Not Found\n")
	}

	StoreHandler(func() []byte { return body }).ServeHTTP(failingResponseWriter{}, httptest.NewRequest(http.MethodGet, "/store", nil))
}

func TestNew(t *testing.T) {
	t.Parallel()
	mux := New(Config{StoreHTTPPath: "/store"}, func() []byte { return []byte("{}") }, func() bool { return true })

	tests := []struct {
		path       string
		wantStatus int
	}{
		{path: "/metrics", wantStatus: http.StatusOK},
		{path: "/__health", wantStatus: http.StatusOK},
		{path: "/store", wantStatus: http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, tt.path, nil))
			if recorder.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, tt.wantStatus)
			}
		})
	}

	t.Run("store disabled by empty path", func(t *testing.T) {
		t.Parallel()
		mux := New(Config{}, func() []byte { return []byte("{}") }, func() bool { return true })
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/store", nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
		}
	})

	t.Run("store disabled by nil source", func(t *testing.T) {
		t.Parallel()
		mux := New(Config{StoreHTTPPath: "/store"}, nil, func() bool { return true })
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/store", nil))
		if recorder.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
		}
	})
}
