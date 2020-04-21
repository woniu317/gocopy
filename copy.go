package gocopy

import (
	"errors"
	"reflect"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

func ConvertToTarget(source, target interface{}, sourceTag string,
	convertFunMap map[string]func(gjson.Result) interface{}) (err error) {

	vo := reflect.ValueOf(target)
	if vo.Kind() != reflect.Ptr || vo.IsNil() {
		return errors.New("target must be a pointer")
	}
	//bytes, err := json.Marshal(source)
	bytes, err := jsoniter.Marshal(source)
	vo = vo.Elem()
	targetType := reflect.TypeOf(target).Elem()
	for i := 0; i < vo.NumField(); i++ {
		f := vo.Field(i)
		if !f.IsValid() || !f.CanSet() {
			continue
		}
		tag := targetType.Field(i).Tag.Get(sourceTag)
		originalData := gjson.GetBytes(bytes, tag)

		if originalData.Type == gjson.Null {
			continue
		}

		funcKey := targetType.Field(i).Tag.Get(sourceTag + "FuncKey")
		if funcKey == "" {
			funcKey = targetType.Field(i).Tag.Get(sourceTag)
		}

		if fn, ok := convertFunMap[funcKey]; ok {
			value := fn(originalData)

			switch reflect.TypeOf(value).Kind() {
			case reflect.String:
				f.SetString(reflect.ValueOf(value).String())
			case reflect.Bool:
				f.SetBool(reflect.ValueOf(value).Bool())
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				f.SetInt(reflect.ValueOf(value).Int())
			case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
				f.SetUint(reflect.ValueOf(value).Uint())
			case reflect.Float32, reflect.Float64:
				f.SetFloat(reflect.ValueOf(value).Float())
			default:
				panic("cannot recognize the convert kind " + reflect.TypeOf(value).Kind().String())
			}
		} else {
			switch f.Kind() {
			case reflect.String:
				f.SetString(originalData.Str)
			case reflect.Bool:
				f.SetBool(originalData.Bool())
			case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
				f.SetInt(originalData.Int())
			case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
				f.SetUint(originalData.Uint())
			case reflect.Float32, reflect.Float64:
				f.SetFloat(originalData.Float())
			default:
				panic("cannot recognize the kind " + f.Kind().String())
			}
		}
	}

	return
}
