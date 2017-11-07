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
	tests  []string
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

func (e *endPoint) Test(err string, testVals ...string) {
	var vals []string
	if len(testVals) != len(e.fields) {
		vals = make([]string, len(e.fields))
		for i := 0; i < len(vals); i++ {
			vals[i] = testVals[i]
		}
	}
	testBody := make([]string, len(e.fields))
	for i, f := range e.fields {
		testBody[i] = fmt.Sprintf("\t\t\t\t\"%s\": \"%s\"", f.name, vals[i])
	}

	e.tests = append(e.tests, fmt.Sprintf(
		"\t&%s{\n"+
			"\t\thttptest.NewRequest(%s, \"%s\",\n"+
			"\t\t\tstrings.NewReader(`{\n"+
			"%s\n"+
			"\t\t\t}`)),\n"+
			"\t\t%s,\n"+
			"\t},\n", e.testType(), e.method, e.path, strings.Join(testBody, ",\n"), err))
}

func (e *endPoint) Tests() {
	fmt.Printf("tests := []*%s{\n"+
		"%s"+
		"}\n", e.testType(), strings.Join(e.tests, "\n"))
}

func (e *endPoint) PassingTests() {
	type set struct {
		val     string
		passing bool
	}

	m := make(map[string][]*set)

	for _, f := range e.fields {
		for _, v := range f.validators {
			for _, p := range v.pVals {
				s := &set{
					val:     p,
					passing: true,
				}
				m[f.name] = append(m[f.name], s)
			}
			for _, n := range v.fVals {
				s := &set{
					val:     n,
					passing: true,
				}
				m[f.name] = append(m[f.name], s)
			}
		}
	}

	fmt.Printf("%v\n", m)
	fmt.Printf("\n\n")

	for i, ss := range m {
		for j, s := range ss {
			fmt.Printf("\"%s\": \"%s\"\n", i, s.val)
			for l, ss2 := range m {
				if i == l {
					continue
				}
				for k, s2 := range ss2 {
					if j == k {
						continue
					}
					fmt.Printf("\"%s\": \"%s\"\n", l, s2.val)
				}
			}
			fmt.Printf("\n\n")
		}

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

type field struct {
	name       string
	Type       string
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

type validField struct {
	name   string
	params []string
	pVals  []string
	fVals  []string
}

func (v *validField) PassWith(s ...string) *validField {
	if len(s) == 0 {
		panic(ErrorEmptyString)
	}

	v.pVals = append(v.pVals, s...)
	return v
}

func (v *validField) FailWith(s ...string) *validField {
	if len(s) == 0 {
		panic(ErrorEmptyString)
	}

	v.fVals = append(v.fVals, s...)
	return v
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

func (f *field) NewValidField(name string, params ...string) *validField {
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
	return vf

}
