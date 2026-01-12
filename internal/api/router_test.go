package api_test

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/api"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

type HelloIn struct{}

type HelloOut struct {
	Message string `json:"message"`
	Name    string `json:"name"`
}

var hello = api.Query(func(ctx api.Context, in HelloIn) (*HelloOut, error) {
	ctx.Logger().Info("hello world")
	return &HelloOut{Message: "Hello World"}, nil
})

func TestApiRouter(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	router := api.NewRouter()
	router.SetLogger(logger)
	router.Get("/hello", hello)
	handler, err := router.BuildHandler()

	require.NoError(t, err)
	s := httptest.NewServer(handler)
	defer s.Close()

	resp, err := http.Get(s.URL + "/hello")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var out HelloOut
	err = json.NewDecoder(resp.Body).Decode(&out)
	require.NoError(t, err)
	require.Equal(t, "Hello World", out.Message)
}

func TestApiRouter_NotFound(t *testing.T) {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Method(http.MethodGet, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	}))
	s := httptest.NewServer(r)
	defer s.Close()
	resp, err := http.Get(s.URL + "/")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}
