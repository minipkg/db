package mongo

import (
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func QueryWhereCondition(pointerOnEntity interface{}) bson.M {
	condition := bson.M{}
	value := reflect.ValueOf(pointerOnEntity).Elem()
	t := value.Type()
	empty := reflect.New(t).Elem()

	for i := 0; i < t.NumField(); i++ {
		switch value.Field(i).Kind() {
		case reflect.Slice, reflect.Map, reflect.Struct:
			continue
		}

		fieldValue := value.Field(i).Interface()
		emptyFieldValue := empty.Field(i).Interface()

		if fieldValue != emptyFieldValue {
			condition[strings.ToLower(t.Field(i).Name)] = fieldValue
		}
	}

	return condition
}
