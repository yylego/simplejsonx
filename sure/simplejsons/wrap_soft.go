// Code generated using sure/sure_pkg_gen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/sure
// Generated from: sure.gen_test.go:25 -> sure.TestGen
// ========== SURE:DO-NOT-EDIT-SECTION:END ==========

package simplejsons

import (
	"github.com/bitly/go-simplejson"
	"github.com/yylego/simplejsonx"
	"github.com/yylego/sure"
)

func Load(data []byte) (object *simplejson.Json) {
	object, err := simplejsonx.Load(data)
	sure.Soft(err)
	return object
}

func LoadWithOption(data []byte, option *simplejsonx.LoadOption) (object *simplejson.Json) {
	object, err := simplejsonx.LoadWithOption(data, option)
	sure.Soft(err)
	return object
}

func Wrap(value any) (object *simplejson.Json) {
	object = simplejsonx.Wrap(value)
	return object
}

func List(elements []any) (objects []*simplejson.Json) {
	objects = simplejsonx.List(elements)
	return objects
}
