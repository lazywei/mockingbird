package mockingbird_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMockingbird(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mockingbird Suite")
}
