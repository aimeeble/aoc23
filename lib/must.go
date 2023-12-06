package lib

import "log"

func Must[T any](v T, err error) T {
	if err != nil {
		log.Fatalf("Must failure: %v", err)
	}
	return v
}
