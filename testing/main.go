package main

import "log"

type Shaper interface {
	Foo()
	A() int
}

type Common struct {
	a int
}

func (c *Common) localFoo() {
	log.Println("inside common localFoo")
}

func (c *Common) Foo() {
	log.Println("inside common Foo")
	c.localFoo()

}

func (c *Common) A() int {
	return c.a
}

type Alpha struct {
	C, D float64
	Common
}

func (a *Alpha) Foo() {
	log.Println("inside alpha Foo")
	a.Common.Foo()
	a.localFoo()
}

func (a *Alpha) localFoo() {
	log.Println("inside alpha localFoo")
}

func main() {
	alpha := &Alpha{
		C: 34.5,
		D: 11.1,

		Common: Common{
			a: 10,
		},
	}

	run(alpha)
}

func run(s Shaper) {

	s.Foo()
}
