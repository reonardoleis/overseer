package database

import "sync"

const (
	FAVORITES = "favorites"
	FUNCTIONS = "functions"
)

var (
	locks = map[string]*sync.Mutex{
		FAVORITES: new(sync.Mutex),
		FUNCTIONS: new(sync.Mutex),
	}
)
