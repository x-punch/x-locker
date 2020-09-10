package locker

import "errors"

var (
	// ErrInvalidLockGroup represents lock group empty error
	ErrInvalidLockGroup = errors.New("lock group cannot by empty")
	// ErrLockGroupNotFound represents lock group not found error
	ErrLockGroupNotFound = errors.New("lock group not found")
)
