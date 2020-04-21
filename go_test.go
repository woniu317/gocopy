package gocopy_test

import (
	"gocopy"
	"testing"

	"github.com/tidwall/gjson"
)

type base struct {
	Name string `json:"name"`
	Sex  int    `json:"sex"`
}

type source struct {
	B base `json:"base"`
}

type target struct {
	Name int `cp:"base.name" cpFuncKey:"base.name"`
	Sex  int `cp:"base.sex"`
}

func BenchmarkCopy(t *testing.B) {
	a := source{B: base{Name: "123", Sex: 3}}
	b := target{Name: 0}
	for i := 0; i < t.N; i++ {
		_ = gocopy.ConvertToTarget(a, &b, "cp", map[string]func(result gjson.Result) interface{}{
			"base.name": func(data gjson.Result) interface{} {
				return data.Int()
			},
		})
	}
}
