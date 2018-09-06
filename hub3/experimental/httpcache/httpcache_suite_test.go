package httpcache_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestHttpcache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Httpcache Suite")
}
