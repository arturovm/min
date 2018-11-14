package min_test

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

// TestMin is the entry point for tests
func TestMin(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Min Suite")
}

// These variables define our own BDD semantics
var (
	Describe = ginkgo.Describe
	Given    = ginkgo.Describe
	When     = ginkgo.Context
	Then     = ginkgo.It
)

// These variables bring Gomega methods into our scope
var (
	Expect       = gomega.Expect
	HaveOccurred = gomega.HaveOccurred
	Equal        = gomega.Equal
)
