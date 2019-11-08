package main

import "log"

// Shaper is the interface
type Shaper interface {
	Foo()
}

// Common has common components for all concrete type
type Common struct {
	lf func()
}

// this is different for each concrete type, does nothing here
func (c *Common) localFoo() {
	log.Println("inside common localFoo")
}

// Foo is called from the outside, does some work, calls localFoo, does other work
func (c *Common) Foo() {
	log.Println("inside common Foo")

	log.Println("running pre-common")
	c.lf()
	log.Println("running post-common")

}

// Alpha implements Shaper
type Alpha struct {
	Common
}

func (a *Alpha) localFoo() {
	log.Println("inside alpha localFoo")
}

func main() {
	alpha := &Alpha{
		Common: Common{},
	}

	alpha.lf = alpha.localFoo

	run(alpha)
}

func run(s Shaper) {
	s.Foo()
}
