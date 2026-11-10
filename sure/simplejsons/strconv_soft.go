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

func Strconv[T any](object *simplejson.Json) T {
	res0, err := simplejsonx.Strconv[T](object)
	sure.Soft(err)
	return res0
}
