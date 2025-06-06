package certifycontainers_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCertifycontainers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certifycontainers Suite")
}
