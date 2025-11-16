package middleware

import (
	"context"
	"net/http"
	"strconv"
)


type limitType struct{}

type offsetType struct{}


func Paginate(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//достаем урл запроса из обьекта реквест
		q := r.URL.Query()
		// достаем из урла значение параметра page
		page, err := strconv.Atoi(q.Get("page"))
		//обрабатываем ошибку
		if err != nil || page < 0{
			page = 0			
		}
		// достаем из урла значение параметра limit
		limit, err := strconv.Atoi(q.Get("limit"))

		if err != nil || limit < 1 || limit > 100 {
			limit = 20
		}

		offset := page * limit
		// кладем значения параметров в контекст под специально созданный уникальный тип и передаем дальше
		ctx := r.Context()
		ctx = context.WithValue(ctx, limitType{}, limit)
        ctx = context.WithValue(ctx, offsetType{}, offset)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
// достаем значения параметров в контекст из под специально созданный уникальный тип и возвращаем
func GetPagination(r *http.Request) (limit, offset int) {
	ctx := r.Context()

	limit, _ = ctx.Value(limitType{}).(int)
	offset, _ = ctx.Value(offsetType{}).(int)

	return
}


