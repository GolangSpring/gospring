package helper

import (
	"fmt"
	"reflect"
)

func Stringify(obj any) string {
	// Use reflection on the value of the struct passed to the method (dynamic receiver)
	val := reflect.ValueOf(obj)

	// Ensure we're working with a struct
	if val.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct.")
		// return a string representation of the value
		return fmt.Sprintf("%v", obj)
	}

	_type := reflect.TypeOf(obj)
	info := _type.Name() + "(" // Start with the struct name and an opening parenthesis

	// Iterate through the struct fields and print their names and values
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := _type.Field(i).Name
		fieldValue := field.Interface()
		// Format each field with name and value
		fieldInfo := fmt.Sprintf("%s: %v", fieldName, fieldValue)

		// Append a comma if it's not the last field
		if i < val.NumField()-1 {
			info += fieldInfo + ", "
		} else {
			info += fieldInfo // No comma for the last field
		}
	}

	info += ")" // Close the parentheses
	return info
}
