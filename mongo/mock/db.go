package mock

import (
	"context"
	"fmt"
	"reflect"

	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongodb "github.com/minipkg/db/mongo"
)

type DB struct {
	mock.Mock
}

func (m DB) Collection(a0 string, a1 ...*options.CollectionOptions) mongodb.ICollection {
	ret := m.Called(a0, a1)

	var r0 mongodb.ICollection
	if rf, ok := ret.Get(0).(func(string, ...[]*options.CollectionOptions) mongodb.ICollection); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongodb.ICollection)
		}
	}

	return r0
}

func (b DB) Close(ctx context.Context) error {
	return nil
}

type Collection struct {
	mock.Mock
	DefaultSingleResult mongodb.ISingleResult
	DefaultCursor       mongodb.ICursor
}

func (m Collection) FindOne(a0 context.Context, a1 interface{}, a2 ...*options.FindOneOptions) mongodb.ISingleResult {
	ret := m.Called(a0, a1, a2)

	var r0 mongodb.ISingleResult
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...[]*options.FindOneOptions) mongodb.ISingleResult); ok {
		r0 = rf(a0, a1, a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongodb.ISingleResult)
		}
	}

	return r0
}

func (m Collection) Find(a0 context.Context, a1 interface{}, a2 ...*options.FindOptions) (mongodb.ICursor, error) {
	ret := m.Called(a0, a1, a2)

	var r0 mongodb.ICursor
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, ...[]*options.FindOptions) mongodb.ICursor); ok {
		r0 = rf(a0, a1, a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mongodb.ICursor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, ...[]*options.FindOptions) error); ok {
		r1 = rf(a0, a1, a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m Collection) InsertOne(a0 context.Context, a1 interface{}) (interface{}, error) {
	ret := m.Called(a0, a1)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) interface{}); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m Collection) UpdateOne(a0 context.Context, a1 interface{}, a2 interface{}) (interface{}, error) {
	ret := m.Called(a0, a1, a2)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, interface{}, interface{}) interface{}); ok {
		r0 = rf(a0, a1, a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}, interface{}) error); ok {
		r1 = rf(a0, a1, a2)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (m Collection) DeleteOne(a0 context.Context, a1 interface{}) (int64, error) {
	ret := m.Called(a0, a1)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, interface{}) int64); ok {
		r0 = rf(a0, a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(int64)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, interface{}) error); ok {
		r1 = rf(a0, a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type SingleResult struct {
	Entity interface{}
	Err    error
}

func (m SingleResult) Decode(out interface{}) error {

	if err := setByPtr(m.Entity, out); err != nil {
		return err
	}
	return m.Err
}

type Cursor struct {
	next int
	Res  []interface{}
}

func (m *Cursor) Next(ctx context.Context) bool {
	m.next++
	return m.next <= len(m.Res)
}

func (m Cursor) Decode(out interface{}) error {

	if err := setByPtr(m.Res[m.next-1], out); err != nil {
		return err
	}
	return nil
}

func setByPtr(in interface{}, out interface{}) error {
	outType := reflect.TypeOf(out)

	if outType.Kind() != reflect.Ptr {
		return fmt.Errorf("Parameter out must be a Ptr")
	}

	outVal := reflect.ValueOf(out)
	outValElem := outVal.Elem()

	if !outValElem.CanSet() {
		return fmt.Errorf("!outValElem.CanSet()")
	}

	inType := reflect.TypeOf(in)

	if inType.Kind() != reflect.Ptr {
		return fmt.Errorf("Entity with data must be a Ptr")
	}

	inVal := reflect.ValueOf(in)
	inValElem := inVal.Elem()

	outValElem.Set(inValElem)

	return nil
}

var _ mongodb.IDB = (*DB)(nil)
var _ mongodb.ICollection = (*Collection)(nil)
var _ mongodb.ISingleResult = (*SingleResult)(nil)
var _ mongodb.ICursor = (*Cursor)(nil)
