package hazart

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateID generates a random alphanumeric string ID
func GenerateID() string {
	return fmt.Sprintf("%d", rand.Intn(900000)+100000)
}

// ExtractID extracts string ID from a struct using reflection
func ExtractID(entity any) string {
	val := reflect.ValueOf(entity)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.Struct {
		idField := val.FieldByName("ID")
		if idField.IsValid() {
			return fmt.Sprintf("%v", idField.Interface())
		}
	}
	return GenerateID()
}
