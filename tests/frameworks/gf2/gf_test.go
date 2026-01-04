package gf

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/purpose168/GoAdmin/tests/common"
)

func TestGf2(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(internalHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}
