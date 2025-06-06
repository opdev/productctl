package certtargets_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCerttargets(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certtargets Suite")
}
