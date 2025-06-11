package jsonschema_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLSP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JSONschema Suite")
}
