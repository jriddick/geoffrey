package irc_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestIrc(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Irc Suite")
}
