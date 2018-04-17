package alltest_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/onsi/ginkgo/reporters"
)

func TestAlltest(t *testing.T) {
	RegisterFailHandler(Fail)
	rs := []Reporter{}
	rs = append(rs, reporters.NewJUnitReporter("alltest.xml"))
	RunSpecsWithDefaultAndCustomReporters(t, "AllTestSuite", rs)
}
