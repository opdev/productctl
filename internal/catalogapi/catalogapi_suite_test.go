package catalogapi_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCatalogapi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Catalogapi Suite")
}
