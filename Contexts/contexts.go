package contexts

import (
	"context"
	"net/http"
)

//Эта middleware функция позволяет обработчику next обращаться к айди аутентифицированного пользователя через контекст (из main.go)
func ServeWithContext(next http.HandlerFunc, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)

		next(w, r)
	}
}
