package requestmeta

import "context"

type ctxKey struct{}

var ContextKey = ctxKey{}

func IntoContext(ctx context.Context, data *RequestDTO) context.Context {
	return context.WithValue(ctx, ContextKey, data)
}

func FromContext(ctx context.Context) (*RequestDTO, bool) {
	data, ok := ctx.Value(ContextKey).(*RequestDTO)
	return data, ok
}
