package gorm

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/minipkg/selection_condition"
)

const DefaultLimit = 1000

func Conditions(db *gorm.DB, conditions *selection_condition.SelectionCondition) *gorm.DB {
	if conditions == nil {
		return db
	}

	if err := conditions.Validate(); err != nil {
		db.AddError(err)
		return db
	}

	db = Where(db, conditions.Where)
	db = SortOrder(db, conditions.SortOrder)
	db = Limit(db, conditions.Limit)
	db = Offset(db, conditions.Offset)

	return db
}

func SortOrder(db *gorm.DB, orders []map[string]string) *gorm.DB {
	if orders == nil {
		return db
	}

	if db.Statement.Schema == nil {
		db.AddError(errors.New("Schema is nil"))
		return db
	}

	for _, order := range orders {
		s := strings.Builder{}

		for k, v := range order {
			field, ok := db.Statement.Schema.FieldsByName[k]
			if !ok {
				db.AddError(errors.Errorf("Can not find a field %q", k))
				return db
			}

			if field.DBName == "" {
				db.AddError(errors.Errorf("DBName in model must be specified"))
				return db
			}
			tableField := field.DBName

			s.WriteString(tableField + " " + v + ", ")
		}
		db = db.Order(strings.Trim(s.String(), ", "))
	}
	return db
}

func Offset(db *gorm.DB, value uint) *gorm.DB {
	if value == 0 {
		return db
	}
	return db.Offset(int(value))
}

func Limit(db *gorm.DB, value uint) *gorm.DB {
	if value == 0 {
		return db.Limit(DefaultLimit)
	}
	return db.Limit(int(value))
}

func Where(db *gorm.DB, conditions interface{}) *gorm.DB {
	if conditions == nil {
		return db
	}

	wcs, ok := conditions.(selection_condition.WhereConditions)
	if ok {
		return WhereConditions(db, wcs)
	}

	wc, ok := conditions.(selection_condition.WhereCondition)
	if ok {
		return WhereCondition(db, wc)
	}

	if !isStruct(conditions) {
		db.AddError(errors.Errorf("conditions must be a selection_condition.WhereConditions, selection_condition.WhereCondition or a struct"))
		return db
	}
	return db.Where(conditions)
}

func isStruct(e interface{}) bool {
	t := reflect.TypeOf(e)

	if t.Kind() == reflect.Ptr {
		t = reflect.Indirect(reflect.ValueOf(e)).Type()
	}
	return t.Kind() == reflect.Struct
}

func WhereConditions(db *gorm.DB, conditions selection_condition.WhereConditions) *gorm.DB {
	if err := conditions.Validate(); err != nil {
		db.AddError(err)
		return db
	}

	for _, condition := range conditions {
		db = WhereCondition(db, condition)
		if db.Error != nil {
			return db
		}
	}
	return db
}

func WhereCondition(db *gorm.DB, condition selection_condition.WhereCondition) *gorm.DB {
	if err := condition.Validate(); err != nil {
		db.AddError(err)
		return db
	}

	if db.Statement.Schema == nil {
		db.AddError(errors.New("Schema is nil"))
		return db
	}

	field, ok := db.Statement.Schema.FieldsByName[condition.Field]
	if !ok {
		db.AddError(errors.Errorf("Can not find a field %q", condition.Field))
		return db
	}

	if field.DBName == "" {
		db.AddError(errors.Errorf("DBName in model must be specified"))
		return db
	}
	tableField := field.Schema.Table + "." + field.DBName

	switch condition.Condition {
	case selection_condition.ConditionEq:
		db = db.Where(map[string]interface{}{tableField: condition.Value})
	case selection_condition.ConditionIn:
		conds, ok := condition.Value.([]interface{})
		if !ok {
			db.AddError(errors.Errorf("Can not assign value condition to slice"))
		}
		db = db.Where(tableField+" IN (?)", conds)
	case selection_condition.ConditionBt:
		conds, ok := condition.Value.([]interface{})
		if !ok {
			db.AddError(errors.Errorf("Can not assign value condition to slice"))
		}
		db = db.Where(tableField+" BETWEEN ? AND ?", conds[0], conds[1])
	case selection_condition.ConditionGt:
		db = db.Where(tableField+" > ?", condition.Value)
	case selection_condition.ConditionGte:
		db = db.Where(tableField+" >= ?", condition.Value)
	case selection_condition.ConditionLt:
		db = db.Where(tableField+" < ?", condition.Value)
	case selection_condition.ConditionLte:
		db = db.Where(tableField+" <= ?", condition.Value)
	}
	return db
}

func keysToSnakeCase(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))

	for key, val := range in {
		out[strcase.ToSnake(key)] = val
	}
	return out
}

func keysToSnakeCaseStr(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))

	for key, val := range in {
		out[strcase.ToSnake(key)] = val
	}
	return out
}
