// Package stagparser provides a generic parser for golang struct tag.
// stagparser can parse tags like the following:
//
//   - `validate:"required,length(min=1, max=10)"`
//   - `validate:"max=10,list=[apple,'star fruits']"`
//
// tags are consists of 'definition'. 'definition' have 3 forms:
//
//   - name only: required
//   - name with a single attribute: max=10
//   - in this case, parse result is name="max", attributes={"max":10}
//   - name with multiple attributes: length(min=1, max=10)
//
// name and attribute must be a golang identifier.
// An attribute value must be one of an int64, a float64, an identifier,
// a string quoted by "'" and an array.
//
//   - int64: 123
//   - float64: 111.12
//   - string: 'ab\tc'
//   - identifiers are interpreted as string in value context
//   - array:  [1, 2, aaa]
//
// You can parse objects just call ParseStruct:
//
//	import "github.com/yuin/stagparser"
//
//	type User struct {
//	  Name string `validate:"required,length(min=4,max=10)"`
//	}
//
//	func main() {
//	  user := &User{"bob"}
//	  definitions, err := stagparser.ParseStruct(user)
//	}
package stagparser

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"
)

// ParseError is an error indicating invalid tag value.
type ParseError interface {
	error
	// Source is a source name
	Source() string

	// Column is a column error occurred
	Column() int

	// Line is a line number error occurred
	Line() int
}

type parseError struct {
	message string
	source  string
	column  int
	line    int
}

func (e *parseError) Error() string {
	return fmt.Sprintf("%s (%d:%d [%s])", e.message, e.line, e.column, e.source)
}

func (e *parseError) Source() string {
	return e.source
}

func (e *parseError) Column() int {
	return e.column
}

func (e *parseError) Line() int {
	return e.line
}

type parser struct {
	source string
	s      scanner.Scanner
}

func newParser(source string) *parser {
	p := &parser{
		source: source,
	}
	p.s.Mode = scanner.ScanIdents | scanner.ScanInts | scanner.ScanFloats | scanner.ScanStrings
	return p
}

func (p *parser) Parse(tag string) ([]Definition, error) {
	p.s.Init(strings.NewReader(tag))
	result := []Definition{}
	for {
		tok := p.s.Scan()
		switch tok {
		case scanner.EOF:
			return result, nil
		case scanner.Ident:
			ident := p.s.TokenText()
			if p.s.Peek() == '=' {
				_ = p.s.Next()
				value, err := p.parseValue()
				if err != nil {
					return nil, err
				}
				arg := map[string]interface{}{
					ident: value,
				}
				result = append(result, newDefinition(ident, arg))
			} else if p.s.Peek() == '(' {
				_ = p.s.Next()
				arg, err := p.parseArgs()
				if err != nil {
					return nil, err
				}
				result = append(result, newDefinition(ident, arg))
			} else if p.s.Peek() == scanner.EOF || p.s.Peek() == ',' {
				result = append(result, newDefinition(ident, map[string]interface{}{}))
			}
		case ',':
			// NOP
		default:
			return nil, p.parseError(fmt.Sprintf("invalid token: %s", p.s.TokenText()))
		}
	}
}

func (p *parser) parseError(message string) error {
	return &parseError{
		message: message,
		column:  p.s.Position.Column,
		line:    p.s.Position.Line,
	}
}

