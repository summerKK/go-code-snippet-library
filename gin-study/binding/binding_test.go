package binding_test

import (
	"log"
	"testing"

	"github.com/matryer/is"
	"github.com/summerKK/go-code-snippet-library/gin-study/binding"
)

func TestValidate(t *testing.T) {
	assertIs := is.New(t)

	type p0 struct {
		Name string `binding:"required"`
		Age  int    `binding:"required"`
	}

	type p2 struct {
		Items []p0 `binding:"required"`
	}

	type p3 struct {
		Items p0 `binding:"required"`
	}

	type p1 struct {
		P0 p0     `binding:"required"`
		P1 []p0   `binding:"required"`
		P3 [][]p0 `binding:"required"`
		P4 string `binding:"required"`
		P5 p2     `binding:"required"`
		P6 p3     `binding:"required"`
		p7 string `binding:"required"`
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
				Age:  29,
			},
		},
		P3: [][]p0{
			[]p0{
				{
					Name: "summer3",
					Age:  30,
				},
				{
					Name: "summer4",
					Age:  31,
				},
			},
		},
		P4: "hello,world",
		P5: p2{
			Items: []p0{
				{
					Name: "summer5",
					Age:  32,
				},
			},
		},
		P6: p3{
			Items: p0{
				Name: "summer6",
				Age:  33,
			},
		},
	}

	err := binding.Validate(v)
	log.Println(err)

	assertIs.True(err == nil)

}
