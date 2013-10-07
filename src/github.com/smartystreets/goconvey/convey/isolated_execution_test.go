package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
	"testing"
)

func TestSingleScope(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "done"
	})

	expectEqual(t, "done", output)
}

func TestSingleScopeWithMultipleConveys(t *testing.T) {
	output := prepare()

	Convey("1", t, func() {
		output += "1"
	})

	Convey("2", t, func() {
		output += "2"
	})

	expectEqual(t, "12", output)
}

func TestNestedScopes(t *testing.T) {
	output := prepare()

	Convey("a", t, func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})
		})
	})

	expectEqual(t, "a aa aaa | ", output)
}

func TestNestedScopesWithIsolatedExecution(t *testing.T) {
	output := prepare()

	Convey("a", t, func() {
		output += "a "

		Convey("aa", func() {
			output += "aa "

			Convey("aaa", func() {
				output += "aaa | "
			})

			Convey("aaa1", func() {
				output += "aaa1 | "
			})
		})

		Convey("ab", func() {
			output += "ab "

			Convey("abb", func() {
				output += "abb | "
			})
		})
	})

	expectEqual(t, "a aa aaa | a aa aaa1 | a ab abb | ", output)
}

func TestSingleScopeWithConveyAndNestedReset(t *testing.T) {
	output := prepare()

	Convey("1", t, func() {
		output += "1"

		Reset(func() {
			output += "a"
		})
	})

	expectEqual(t, "1a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndReset(t *testing.T) {
	output := prepare()

	Convey("reset after each nested convey", t, func() {
		Convey("first output", func() {
			output += "1"
		})

		Convey("second output", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})
	})

	expectEqual(t, "1a2a", output)
}

func TestSingleScopeWithMultipleRegistrationsAndMultipleResets(t *testing.T) {
	output := prepare()

	Convey("each reset is run at end of each nested convey", t, func() {
		Convey("1", func() {
			output += "1"
		})

		Convey("2", func() {
			output += "2"
		})

		Reset(func() {
			output += "a"
		})

		Reset(func() {
			output += "b"
		})
	})

	expectEqual(t, "1ab2ab", output)
}

func TestPanicAtHigherLevelScopePreventsChildScopesFromRunning(t *testing.T) {
	output := prepare()

	Convey("This step panics", t, func() {
		Convey("this should NOT be executed", func() {
			output += "1"
		})

		panic("Hi")
	})

	expectEqual(t, "", output)
}

func TestPanicInChildScopeDoes_NOT_PreventExecutionOfSiblingScopes(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling should execute", func() {
			output += "2"
		})
	})

	expectEqual(t, "2", output)
}

func TestResetsAreAlwaysExecutedAfterScopePanics(t *testing.T) {
	output := prepare()

	Convey("This is the parent", t, func() {
		Convey("This step panics", func() {
			panic("Hi")
			output += "1"
		})

		Convey("This sibling step does not panic", func() {
			output += "a"

			Reset(func() {
				output += "b"
			})
		})

		Reset(func() {
			output += "2"
		})
	})

	expectEqual(t, "2ab2", output)
}

func TestSkipTopLevel(t *testing.T) {
	output := prepare()

	SkipConvey("hi", t, func() {
		output += "This shouldn't be executed!"
	})

	expectEqual(t, "", output)
}

func TestSkipNestedLevel(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "yes"

		SkipConvey("bye", func() {
			output += "no"
		})
	})

	expectEqual(t, "yes", output)
}

func TestSkipNestedLevelSkipsAllChildLevels(t *testing.T) {
	output := prepare()

	Convey("hi", t, func() {
		output += "yes"

		SkipConvey("bye", func() {
			output += "no"

			Convey("byebye", func() {
				output += "no-no"
			})
		})
	})

	expectEqual(t, "yes", output)
}

func prepare() string {
	runner = execution.NewRunner()
	reporting.QuietMode()
	return ""
}
