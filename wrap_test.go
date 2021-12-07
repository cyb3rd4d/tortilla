package tortilla

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/suite"
)

type tortillaTestSuite struct {
	suite.Suite
}

func (s *tortillaTestSuite) TestItCreatesATortillaFromAnError() {
	err := errors.New("initial error")
	got := New(err)

	s.ErrorIs(got, err)
}

func (s *tortillaTestSuite) TestItWrapsAnError() {
	err := errors.New("initial error")
	t := New(err)

	wrappedWith := errors.New("wrapping error")
	got := t.Wrap(wrappedWith)

	s.ErrorIs(got, wrappedWith)
	s.EqualError(got, "wrapping error. initial error.")
}

func (s *tortillaTestSuite) TestItAddsErrorsWithoutWrapping() {
	err := errors.New("initial error")
	t := New(err)

	second := errors.New("second")
	got := t.Add(second).Add(errors.New("third"))

	s.ErrorIs(got, err)
	s.NotErrorIs(got, second)
	s.EqualError(got, "initial error: third, second.")
}

func (s *tortillaTestSuite) TestItCreatesATortillaFromATortilla() {
	lastWrap := errors.New("second")
	t1 := New(errors.New("first")).Wrap(lastWrap)
	got := New(t1)

	s.ErrorIs(got, lastWrap)
	s.EqualError(got, "second. first.")
}

func (s *tortillaTestSuite) TestItCanBeRolledOut() {
	t := New(newError("first")).
		Wrap(newError("second")).
		Add(newError("third")).
		Add(newError("fourth")).
		Wrap(newError("fifth")).
		Wrap(newError("sixth")).
		Add(newError("seventh"))

	s.EqualError(t, "sixth: seventh. fifth. second: fourth, third. first.")

	got := t.RollOut()
	expected := Stack{
		{
			"sixth": []string{"seventh"},
		},
		{
			"fifth": []string{},
		},
		{
			"second": []string{"fourth", "third"},
		},
		{
			"first": []string{},
		},
	}

	s.Equal(expected, got)

	expectedPrettyPrint := `sixth:
....seventh
fifth:
second:
....fourth
....third
first:`

	s.Equal(expectedPrettyPrint, got.PrettyPrint())
}

func TestTortilla(t *testing.T) {
	suite.Run(t, new(tortillaTestSuite))
}

func newError(msg string) error {
	return errors.New(msg)
}
