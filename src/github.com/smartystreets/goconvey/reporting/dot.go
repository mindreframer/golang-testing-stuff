package reporting

import "fmt"

import (
	"github.com/smartystreets/goconvey/printing"
)

func (self *dot) BeginStory(story *StoryReport) {}

func (self *dot) Enter(scope *ScopeReport) {}

func (self *dot) Report(report *AssertionReport) {
	if report.Error != nil {
		fmt.Print(redColor)
		self.out.Insert(dotError)
	} else if report.Failure != "" {
		fmt.Print(yellowColor)
		self.out.Insert(dotFailure)
	} else if report.Skipped {
		fmt.Print(yellowColor)
		self.out.Insert(dotSkip)
	} else {
		fmt.Print(greenColor)
		self.out.Insert(dotSuccess)
	}
	fmt.Print(resetColor)
}

func (self *dot) Exit() {}

func (self *dot) EndStory() {}

func NewDotReporter(out *printing.Printer) *dot {
	self := dot{}
	self.out = out
	return &self
}

type dot struct {
	out *printing.Printer
}
