package certifyhelmcharts_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCertifyhelmcharts(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Certifyhelmcharts Suite")
}
