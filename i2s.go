package main

import (
	"fmt"
	"reflect"
)

func i2s(data interface{}, out interface{}) error {
	if reflect.ValueOf(out).Kind() != reflect.Ptr {
		return fmt.Errorf("out is not a pointer")
	}
	vOut := reflect.ValueOf(out).Elem()
	tOut := vOut.Type()
	kOut := vOut.Kind()
	switch kOut {
	case reflect.Struct:
		vData, ok := data.(map[string]interface{})
		if !ok {
			return fmt.Errorf("input data must be of type map[string]interface{} for %v\n", vOut)
		}
		for i := 0; i < vOut.NumField(); i++ {
			fieldName := tOut.Field(i).Name
			fieldPtr := vOut.Field(i).Addr()
			subData, ok := vData[fieldName]
			if !ok {
				return fmt.Errorf("no input data provided for %v\n", vOut)
			}
			err := i2s(subData, fieldPtr.Interface())
			if err != nil {
				return err
			}
		}
	case reflect.Slice:
		vData, ok := data.([]interface{})
		if !ok {
			return fmt.Errorf("input data must be of type []interface{} for %v\n", vOut)
		}
		for _, elem := range vData {
			ptr := reflect.New(vOut.Type().Elem())
			err := i2s(elem, ptr.Interface())
			if err != nil {
				return err
			}
			vOut.Set(reflect.Append(vOut, ptr.Elem()))
		}
	case reflect.Bool:
		vData, ok := data.(bool)
		if !ok {
			return fmt.Errorf("input data must be of type bool for %v\n", vOut)
		}
		vOut.SetBool(vData)
	case reflect.Int:
		vData, ok := data.(float64)
		if !ok {
			return fmt.Errorf("input data must be of type float64 for %v\n", vOut)
		}
		vOut.SetInt(int64(vData))
	case reflect.String:
		vData, ok := data.(string)
		if !ok {
			return fmt.Errorf("input data must be of type string for %v\n", vOut)
		}
		vOut.SetString(vData)
	default:
		return fmt.Errorf("unsupported type of input data\n")
	}
	return nil
}
