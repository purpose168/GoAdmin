package fasthttp

import (
	"net/http"
	"testing"

	"github.com/gavv/httpexpect"
	"github.com/purpose168/GoAdmin/tests/common"
)

func TestFasthttp(t *testing.T) {
	common.ExtraTest(httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewFastBinder(internalHandler()),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
	}))
}
