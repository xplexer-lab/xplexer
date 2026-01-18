package binder

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBinder_Bind(t *testing.T) {
	b := NewDefault()

	t.Run("Scalar Types", func(t *testing.T) {
		type Request struct {
			Name   string  `query:"name"`
			Age    int     `query:"age"`
			Active bool    `query:"active"`
			Score  float64 `query:"score"`
		}

		r := httptest.NewRequest(http.MethodGet, "/?name=Gopher&age=10&active=true&score=99.9", nil)
		var req Request

		err := b.Bind(r, &req)

		require.NoError(t, err)
		require.Equal(t, "Gopher", req.Name)
		require.Equal(t, 10, req.Age)
		require.True(t, req.Active)
		require.InDelta(t, 99.9, req.Score, 0.0001)
	})

	t.Run("Slices support", func(t *testing.T) {
		type Request struct {
			Tags []string `query:"tag"`
			Nums []int    `query:"num"`
		}

		r := httptest.NewRequest(http.MethodGet, "/?tag=a&tag=b&num=1&num=2", nil)
		var req Request

		err := b.Bind(r, &req)
		require.NoError(t, err)

		require.Len(t, req.Tags, 2)
		require.Equal(t, []string{"a", "b"}, req.Tags)

		require.Len(t, req.Nums, 2)
		require.Equal(t, []int{1, 2}, req.Nums)
	})

	t.Run("Precedence Logic", func(t *testing.T) {
		type Request struct {
			PathFirst       string `path:"my_key" query:"my_key"`
			QueryFirst      string `query:"my_key" path:"my_key"`
			FallbackToQuery string `path:"missing_key" query:"my_key"`
		}

		r := httptest.NewRequest(http.MethodGet, "/?my_key=value_from_query", nil)
		r.SetPathValue("my_key", "value_from_path")

		var req Request
		err := b.Bind(r, &req)
		require.NoError(t, err)

		require.Equal(t, "value_from_path", req.PathFirst)
		require.Equal(t, "value_from_query", req.QueryFirst)
		require.Equal(t, "value_from_query", req.FallbackToQuery)
	})

	t.Run("Nested Structs & Pointers", func(t *testing.T) {
		type Metadata struct {
			RequestID string `header:"X-Request-ID"`
		}

		type Request struct {
			Meta    Metadata
			MetaPtr *Metadata
		}

		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Request-ID", "123-abc")

		var req Request
		err := b.Bind(r, &req)
		require.NoError(t, err)

		require.Equal(t, "123-abc", req.Meta.RequestID)

		require.NotNil(t, req.MetaPtr)
		require.Equal(t, "123-abc", req.MetaPtr.RequestID)
	})

	t.Run("Validation & Errors", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		err := b.Bind(r, nil)
		require.Error(t, err)
		require.Contains(t, err.Error(), "non-nil pointer")

		type Req struct{}
		var valReq Req
		err = b.Bind(r, valReq)
		require.Error(t, err)
		require.Contains(t, err.Error(), "non-nil pointer")

		type IntReq struct {
			Val int `query:"val"`
		}
		rBad := httptest.NewRequest(http.MethodGet, "/?val=not_a_number", nil)
		var intReq IntReq
		err = b.Bind(rBad, &intReq)
		require.Error(t, err)
	})
}
