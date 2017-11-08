package apidef

import (
	"testing"
)

var (
	tpath   = []byte("/path")
	tmethod = "method"
)

func Test_NewEndpoint(t *testing.T) {
	endpoint := NewEndpoint(tpath, tmethod)
	f := endpoint.NewField("name", "string").PassWith("James123", "user333598").FailWith("user123", "123", "")
	_ = f.NewValidField("alphanumeric").NewValidField("length", "8", "64")

	f = endpoint.NewField("email", "string").PassWith("user@email.com", "jamees@gmail.com").FailWith("notanemail", "1@m.com")
	_ = f.NewValidField("email").NewValidField("length", "8", "64")

	f = endpoint.NewField("password", "string").PassWith("abc123!!??").FailWith("user", "", "12345678")
	_ = f.NewValidField("alphanumeric").NewValidField("length", "8", "64")

	f = endpoint.NewField("something", "int").PassWith("23").FailWith("aaa")
	_ = f.NewValidField("numeric")

	endpoint.Struct()

	endpoint.Tests()

}

// field 1
// pass
// pass2
// fail
// fail2
// fail3

// field 2
// pass
// fail
