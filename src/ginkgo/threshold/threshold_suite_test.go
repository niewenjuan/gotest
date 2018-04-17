package threshold_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/onsi/ginkgo/reporters"
)

func TestThreshold(t *testing.T) {
	RegisterFailHandler(Fail)
	rs := []Reporter{}
	rs = append(rs, reporters.NewJUnitReporter("threshold.xml"))
	RunSpecsWithDefaultAndCustomReporters(t, "ThresholdTestSuite", rs)
}
