package file_test

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/pivotalservices/file-downloader-resource/file"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMatcher(t *testing.T) {
	spec.Run(t, "Matcher", testMatcher, spec.Report(report.Terminal{}))
}

func testMatcher(t *testing.T, when spec.G, it spec.S) {

	it.Before(func() {
		RegisterTestingT(t)

	})
	it.After(func() {

	})
	when("matches file path", func() {
		it("returns true", func() {
			match, err := file.Matches("elastic-runtime/cf-2.3.3-build.10.pivotal", "elastic-runtime", "cf*.pivotal", "2.3.3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(match).Should(BeTrue())
		})
	})
	when("doesn't match file path", func() {
		it("returns false", func() {
			match, err := file.Matches("elastic-runtime/cf-2.3.3-build.10.pivotal", "elastic-runtime", "cf*.pivotal", "2.3.4")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(match).Should(BeFalse())
		})
		it("returns false", func() {
			match, err := file.Matches("elastic-runtime/cf-2.3.3-build.10.pivotal", "elastic-runtimer", "cf*.pivotal", "2.3.3")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(match).Should(BeFalse())
		})
	})
}
