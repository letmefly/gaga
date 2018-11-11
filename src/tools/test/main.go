package main

import (
	"log"
)

type Person struct {
	name string
}

func (p *Person) PrintName() {
	log.Println("this is persion")
}

func (p *Person) PrintAge() {
	log.Println("persion age")
}

type Coder struct {
	Person
}

func (c *Coder) PrintName() {
	log.Println("this is coder")
}

func main() {
	coder := Coder{}

	coder.PrintName()
	coder.PrintAge()
}
