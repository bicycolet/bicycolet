package middleware

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/go-kit/kit/log"
)

func TestMiddleware(t *testing.T) {
	t.Parallel()

	t.Run("run", func(t *testing.T) {
		builder := New(nil)

		var called []int
		for i := 0; i < 7; i++ {
			func(i int) {
				builder.Add(func(next http.Handler, logger log.Logger) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						called = append(called, i)
						next.ServeHTTP(w, r)
					})
				})
			}(i)
		}

		handler := builder.Build()
		handler.ServeHTTP(nil, nil)

		if expected, actual := []int{0, 1, 2, 3, 4, 5, 6}, called; !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected: %v, actual: %v", expected, actual)
		}
	})
}
