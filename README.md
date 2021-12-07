![Tortilla logo](./tortilla.png)
*(Author: [Micheile](https://unsplash.com/@micheile), original: https://unsplash.com/photos/1zyj8nOdwPs)*

# Tortilla

A Go package to wrap your errors as easily as your tortillas. *Bon appétit!*

# Purpose

Errors in Go is a very divisive topic. There are pros and cons, and for me the thing I really miss
compared to "traditional" exceptions mechanism is the ability to keep an history of the errors you
handled in the lifetime of your program. Although it is possible to wrap or embed errors to keep an
information of what happened in a lower level in vanilla Go with `fmt.Errorf`, I found really
hard to set up a strong but simple, standardized but not too "cumbersome" error handling.

That's why a thought about creating a package to make my error handling easier, and because Go wraps
errors in other errors, (and because I love Tex-Mex cuisine as well), I liked the idea of wrapping all
that things with big Tortillas!

# Examples

OK, enough talk. Let's take some examples to illustrate what you can do with a Tortilla!
These examples can be found in [the examples directory of the project](./examples).

## Wrap and pretty print

In this example we simulate a call to a function that fetches data from a DB and stores it in cache.
The cache returns an error, a Tortilla is created from it and we wrap with our business error
`errDataFetching`. Then we can use `errors.Is` to make a decision regarding the error. Please note that
only the last error used to wrap a Tortilla can be matched by `errors.Is` or `errors.As`.

```go
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
```

Here we simply re-create our Tortilla with `tortilla.New` to call `RollOut().PrettyPrint()`. This will
generate a string with a visual representation of the encountered errors in the reverse order of creation:

```
unable to fetch data:
some cache error:
```

## Adding errors and hierarchy

Now what if you want to add some errors without using them to wrap your Tortilla?

```go
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
			New(errObfuscate).
			Add(err).
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
```

* First, in `encode` function we create a Tortilla `errEncode` and we add the encryption error into it.
* Then in `obsufcate` function a new Tortilla `errObfuscate` wraps the previous Tortilla and add another
error to explain what happaned.
* Finally in `main` the error matched `errObfuscate` and the Tortilla is rolled-out to be displayed in a
human-readable string:

```
unable to obsufcate data:
....some context of what happened
....encoding failed: encryption error.
```

Here we see the hierachy of the errors: the first line is the last error used to wrap the Tortilla.
The second and third are shifted to identify them as "children" of the first one.

As you can see, the third line contains the encoding error followed by the encryption error. It's because
we added the error returned by `encode`, which was flattened.

If you want to keep the exact hierachy, you must create a Torilla from the returned error, and then wrapping
with your own error:

```go
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
```

```
unable to obsufcate data:
....some context of what happened
encoding failed:
....encryption error
```
