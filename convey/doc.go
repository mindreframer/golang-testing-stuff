// Package convey contains all of the public-facing entry points to this project.
// This means that it should never be required of the user to import any other
// packages from this project as they serve internal purposes.
package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
)

// Convey is the method intended for use when declaring the scopes
// of a specification. Each scope has a description and a func()
// which may contain other calls to Convey(), Reset() or Should-style
// assertions. Convey calls can be nested as far as you see fit.
//
// IMPORTANT NOTE: The top-level Convey() within a Test method
// must conform to the following signature:
//
//     Convey(description string, t *testing.T, action func())
//
// All other calls should like like this (no need to pass in *testing.T):
//
//     Convey(description string, action func())
//
// Don't worry, the goconvey will panic if you get it wrong so you can fix it.
//
// See the examples package for, well, examples.
func Convey(items ...interface{}) {
	entry := discover(items)
	register(entry)
}

// SkipConvey is analagous to Convey except that the scope is not executed
// (which means that child scopes defined within this scope are not run either).
// The reporter will be notified that this step was skipped.
func SkipConvey(items ...interface{}) {
	entry := discover(items)
	entry.Action = execution.NewSkippedAction(skipReport)
	register(entry)
}

func register(entry *execution.Registration) {
	if entry.Test != nil {
		runner.Begin(entry)
		runner.Run()
	} else {
		runner.Register(entry)
	}
}

func skipReport() {
	reporter.Report(reporting.NewSkipReport())
}

// Reset registers a cleanup function to be run after each Convey()
// in the same scope. See the examples package for a simple use case.
func Reset(action func()) {
	runner.RegisterReset(execution.NewAction(action))
}

// So is the means by which assertions are made against the system under test.
// The majority of exported names in this package begin with the word 'Should'
// and describe how the first argument (actual) should compare with any of the
// final (expected) arguments. How many final arguments are accepted depends on
// the particular assertion that is passed in as the assert argument.
// See the examples package for use cases and the assertions package for
// documentation on specific assertion methods.
func So(actual interface{}, assert assertion, expected ...interface{}) {
	if result := assert(actual, expected...); result == assertionSuccess {
		reporter.Report(reporting.NewSuccessReport())
	} else {
		reporter.Report(reporting.NewFailureReport(result))
	}
}

// SkipSo is analagous to So except that the assertion that would have been passed
// to So is not executed and the reporter is notified that the assertion was skipped.
func SkipSo(stuff ...interface{}) {
	skipReport()
}

// assertion is an alias for a function with a signature that the convey.So()
// method can handle. Any future or custom assertions should conform to this
// method signature. The return value should be an empty string if the assertion
// passes and a well-formed failure message if not.
type assertion func(actual interface{}, expected ...interface{}) string

const assertionSuccess = ""
