Validator is a Go language package used for validating struct values.

It has following features:
1. Uses struct tags to perform struct validation.
2. Can also register validation rules by map.
3. Allows validation of private fields.
---
### Example
```go
package main

import (
	"log"
	"github.com/guan-wei-huang/validator"
)

type Test struct {
	Num   int     `validate:"gt=2"`
	Float float32 `validate:"ls=3.4"`
	Str   string  `validate:"required,len=5"`
}

func main() {
	s := Test{
		Num:   4,
		Float: 3.2,
		Str:   "hello",
	}

	v := validator.New()
	if err := v.ValidateStruct(s); err != nil {
		log.Println(err)
	}
}
```

---
#### validation through map
```go

package main

import (
	"log"

	"github.com/guan-wei-huang/validator"
)

type Test struct {
	Num   int
	Float float32
	Str   string
}

func main() {
	s := Test{
		Num:   5,
		Float: 3.2,
		Str:   "hello",
	}

	v := validator.New()
	if err := v.RegisterMapRule(Test{}, map[string]interface{}{
		"Num":   "gt=4",
		"Float": "ls=3.5",
		"str":   "required,len=5",
	}); err != nil {
		log.Fatal(err)
	}

	if err := v.ValidateStruct(s); err != nil {
		log.Fatal(err)
	}
}

```