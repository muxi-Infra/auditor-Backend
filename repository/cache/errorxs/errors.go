package errorxs

import (
	"errors"
	"fmt"
)

type CacheNotFoundError struct {
	Err error
}

func (e CacheNotFoundError) Error() string {
	return fmt.Sprintf("cache not found: %v", e.Err)
}

func (e CacheNotFoundError) Unwrap() error {
	return e.Err
}

func ToCacheNotFoundError(err error) error {
	return CacheNotFoundError{Err: err}
}
func IsCacheNotFoundError(err error) bool {
	var target CacheNotFoundError
	return errors.As(err, &target)
}
