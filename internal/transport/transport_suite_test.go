package transport_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// fakeHTTPEndpoint is just used to test transport RoudnTrips.
// No server is set up here, so it's expected to ECONNREFUSED
const fakeHTTPEndpoint = "http://localhost:31425"

func TestTransport(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transport Suite")
}
