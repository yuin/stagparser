# stagparser - generic parser for golang struct tag
[![GoDoc](https://godoc.org/github.com/yuin/stagparser?status.svg)](https://godoc.org/github.com/yuin/stagparser)

Package stagparser provides a generic parser for golang struct tag.
stagparser can parse tags like the following:

- `validate:"required,length(min=1, max=10)"`
- `validate:"max=10,list=[apple,'star fruits']"`

tags consist of 'definition'. 'definition' has 3 forms:

- name only: `required`
- name with a single attribute: `max=10`
    - in this case, parse result is name=`"max"`, attributes=`{"max":10}`
- name with multiple attributes: `length(min=1, max=10)`

name and attribute must be a golang identifier.
An attribute value must be one of an int64, a float64, an identifier,
a string quoted by `'` and an array.

* int64: `123`
* float64: `111.12`
* string: `'ab\tc'`
  * identifiers are interpreted as string in value context
* array: `[1, 2, aaa]`

You can parse objects just calling ParseStruct:

```
import "github.com/yuin/stagparser"

type User struct {
  Name string `validate:"required,length(min=4,max=10)"`
}

func main() {
  user := &User{"bob"}
  definitions, err := stagparser.ParseStruct(user)
}
```

## License
MIT

## Author
Yusuke Inuzuka
