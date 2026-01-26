package restapi_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
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

	var hello = restapi.Operation(func(ctx context.Context, in HelloIn) (*HelloOut, error) {
		logger.FromContext(ctx).Info("hello world")
		return &HelloOut{Message: "Hello World"}, nil
	})

	log := logger.NewDummy()

	router := restapi.NewRouter()
	router.SetLogger(log)
	router.Use(logger.Inject(log))
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

var calls struct {
	mu    sync.Mutex
	calls int
}

func resetCalls() {
	calls.mu.Lock()
	defer calls.mu.Unlock()
	calls.calls = 0
}

func getCalls() int {
	calls.mu.Lock()
	defer calls.mu.Unlock()
	return calls.calls
}

func incrCalls() {
	calls.mu.Lock()
	defer calls.mu.Unlock()
	calls.calls++
}

type sanitizer struct{}

func (sa *sanitizer) Sanitize() {
	incrCalls()
}

func TestRouter_Operation_Common(t *testing.T) {
	resetCalls()

	type getUserQueryOut struct {
		Greet string `json:"greet"`
	}

	type operationIn struct {
		*sanitizer
		Id string `path:"user_id"`
	}

	var getUserQuery = restapi.Operation(func(ctx context.Context, in operationIn) (*getUserQueryOut, error) {
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
	require.Equal(t, 1, getCalls(), "invokes sanitize metthod")
}

type opIn struct {
	Name string `json:"string"`
}

func (in *opIn) Sanitize() {
	in.Name = strings.TrimSpace(in.Name)
}

type opOut struct {
	Name string `json:"message"`
}

func TestOperation(t *testing.T) {
	t.Run("sanitizer works with pointers", func(t *testing.T) {

		r := restapi.NewRouter()
		r.SetLogger(logger.NewDummy())
		r.Post("/operation", restapi.Operation(func(ctx context.Context, in opIn) (*opOut, error) {
			return &opOut{
				Name: fmt.Sprintf("%s", in.Name),
			}, nil
		}))

		handler, err := r.BuildHandler()

		require.NoError(t, err)

		server := httptest.NewServer(handler)
		test := httpexpect.Default(t, server.URL)

		var out opOut
		test.
			POST("/operation").
			WithJSON(&opIn{
				Name: "\t master ",
			}).
			Expect().
			Status(http.StatusOK).
			JSON().
			Decode(&out)

		require.Equal(t, "master", out.Name)
	})
}
