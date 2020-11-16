package binding_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/summerKK/go-code-snippet-library/gin-study/binding"
)

func TestValidate(t *testing.T) {
	type p0 struct {
		Name string `binding:"required"`
		Age  int    `binding:"required"`
	}

	type p1 struct {
		P0 p0   `binding:"required"`
		P1 []p0 `binding:"required"`
	}

	v := p1{
		P0: p0{
			Name: "summer0",
			Age:  28,
		},
		P1: []p0{
			{
				Name: "summer1",
				Age:  28,
			},
			{
				Name: "summer2",
				Age:  0,
			},
		},
	}

	err := binding.Validate(v)

	assert.NotEqual(t, err, nil)
	assert.IsEqual("Required Age", err.Error())
}
