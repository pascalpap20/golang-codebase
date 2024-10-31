package logging

import (
	"runtime"
)

type StackFrame struct {
	Func string `json:"func"`
	File string `json:"file"`
	Line int    `json:"line"`
}

func GetTrace(skip int) []StackFrame {
	stacktraces := []StackFrame{}

	// Ask runtime.Callers for up to 10 PCs, including runtime.Callers itself.
	pc := make([]uintptr, 20)
	n := runtime.Callers(2+skip, pc)
	if n == 0 {
		// No PCs available. This can happen if the first argument to
		// runtime.Callers is large.
		//
		// Return now to avoid processing the zero Frame that would
		// otherwise be returned by frames.Next below.
		return stacktraces
	}

	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)

	// Loop to get frames.
	// A fixed number of PCs can expand to an indefinite number of Frames.
	for {
		frame, more := frames.Next()

		// Process this frame.
		//
		// To keep this example's output stable
		// even if there are changes in the testing package,
		// stop unwinding when we leave package runtime.
		//if !strings.Contains(frame.File, "runtime/") {
		//	break
		//}

		stacktraces = append(stacktraces, StackFrame{
			frame.Function, frame.File, frame.Line,
		})

		// fmt.Printf("- more:%v | %s\n", more, frame.Function)

		// Check whether there are more frames to process after this one.
		if !more {
			break
		}
	}

	return stacktraces
}
