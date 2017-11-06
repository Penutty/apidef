// Licensing will go here.

// apidef.go is designed to accept api defintions.
// Use these definitions to generate unit-tests and structs for input validation.
// Designed to be used with Go Generate.
package apidef

import (
	"errors"
	"fmt"
	valid "github.com/asaskevich/govalidator"
	"strconv"
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
)

type endPointKey struct {
	path   []byte
	method []byte
}

type endPoint struct {
	endPointKey
	fields []*field
}

func NewEndpoint(path []byte, method []byte) *endPoint {
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
	testType := string(append(strings.ToUpper(e.method), e.path[1:]...)) + "Test"

	fieldErrors = make([]string, len(e.fields))
	for i, f := range e.fields {
		fieldErrors[i] = f.Errors()
	}

	fieldTests = make([]string, len(e.fields))
	for _, f := range e.fields {
		fieldTests[i] = f.Tests()
	}

	fmt.Printf(`
var (
%s
)

type %s struct {
	req *http.Request
 	err error
}

tests := []*%s {
%s
}
`, strings.Join(fieldErrors, "\n"), testType, testType, strings.Join(fieldTests))
}

func (e *endPoint) Struct() {
	fmt.Printf("type body struct {\n")
	for _, f := range e.fields {
		fmt.Printf("%s\n", f)
	}
	fmt.Printf("}\n")

}

type field struct {
	name       []byte
	Type       []byte
	validators []*validField
}

func (e *endPoint) NewField(name []byte, Type []byte) *field {
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

func (f *field) Errors() string {
	validatorsErrors := make([]string, len(f))
	for i, v := range f.validators {
		validatorsErrors[i] = v.Errors()
	}
	return strings.Join(validatorsErrors, "\n")
}

func (f *field) Tests() string {
	validatorsTests := make([]string, len(f))
	for i, v := range f.validators {
		validatorsTests[i] = v.Tests()
	}
	return Sprintf(`
	&postUserTest{
		http.test.NewRequest(http.MethodPost, "/user", 
			strings.NewReader("{
				%s
			}")),
		ErrorApiDef%s,
	},
	)`, strings.Join(validatorsTests, ",\n"), f.name)
}

type validField struct {
	name     []byte
	min      uint64
	max      uint64
	passVals []string
	failVals []string
}

func (v *validField) String() string {
	if v.hasParams() {
		return fmt.Sprintf("%s(%v|%v)", v.name, strconv.Itoa(int(v.min)), strconv.Itoa(int(v.max)))
	} else {
		return fmt.Sprintf("%s", v.name)
	}
}

func (v *validField) Errors() string {
	return fmt.Sprintf("\tErrorApiDef%s = errors.New(\"Invalid field. Criteria \"%s\" not met.\")", string(v.name), string(v.name))
}

func (v *validField) Tests() string {
	if len(def) <= 0 {
		panic(ErrorInvalidDefault)
	}

	if v.hasParams() {

	} else {

	}
}

func (v *validField) hasParams() bool {
	if v.min == 0 &&
		v.max == 0 &&
		len(v.passVal) <= 0 &&
		len(v.failVal) == 0 {
		return false
	}
	return true
}

func (v *validField) PassWith(val string) *validField {
	v.passVals = append(v.passVals, val)
	return v
}

func (v *validField) FailWith(val string) *validField {
	v.failVals = append(v.failVals, val)
	return v
}

func (f *field) NewValidField(name []byte) *validField {
	switch {
	case len(name) <= 0:
		panic(ErrorInvalidName)
	}

	vf := &validField{
		name: name,
	}

	f.validators = append(f.validators, vf)
	return vf
}

func (f *field) NewValidFieldParams(name []byte, min uint64, max uint64) *validField {
	if len(name) <= 0 {
		panic(ErrorInvalidName)
	}

	vf := &validField{
		name: name,
		min:  min,
		max:  max,
	}
	f.validators = append(f.validators, vf)
	return vf
}
