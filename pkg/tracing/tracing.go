package tracing

import "golang.org/x/net/context"

type (
	Key int
)

const (
	ClientIDKey Key = iota
)

func GetClientIDFromContext(ctx context.Context) string {
	clientID := ""
	ctxClientID := ctx.Value(ClientIDKey)
	if ctxClientID != nil {
		clientID = ctxClientID.(string)
	}

	return clientID
}

func SetClientIDInContext(ctx context.Context, clientID string) context.Context {
	return context.WithValue(ctx, ClientIDKey, clientID)
}
