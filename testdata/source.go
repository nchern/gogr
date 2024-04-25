package main

import (
	"errors"
	"log"
)

type ComplexInterface interface {
	A(i int)
	B() (int, error)
	C(s string) error

	// D is a multiline method declaration
	D(s string,
		x int,
		y int) (string, error)
}

type Foo struct {
	// multiple fields on one line
	x, y, z float64
}

// a method of Foo struct declared in another file
func (f *Foo) String() string {
	return "abc"
}

type ChainCalls struct {
}

func (c ChainCalls) Fuzz(x string, b *Foo,
	args ...string) ChainCalls {
	return c
}

func (c ChainCalls) Buzz(i ...int) ChainCalls {
	return c
}

func (c ChainCalls) Int() int {
	return 0
}

func (c ChainCalls) Str() (string, error) {
	return "", errors.New(
		"boom" +
			"bam")
}

// exmaple of multiline declaration
func MultilineDecl(
	s string,
	i int,
	z float64) error {

	// multiline call
	log.Printf("foo",
		"bar",
		"fuzz",
		"buzz",
	)

	// multiline if clause
	if isBla := true; isBla &&
		true ||
		false {
		log.Printf("lala")
	}

	chain := ChainCalls{}

	chain.Fuzz("a",
		"b",
		"c").
		Buzz(
			1,
			2,
			3)

	for i := chain.Fuzz(
		"a",
		nil,
		"c").
		Int(); i < 10 ||
		false ||
		true; i++ {
		log.Println("asds")
	}

	return nil
}
