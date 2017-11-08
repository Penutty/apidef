// Licensing will go here.

// apidef.go is designed to accept api defintions.
// Use these definitions to generate unit-tests and structs for input validation.
// Designed to be used with Go Generate.
package apidef

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrorInvalidEndpointParameters = errors.New("A path and method must be passed via the -path and -method flags")
	ErrorInvalidGenerateParameters = errors.New("The flag genTests OR genStructs must be set. Both may not be set.")
	ErrorApiEndPointDNE            = errors.New("The API endpoint as specified by -path and -method is not defined in main.")
	ErrorInvalidPath               = errors.New("path must be a valid string.")
	ErrorInvalidMethod             = errors.New("method must be a valid string.")
	ErrorInvalidName               = errors.New("name must be a valid string.")
	ErrorInvalidType               = errors.New("type must be a vaild string.")
	ErrorEmptyString               = errors.New("string value must not be empty.")
)

type endPointKey struct {
	path   []byte
	method string
}

type endPoint struct {
	endPointKey
	fields []*field
}

func NewEndpoint(path []byte, method string) *endPoint {
	switch {
	case len(path) <= 0:
		panic(ErrorInvalidPath)
	case len(method) <= 0:
		panic(ErrorInvalidMethod)
	}

	return &endPoint{
		endPointKey: endPointKey{
			path:   path,
			method: method,
		},
	}
}

func (e *endPoint) Tests() {
	fieldTvals := make([][]*testVal, len(e.fields))
	for i, f := range e.fields {
		fieldTvals[i] = f.testVals
	}

	fmt.Printf("tests := []*%s{\n", e.testType())
	e.testCombs(fieldTvals, make([]*testVal, 0))
	fmt.Printf("}\n")
}

func (e *endPoint) test(ts []*testVal) {
	passing := "true"
	ss := make([]string, len(ts))
	for i, t := range ts {
		ss[i] = fmt.Sprintf("\t\t\t\t\"%s\": \"%s\"", t.field, t.val)
		if t.passing == false {
			passing = "false"
		}
	}
	fmt.Printf(
		"\t&%s{\n"+
			"\t\thttptest.NewRequest(%s, \"%s\",\n"+
			"\t\t\tstrings.NewReader(`{\n"+
			"%s\n"+
			"\t\t\t}`)),\n"+
			"\t\t%s,\n"+
			"\t},\n", e.testType(), e.method, e.path, strings.Join(ss, ",\n"), passing)
}

func (e *endPoint) testCombs(m [][]*testVal, ts []*testVal) {
	if len(m) == 1 {
		for _, v := range m[0] {
			e.test(append(ts, v))
		}
		return
	}
	for _, v := range m[0] {
		e.testCombs(m[1:], append(ts, v))
	}
}

func (e *endPoint) testType() string {
	return strings.ToUpper(e.method) + string(e.path[1:]) + "Test"
}

func (e *endPoint) Struct() {
	fmt.Printf("type body struct {\n")
	for _, f := range e.fields {
		fmt.Printf("%s\n", f)
	}
	fmt.Printf("}\n")
}

type testVal struct {
	field   string
	val     string
	passing bool
}

type field struct {
	name       string
	Type       string
	testVals   []*testVal
	validators []*validField
}

func (e *endPoint) NewField(name string, Type string) *field {
	switch {
	case len(name) <= 0:
		panic(ErrorInvalidName)
	case len(Type) <= 0:
		panic(ErrorInvalidType)
	}

	f := &field{
		name: name,
		Type: Type,
	}
	e.fields = append(e.fields, f)
	return f
}

func (f *field) String() string {
	vs := make([]string, len(f.validators))
	for i, v := range f.validators {
		vs[i] = fmt.Sprintf("%s", v)
	}
	return fmt.Sprintf("\t%s %s `valid: \"%s\"`", f.name, f.Type, strings.Join(vs, ","))
}

func (f *field) PassWith(s ...string) *field {
	if len(s) == 0 {
		panic(ErrorEmptyString)
	}
	for _, v := range s {
		f.testVals = append(f.testVals, &testVal{field: f.name, val: v, passing: true})
	}
	return f
}

func (f *field) FailWith(s ...string) *field {
	if len(s) == 0 {
		panic(ErrorEmptyString)
	}
	for _, v := range s {
		f.testVals = append(f.testVals, &testVal{field: f.name, val: v, passing: false})
	}
	return f
}

type validField struct {
	name   string
	params []string
}

func (v *validField) String() string {
	if v.hasParams() {
		return fmt.Sprintf("%s(%s)", v.name, strings.Join(v.params, "|"))
	} else {
		return fmt.Sprintf("%s", v.name)
	}
}

func (v *validField) hasParams() bool {
	if len(v.params) == 0 {
		return false
	}
	return true
}

func (f *field) NewValidField(name string, params ...string) *field {
	switch {
	case len(name) <= 0:
		panic(ErrorInvalidName)
	}

	vf := new(validField)
	if len(params) == 0 {
		vf = &validField{
			name: name,
		}
	} else {
		vf = &validField{
			name:   name,
			params: params,
		}
	}

	f.validators = append(f.validators, vf)
	return f

}