func (p *parser) parseValue() (interface{}, error) {
	switch p.s.Peek() {
	case '\'':
		return p.parseString(p.s.Next())
	case '[':
		return p.parseArray(p.s.Next())
	default:
		tok := p.s.Scan()
		switch tok {
		case scanner.Ident:
			return p.s.TokenText(), nil
		case scanner.String, scanner.Int, scanner.Float, '-':
			mul := 1
			if tok == '-' {
				mul = -1
				tok = p.s.Scan()
			}
			if tok == scanner.String {
				str := p.s.TokenText()
				return str[1 : len(str)-1], nil
			} else if tok == scanner.Int {
				v, err := strconv.ParseInt(p.s.TokenText(), 10, 64)
				if err != nil {
					return nil, err
				}
				return int64(mul) * v, err
			} else if tok == scanner.Float {
				v, err := strconv.ParseFloat(p.s.TokenText(), 64)
				if err != nil {
					return nil, err
				}
				return float64(mul) * v, err
			}
		default:
			return nil, p.parseError(fmt.Sprintf("invalid value: '%s'", p.s.TokenText()))
		}
	}
	return nil, p.parseError(fmt.Sprintf("invalid value: '%s'",
		string([]rune{p.s.Peek()})))
}

func (p *parser) parseString(_ rune) (string, error) {
	var buf bytes.Buffer
	ch := p.s.Next()
	for ch != '\'' {
		if ch == '\n' || ch == '\r' || ch < 0 {
			return "", p.parseError("unterminated string")
		}
		if ch == '\\' {
			s, err := p.parseEscape(ch)
			if err != nil {
				return "", err
			}
			buf.WriteString(s)
		} else {
			buf.WriteRune(ch)
		}
		ch = p.s.Next()
	}
	return buf.String(), nil
}

func (p *parser) parseEscape(_ rune) (string, error) {
	ch := p.s.Next()
	switch ch {
	case 'a':
		return "\a", nil
	case 'b':
		return "\b", nil
	case 'f':
		return "\f", nil
	case 'n':
		return "\n", nil
	case 'r':
		return "\r", nil
	case 't':
		return "\t", nil
	case 'v':
		return "\v", nil
	case '\\':
		return "\\", nil
	case '"':
		return "\"", nil
	case '\'':
		return "'", nil
	}
	return "", p.parseError(fmt.Sprintf("invalid escape sequence: %s", string(ch)))
}

func (p *parser) parseArray(_ rune) ([]interface{}, error) {
	result := []interface{}{}
	for {
		value, err := p.parseValue()
		if err != nil {
			return result, err
		}
		result = append(result, value)
		next := p.s.Next()
		if next == ']' {
			return result, nil
		}
		if next == ',' {
			continue
		}
		return result, p.parseError(fmt.Sprintf(", expected but got %s", string(next)))
	}
}

func (p *parser) parseArgs() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	for {
		tok := p.s.Scan()
		if tok != scanner.Ident {
			return result, p.parseError(fmt.Sprintf("invalid attribute name: %s", p.s.TokenText()))
		}
		name := p.s.TokenText()
		eq := p.s.Next()
		if eq != '=' {
			return result, p.parseError(fmt.Sprintf("= expected but got %s", string(eq)))
		}
		value, err := p.parseValue()
		if err != nil {
			return result, err
		}
		result[name] = value
		next := p.s.Next()
		if next == ')' {
			return result, nil
		}
		if next == ',' {
			continue
		}
		return result, p.parseError(fmt.Sprintf(") or , expected but got %s", string(next)))
	}
}

// ParseTag parses a given tag value.
func ParseTag(value string, name string) ([]Definition, error) {
	p := newParser(name)
	return p.Parse(value)
}

// ParseStruct parses struct tags of given object. map key is a field name.
func ParseStruct(obj interface{}, tag string) (map[string][]Definition, error) {
	result := map[string][]Definition{}
	r := reflect.ValueOf(obj)
	if r.Kind() == reflect.Ptr {
		obj = r.Elem().Interface()
	}
	rv := reflect.TypeOf(obj)
	for i := 0; i < rv.NumField(); i++ {
		f := rv.Field(i)
		value := f.Tag.Get(tag)
		if len(value) == 0 {
			continue
		}
		defs, err := ParseTag(value, rv.Name()+"."+f.Name)
		if err != nil {
			return nil, err
		}
		result[f.Name] = defs
	}
	return result, nil
}
