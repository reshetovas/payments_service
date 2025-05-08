package ctxutils

import "context"

type ctxKey string //Чтобы данные не перезатерлись если какая-либо встроенная библиотека решит создать запись с таким же ключом

const userIDKey ctxKey = "user_id"

func WithUserID(ctx context.Context, id int) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

func GetUserID(ctx context.Context) (int, bool) {
	id, ok := ctx.Value(userIDKey).(int)
	return id, ok
}
