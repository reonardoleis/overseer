package database

import "sync"

const (
	FAVORITES = "favorites"
)

var (
	locks = map[string]*sync.Mutex{
		FAVORITES: new(sync.Mutex),
	}
)
