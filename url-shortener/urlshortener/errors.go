package urlshortener

import "errors"

var (
	ErrStorage             = errors.New("storage_error")
	ErrNotFound            = errors.New("not_found")
	ErrKeyGenerationFailed = errors.New("key_generation_failed")
)
