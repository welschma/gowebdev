package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
	Age  int
	Weight float64
	Slice []string
	Map map[string]string
}

func (u User) SayHello() string {
	return "Hello " + u.Name
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")

	if err != nil {
		panic(err)
	}

	user := User{
		Name: "Jon Calhoun",
		Age:  123,
		Weight: 123.456,
		Slice: []string{"a", "b", "c"},
		Map: map[string]string{"a a": "b", "c": "d"},
	}

	err = t.Execute(os.Stdout, user)

	if err != nil {
		panic(err)
	}
}
