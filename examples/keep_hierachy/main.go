package main

import (
	"errors"
	"log"

	"github.com/MartialGeek/tortilla"
)

var (
	errObfuscate = errors.New("unable to obsufcate data")
	errEncode    = errors.New("encoding failed")
)

func main() {
	err := obsufcate()
	if err != nil {
		if errors.Is(err, errObfuscate) {
			log.Fatal(tortilla.New(err).RollOut().PrettyPrint())
		}
	}
}

func obsufcate() error {
	err := encode()
	if err != nil {
		err = tortilla.
			New(err).
			Wrap(errObfuscate).
			Add(errors.New("some context of what happened"))
	}

	return err
}

func encode() error {
	err := encrypt()
	if err != nil {
		err = tortilla.New(errEncode).Add(err)
	}

	return err
}

func encrypt() error {
	return errors.New("encryption error")
}
