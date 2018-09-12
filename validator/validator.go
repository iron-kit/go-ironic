package validator

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	validate "gopkg.in/validator.v2"
	"reflect"
)

var hasInitial = false

func initialValidator() {
	if hasInitial {
		return
	}
	validate.SetValidationFunc("objectid", func(v interface{}, params string) error {
		st := reflect.ValueOf(v)

		if st.Kind() != reflect.String {
			return fmt.Errorf("%s must be a string", st.Type().Name())
		}

		if len(v.(string)) == 0 {
			return nil
		}

		if !bson.IsObjectIdHex(st.String()) {
			return fmt.Errorf("%s must be a object id", st.Type().Name())
		}

		return nil
	})
}

// Validate is a custom validte function
func Validate(v interface{}) error {
	initialValidator()
	return validate.Validate(v)
}
