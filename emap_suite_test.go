// Copyright(c) 2016 Ethan Zhuang <zhuangwj@gmail.com>.

package emap

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestProxy(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "EMap Suite")
}

var _ = Describe("Test initialization", func() {
	Context("Register suite setup and teardown function", func() {
		BeforeSuite(func() {
		})

		AfterSuite(func() {

		})
	})
})
