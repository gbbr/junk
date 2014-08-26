/*
This package demonstrates the dynamic typing capabilities of the Go programming language
by using its reflect package
*/
package main

import (
	"fmt"
	"reflect"
	"strings"
)

var indent = 1

// Gets a reflect.Value's type. If it contains a package prefix,
// it removes it
func getValueType(val reflect.Value) string {
	parts := strings.Split(val.Type().String(), ".")
	if len(parts) > 1 {
		return parts[1]
	}

	return parts[0]
}

// Describes the fields and values of a structure
// passed as a reflect.value, if further structures
// are found it recurses and continue describing each one.
func describeStruct(object reflect.Value) {
	fmt.Printf(", with fields:\r\n\r\n")

	for i := 0; i < object.NumField(); i++ {
		field := object.Field(i)

		fmt.Print(strings.Repeat(" ", indent*3) + "- ")
		if field.Kind() == reflect.Ptr {
			fmt.Print("pointer to ")
			field = field.Elem()
		}

		fmt.Print(getValueType(field))

		if field.Kind() == reflect.Struct {
			indent++
			describeStruct(field)
			indent--
		}

		fmt.Println()
	}
}

// Prints formatted information about the passed
// value to standard output.
func TellMeAbout(obj interface{}) {
	object := reflect.ValueOf(obj)

	fmt.Print("\r\nYou've passed ")

	// If we have a pointer get to the value of it
	if reflect.TypeOf(obj).Kind() == reflect.Ptr {
		fmt.Print("a pointer to ")
		object = object.Elem()
	}

	// If we have a basic type, say this and return
	if object.Kind() != reflect.Struct {
		fmt.Println("a", getValueType(object))
		return
	}

	fmt.Print(getValueType(object))
	
	// If we have a structure start describing it
	describeStruct(object)

	fmt.Println()
}

type Post struct {
	Id   int
	Name string
	User *User
	age  int32
}

type User struct {
	Id      int64
	Initial rune
	Parents [2]string
}

func main() {
	TellMeAbout(&Post{
		1, "About me", &User{1, '∞', [2]string{"Mom", "Dad"}}, 2,
	})
}
