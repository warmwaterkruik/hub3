package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestRapidSaas(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "RapidSaas Suite")
}
