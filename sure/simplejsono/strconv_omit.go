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

func Strconv[T any](object *simplejson.Json) T {
	res0, err := simplejsonx.Strconv[T](object)
	sure.Omit(err)
	return res0
}
