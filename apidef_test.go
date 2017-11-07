package apidef

import (
	"testing"
)

var (
	tpath        = []byte("/path")
	tmethod      = "method"
	tfield       = "field"
	tfield2      = "field2"
	tvalidField  = "validator"
	tvalidField2 = "validator2"
	tType        = "string"
	pVal         = "pass"
	pVal2        = "pass2"
	fVal         = "fail"
	fVal2        = "fail2"
	fVal3        = "fail3"
)

func Test_NewEndpoint(t *testing.T) {
	endpoint := NewEndpoint(tpath, tmethod)
	_ = endpoint.NewField(tfield, tType).NewValidField(tvalidField).PassWith(pVal, pVal2).FailWith(fVal, fVal2, fVal3)
	_ = endpoint.NewField(tfield2, tType).NewValidField(tvalidField2).PassWith(pVal).FailWith(fVal)

	//	expected := "type body struct {\n" +
	//		"\ttestfield string `valid: \"alpha\"\n" +
	//		"}\n"
	endpoint.Struct()

	endpoint.Test("nil", "james123", "useremail@email.com", "testspass")
	endpoint.Tests()

	endpoint.PassingTests()
}
