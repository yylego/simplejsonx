// Code generated using sure/sure_pkg_gen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/sure
// Generated from: testing.go:2036 -> testing.tRunner
// ========== SURE:DO-NOT-EDIT-SECTION:END ==========

package simplejsonm

import (
	"github.com/bitly/go-simplejson"
	"github.com/yylego/simplejsonx"
	"github.com/yylego/sure"
)

func Extract[T any](object *simplejson.Json, key string) T {
	res0, err := simplejsonx.Extract[T](object, key)
	sure.Must(err)
	return res0
}

func Inspect[T any](object *simplejson.Json, key string) T {
	res0, err := simplejsonx.Inspect[T](object, key)
	sure.Must(err)
	return res0
}

func Resolve[T any](object *simplejson.Json) T {
	res0, err := simplejsonx.Resolve[T](object)
	sure.Must(err)
	return res0
}

func GetList(object *simplejson.Json, key string) (objects []*simplejson.Json) {
	objects, err := simplejsonx.GetList(object, key)
	sure.Must(err)
	return objects
}

func Inquire[T any](object *simplejson.Json, key string) (T, bool) {
	res0, res1, err := simplejsonx.Inquire[T](object, key)
	sure.Must(err)
	return res0, res1
}

func Attempt[T any](object *simplejson.Json, key string) (T, bool) {
	res0, res1 := simplejsonx.Attempt[T](object, key)
	return res0, res1
}

func Explore[T any](object *simplejson.Json, path string) (T, bool) {
	res0, res1, err := simplejsonx.Explore[T](object, path)
	sure.Must(err)
	return res0, res1
}

func Harvest[T any](object *simplejson.Json, key string, option *simplejsonx.HarvestOption[T]) T {
	res0, err := simplejsonx.Harvest[T](object, key, option)
	sure.Must(err)
	return res0
}
