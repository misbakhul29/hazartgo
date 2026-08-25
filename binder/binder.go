package binder

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// BindAndValidate parses path parameters, query parameters, headers, and request body into a target struct, then validates it.
func BindAndValidate(r *http.Request, pathParams map[string]string, target any) error {
	val := reflect.ValueOf(target)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("target must be a pointer to a struct")
	}

	elem := val.Elem()
	typ := elem.Type()

	// 1. Parse JSON Body (if HTTP Method has body)
	if r.Body != nil && r.ContentLength > 0 && r.Header.Get("Content-Type") == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(target); err != nil {
			return fmt.Errorf("invalid json body: %w", err)
		}
	}

	// 2. Parse Query, Path, and Header Tags
	queryVals := r.URL.Query()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := elem.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		// Bind Path parameters (`path:"id"`)
		if pathTag := field.Tag.Get("path"); pathTag != "" {
			if paramVal, exists := pathParams[pathTag]; exists {
				setFieldValue(fieldVal, paramVal)
			}
		}

		// Bind Query parameters (`query:"search"`)
		if queryTag := field.Tag.Get("query"); queryTag != "" {
			if qVal := queryVals.Get(queryTag); qVal != "" {
				setFieldValue(fieldVal, qVal)
			}
		}

		// Bind Header parameters (`header:"X-Api-Key"`)
		if headerTag := field.Tag.Get("header"); headerTag != "" {
			if hVal := r.Header.Get(headerTag); hVal != "" {
				setFieldValue(fieldVal, hVal)
			}
		}
	}

	// 3. Perform Validation via struct tags
	if err := validate.Struct(target); err != nil {
		return err
	}

	return nil
}

func setFieldValue(field reflect.Value, valStr string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(valStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var intVal int64
		fmt.Sscanf(valStr, "%d", &intVal)
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var uintVal uint64
		fmt.Sscanf(valStr, "%d", &uintVal)
		field.SetUint(uintVal)
	case reflect.Bool:
		field.SetBool(valStr == "true" || valStr == "1")
	}
}
