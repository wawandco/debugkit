package main

import "github.com/wawandco/debugkit"

type Address struct {
	City  string
	State string
}

type Person struct {
	Name    string
	Age     int
	Address Address
	Tags    []string
}

func main() {
	p := Person{
		Name: "Alice",
		Age:  30,
		Address: Address{
			City:  "NYC",
			State: "NY",
		},
		Tags: []string{"developer", "gopher"},
	}

	// test for more nested structs
	type Company struct {
		Name      string
		CEO       Person
		Employees []Person
	}

	c := Company{
		Name: "Tech Co",
		CEO:  p,
		Employees: []Person{
			p,
			{
				Name: "Bob",
				Age:  25,
				Address: Address{
					City:  "SF",
					State: "CA",
				},
				Tags: []string{"designer", "artist"},
			},
			{
				Name: "Charlie",
				Age:  28,
				Address: Address{
					City:  "LA",
					State: "CA",
				},
				Tags: []string{"manager", "leader"},
			},
		},
	}

	debugkit.Dump(&c)
}
