// Code generated using sure/sure_pkg_gen. DO NOT EDIT.
// This file was auto generated via github.com/yylego/sure
// Generated from: zero_tern.go:31 -> zerotern.VF[...]
// ========== SURE:DO-NOT-EDIT-SECTION:END ==========

package simplejsono

import (
	"github.com/bitly/go-simplejson"
	"github.com/yylego/simplejsonx"
	"github.com/yylego/sure"
)

func Load(data []byte) (object *simplejson.Json) {
	object, err := simplejsonx.Load(data)
	sure.Omit(err)
	return object
}

func Wrap(value interface{}) (object *simplejson.Json) {
	object = simplejsonx.Wrap(value)
	return object
}

func List(elements []interface{}) (objects []*simplejson.Json) {
	objects = simplejsonx.List(elements)
	return objects
}
