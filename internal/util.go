package internal

import (
	"fmt"
	"log"
	"reflect"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
)

func ColorizeStringRed(source string) string {
	return colorRed + source + colorReset
}

func ColorizeStringGreen(source string) string {
	return colorGreen + source + colorReset
}

func ColorizeStringYellow(source string) string {
	return colorYellow + source + colorReset
}

func ColorizeStringBlue(source string) string {
	return colorBlue + source + colorReset
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}

	return false
}

func IsString(i interface{}) bool {
	switch i.(type) {
	case string, *string:
		return true
	}

	return false
}

func GetActualInterfaceValue(i interface{}) any {
	rv := reflect.ValueOf(i)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv.Interface()
}

func GetKeys[K comparable, V any](source map[K]V) []K {
	keys := make([]K, 0, len(source))

	for key := range source {
		keys = append(keys, key)
	}

	return keys
}

func PrintFatal(message string, error error) {
	if error == nil && message == "" {
		return
	}

	if message == "" {
		log.Fatal(ColorizeStringRed(error.Error()))
		return
	}

	if error == nil {
		log.Fatal(ColorizeStringRed(message))
		return
	}

	log.Fatalf(ColorizeStringRed(fmt.Sprintf("%s: %s", message, error.Error())))
}

func PrintMessage(message string, args ...any) {
	if message == "" {
		return
	}

	log.Printf(ColorizeStringGreen(message), args...)
}

func PrintWarning(message string, args ...any) {
	if message == "" {
		return
	}

	log.Printf(ColorizeStringYellow(message), args...)
}

func PrintError(message string, error error) {
	if error == nil && message == "" {
		return
	}

	if message == "" {
		log.Println(ColorizeStringRed(error.Error()))
		return
	}

	if error == nil {
		log.Println(ColorizeStringRed(message))
		return
	}

	log.Println(ColorizeStringRed(fmt.Sprintf("%s: %s", message, error.Error())))
}

func ToInterface[T any](source T) interface{} {
	var result interface{} = source
	return result
}

func Clone(oldObj interface{}) interface{} {
	newObj := reflect.New(reflect.TypeOf(oldObj).Elem())
	oldVal := reflect.ValueOf(oldObj).Elem()
	newVal := newObj.Elem()
	for i := 0; i < oldVal.NumField(); i++ {
		newValField := newVal.Field(i)
		if newValField.CanSet() {
			newValField.Set(oldVal.Field(i))
		}
	}

	return newObj.Interface()
}

func MapArray[T, R any](source []T, mapper func(T) R) []R {
	if len(source) == 0 {
		return make([]R, 0)
	}

	result := make([]R, 0, len(source))
	for i := range source {
		result = append(result, mapper(source[i]))
	}

	return result
}

func UnNilArray[T any](source []T) []T {
	result := make([]T, 0)

	for i := range source {
		if IsNil(source[i]) {
			continue
		}

		result = append(result, source[i])
	}

	return result
}
