package apidef

import (
	"testing"
)

var (
	tpath       = []byte("testpath")
	tmethod     = []byte("testmethod")
	tfield      = []byte("testfield")
	tvalidField = []byte("testValidField")
	tType       = []byte("string")
)

func Test_NewEndpoint(t *testing.T) {
	endpoint := NewEndpoint(tpath, tmethod)
	_ = endpoint.NewField(tfield, tType).NewValidField(tvalidField)
	//	expected := "type body struct {\n" +
	//		"\ttestfield string `valid: \"alpha\"\n" +
	//		"}\n"
	endpoint.Struct()
	endpoint.Tests()
}
