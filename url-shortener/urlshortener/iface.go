package urlshortener

type Manager interface {
	CreateShortcut(fullURL string) (string, error)
	ResolveShortcut(key string) (string, bool)
}
