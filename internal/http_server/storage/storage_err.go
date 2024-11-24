package storage

import "errors"

var (
	ErrGroupExists  = errors.New("group already exists")
	ErrSongExists   = errors.New("song already exists for this group")
	ErrSongNotFound = errors.New("song not found")
)
