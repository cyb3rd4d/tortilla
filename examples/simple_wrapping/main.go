package main

import (
	"errors"
	"log"

	"github.com/MartialGeek/tortilla"
)

var (
	errDataFetching = errors.New("unable to fetch data")
)

func main() {
	err := fetchSomeData()
	if err != nil {
		if errors.Is(err, errDataFetching) {
			log.Fatal(tortilla.New(err).RollOut().PrettyPrint())
		}

		log.Println("unknown error:", err)
	}
}

func fetchSomeData() error {
	cacheErr := cache()
	if cacheErr != nil {
		return tortilla.New(cacheErr).Wrap(errDataFetching)
	}

	return nil
}

func cache() error {
	return errors.New("some cache error")
}
