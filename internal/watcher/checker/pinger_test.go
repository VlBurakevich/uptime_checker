package checker

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"uptime-checker/internal/shared/dto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPinger_Ping_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	pinger := NewPinger(5 * time.Second)
	task := dto.SiteCheckTask{
		SiteID:    uuid.New(),
		URL:       server.URL,
		CreatedAt: time.Now(),
	}

	result := pinger.Ping(context.Background(), task)

	assert.True(t, result.IsUp)
	assert.Equal(t, http.StatusOK, result.StatusCode)
	assert.NotZero(t, result.LatencyMs)
}

func TestPinger_Ping_ServerDown(t *testing.T) {
	pinger := NewPinger(5 * time.Second)
	task := dto.SiteCheckTask{
		SiteID: uuid.New(),
		URL:    "http://localhost:9999",
	}

	result := pinger.Ping(context.Background(), task)

	assert.False(t, result.IsUp)
	assert.NotEmpty(t, result.Error)
}

func TestPinger_Ping_Non200(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	pinger := NewPinger(5 * time.Second)
	task := dto.SiteCheckTask{
		SiteID: uuid.New(),
		URL:    server.URL,
	}

	result := pinger.Ping(context.Background(), task)

	assert.False(t, result.IsUp)
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
}

func TestPinger_Ping_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	pinger := NewPinger(100 * time.Millisecond)
	task := dto.SiteCheckTask{
		SiteID: uuid.New(),
		URL:    server.URL,
	}

	result := pinger.Ping(context.Background(), task)

	assert.False(t, result.IsUp)
	assert.NotEmpty(t, result.Error)
}

func TestPinger_Ping_InvalidURL(t *testing.T) {
	pinger := NewPinger(5 * time.Second)
	task := dto.SiteCheckTask{
		SiteID: uuid.New(),
		URL:    "not-a-valid-url",
	}

	result := pinger.Ping(context.Background(), task)

	assert.False(t, result.IsUp)
	assert.NotEmpty(t, result.Error)
}

func TestPinger_Ping_CancelledContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	pinger := NewPinger(5 * time.Second)
	task := dto.SiteCheckTask{
		SiteID: uuid.New(),
		URL:    server.URL,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	result := pinger.Ping(ctx, task)

	assert.False(t, result.IsUp)
}
