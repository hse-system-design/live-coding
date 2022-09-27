package urlshortener

import "context"

type Manager interface {
	CreateShortcut(ctx context.Context, fullURL string) (string, error)

	// ResolveShortcut returns ErrNotFound if there is no shortcut with the specified key.
	ResolveShortcut(ctx context.Context, key string) (string, error)
}
