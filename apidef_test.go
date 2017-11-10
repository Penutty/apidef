package apidef

import (
	"bytes"
	"fmt"
	"testing"
)

var (
	tpath   = []byte("/path")
	tmethod = "method"
	name    = "name"
	email   = "email"
	sString = "string"
	validAN = "alphanumeric"
	validL  = "length"
	validE  = "email"
	lMin    = "8"
	lMax    = "64"

	namePass  = "James123"
	nameFail  = "user123"
	emailPass = "user@email.com"
	emailFail = "notanemail"
)

func Test_NewEndpoint(t *testing.T) {
	endpoint := NewEndpoint(tpath, tmethod)
	f := endpoint.NewField(name, sString).PassWith("James123").FailWith("user123")
	_ = f.NewValidField(validAN).NewValidField(validL, lMin, lMax)

	f = endpoint.NewField(email, sString).PassWith("user@email.com").FailWith("notanemail")
	_ = f.NewValidField(validE).NewValidField(validL, lMin, lMax)

	t.Run("Struct", func(t *testing.T) {
		sb := make([]byte, 0)
		structBuf := bytes.NewBuffer(sb)
		endpoint.Struct(structBuf)

		expected := fmt.Sprintf("type body struct {\n"+
			"\t%s %s `valid: \"%s,%s(%s|%s)\"`\n"+
			"\t%s %s `valid: \"%s,%s(%s|%s)\"`\n"+
			"}\n", name, sString, validAN, validL, lMin, lMax, email, sString, validE, validL, lMin, lMax)

		if structBuf.String() != expected {
			t.Log(structBuf)
			t.Log(expected)
			t.Fatal()
		}

	})

	t.Run("Tests", func(t *testing.T) {
		tb := make([]byte, 0)
		testsBuf := bytes.NewBuffer(tb)
		endpoint.Tests(testsBuf)

		expected := fmt.Sprintf("type %s struct {\n"+
			"\treq *http.Request\n"+
			"\tpassing bool\n"+
			"}\n"+
			"tests := []*%s{\n", endpoint.testType(), endpoint.testType())

		vars := [][]string{
			[]string{namePass, emailPass, "true"},
			[]string{namePass, emailFail, "false"},
			[]string{nameFail, emailPass, "false"},
			[]string{nameFail, emailFail, "false"},
		}

		for _, v := range vars {
			expected += fmt.Sprintf("\t&%s{\n"+
				"\t\thttptest.NewRequest(%s, \"%s\",\n"+
				"\t\t\tstrings.NewReader(`{\n"+
				"\t\t\t\t\"%s\": \"%s\",\n"+
				"\t\t\t\t\"%s\": \"%s\"\n"+
				"\t\t\t}`)),\n"+
				"\t\t%s,\n"+
				"\t},\n", endpoint.testType(), tmethod, tpath, name, v[0], email, v[1], v[2])
		}

		expected += fmt.Sprintf("}\n")
		if testsBuf.String() != expected {
			t.Log(testsBuf)
			t.Log(expected)
			t.Fatal()
		}

	})
}
