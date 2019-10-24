package view

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"reflect"
)

// SetField struct field set value
func setField(sval interface{}, name string, val interface{}) error {

	var (
		sv = reflect.ValueOf(sval)
		fv reflect.Value
		ft reflect.Type
	)

	for sv.Kind() == reflect.Ptr {
		sv = sv.Elem()
	}

	fv = sv.FieldByName(name)
	if !fv.IsValid() {
		return fmt.Errorf("no such field: %s in obj", name)
	}
	if !fv.CanSet() {
		return fmt.Errorf("cannot set %s field value", name)
	}

	ft = fv.Type()
	if ft != reflect.ValueOf(val).Type() {
		return fmt.Errorf("provided value type didn't match obj field type")
	}

	fv.Set(reflect.ValueOf(val))

	return nil
}

// ConvertMapToStruct params map fill struct
func ConvertMapToStruct(params map[string]interface{}, val interface{}) {
	for field, fieldVal := range params {
		if err := setField(val, field, fieldVal); err != nil {
			fmt.Println(err)
		}
	}
}

// Capitalize: change first character to upper
// 改变字符串首字母为大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str)
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 {
				vv[i] -= 32
				upperStr += string(vv[i])
			} else {
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}

func Sha1String(key []byte, sum []byte) string {
	s := sha1.New()
	s.Write(key)
	return hex.EncodeToString(s.Sum(sum))
}
