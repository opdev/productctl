package archivecomponent_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDeleteProductListing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ArchiveComponent Suite")
}
