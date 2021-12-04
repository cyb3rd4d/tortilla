package tortilla

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

const stackPrettyPrintTpl = `
{{- range $layer := .}}{{range $target, $chain := .}}{{$target}}:{{range $chain}}
....{{.}}{{end}}
{{end}}{{end}}`

var parsedTpl *template.Template

func init() {
	tpl, err := template.New("pretty-print").Parse(stackPrettyPrintTpl)
	if err != nil {
		err = fmt.Errorf("tortilla: unable to parse the pretty print tpl: %s", err)
		panic(err)
	}

	parsedTpl = tpl
}

// Stack is an alias for the slice of layers returned by Tortilla.RollOut()
type Stack []map[string][]string

// PrettyPrint allows to render a visual string of the error stack.
// The format is a slice of maps, each map is "wrapping error": []list of wrapped error. Example:
//
// last error used to wrap others:
// ....last wrapped error
// ....older error
// older wrapping error:
// ....blablabla
func (s Stack) PrettyPrint() string {
	output := new(bytes.Buffer)
	err := parsedTpl.Execute(output, s)
	if err != nil {
		return "Pretty print error"
	}

	return strings.TrimSpace(output.String())
}

type layer struct {
	target error
	chain  []error
}

// Tortilla holds the layers of the errors added in the stack.
// Create a Tortilla with New(err).
//
// Then use Wrap(err) to wrap your Tortilla with the a new error. This is the
// equivalent of using fmt.Errorf("%w %s", wrapWithErr, initialErr)
//
// You can also add an error in the stack without wrapping with the Add(err) method. This
// can be useful to keep the history of what happened in your program without "typing"
// with (i.e. errors.Is or errors.As won't return true).
//
// If an error is printed with Error(), a string will be generated with the errors in the
// stack in an inlined form. The errors are sorted in reverse order of creation.
//
// Of course a Tortilla can be rolled out! The RollOut method returns the layers of your
// Tortilla as a Stack type (an alias of []map[string][]string). Then you can use the
// method Stack.PrettyPrint method to display a hierarchy of the errors wrapped and added
// in your Tortilla lifetime.
type Tortilla struct {
	layers []layer
}

// Error returns a flattened string composed by the errors in the stack.
func (t Tortilla) Error() string {
	var msg string
	for _, layer := range t.layers {
		msg += buildLayerMsg(layer)
	}

	return strings.TrimSpace(msg)
}

// Unwrap allows to compare a Tortilla with errors.Is and errors.As.
func (t Tortilla) Unwrap() error {
	return t.layers[0].target
}

// Wrap wraps your Tortilla with a new error.
func (t Tortilla) Wrap(target error) Tortilla {
	return t.appendLayer(target, nil)
}

// Add adds a new error in the error chain of the last wrapping.
func (t Tortilla) Add(err error) Tortilla {
	chain := []error{err}
	chain = append(chain, t.layers[0].chain...)
	layers := t.layers
	layers[0].chain = chain

	return Tortilla{layers: layers}
}

// RollOut allows you to see what's inside your Tortilla.
func (t Tortilla) RollOut() Stack {
	s := make(Stack, 0, len(t.layers))
	for _, l := range t.layers {
		c := make([]string, 0, len(l.chain))
		for _, e := range l.chain {
			c = append(c, e.Error())
		}

		layer := map[string][]string{
			l.target.Error(): c,
		}

		s = append(s, layer)
	}

	return s
}

func (t Tortilla) appendLayer(target error, chain []error) Tortilla {
	layers := []layer{{target: target, chain: chain}}
	layers = append(layers, t.layers...)

	return Tortilla{
		layers: layers,
	}
}

// New creates a new Tortilla from the given error.
// If err is an existing Tortilla, the returned value is that Tortilla.
// Otherwise a new Tortilla is created from err.
func New(err error) Tortilla {
	t, ok := err.(Tortilla)
	if !ok {
		return Tortilla{
			layers: []layer{
				{target: err},
			},
		}
	}

	return t
}

func buildChainMsg(chain []error) (msg string) {
	for _, err := range chain {
		msg += err.Error() + ", "
	}

	msg = strings.TrimRight(msg, ", ")
	return
}

func buildLayerMsg(layer layer) (msg string) {
	msg = layer.target.Error() + ": "
	msg += buildChainMsg(layer.chain)
	msg = strings.TrimRight(msg, ": ")
	msg += ". "

	return
}
