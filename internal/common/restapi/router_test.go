package restapi_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/common/restapi"
)

func TestApiRouter(t *testing.T) {
	type HelloIn struct{}

	type HelloOut struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	}

	var hello = restapi.Query(func(ctx restapi.Context, in HelloIn) (*HelloOut, error) {
		ctx.Logger().Info("hello world")
		return &HelloOut{Message: "Hello World"}, nil
	})

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	router := restapi.NewRouter()
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
