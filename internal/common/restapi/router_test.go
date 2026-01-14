package restapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
	"github.com/xplexer-lab/xplexer/internal/common/logger"
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

	router := restapi.NewRouter()
	router.SetLogger(logger.NewDummy())
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

func TestRouter_Query(t *testing.T) {

	type getUserQueryOut struct {
		Greet string `json:"greet"`
	}
	var getUserQuery = restapi.Query(func(ctx restapi.Context, in struct {
		Id string `path:"user_id"`
	}) (*getUserQueryOut, error) {
		return &getUserQueryOut{
			Greet: fmt.Sprintf("hello user %s", in.Id),
		}, nil
	})

	r := restapi.NewRouter()
	r.SetLogger(logger.NewDummy())
	r.Get("/user/{user_id}", getUserQuery)
	handler, err := r.BuildHandler()

	require.NoError(t, err)

	server := httptest.NewServer(handler)
	test := httpexpect.Default(t, server.URL)

	var out getUserQueryOut
	test.
		GET("/user/1234").
		Expect().
		Status(http.StatusOK).
		JSON().
		Decode(&out)

	require.Equal(t, "hello user 1234", out.Greet)
}
