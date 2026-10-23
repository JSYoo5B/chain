package routing

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/JSYoo5B/chain"
)

const quotedDirection = "quoted\n\"label\""
const rawDirection = `raw "label"`

// Handwritten plans are an independent reference. Leaf plans must be explicit:
// NewWorkflow otherwise connects members sequentially by default.
func manualRoutes(start, success, failure, abort, quoted, raw chain.Action[int]) *chain.Workflow[int] {
	w := chain.NewWorkflow("Routes", start, success, failure, abort, quoted, raw)
	w.SetRunPlan(start, chain.RunPlan[int]{
		chain.Success: success, chain.Failure: failure, chain.Abort: abort,
		quotedDirection: quoted, rawDirection: raw,
	})
	for _, leaf := range []chain.Action[int]{success, failure, abort, quoted, raw} {
		w.SetRunPlan(leaf, chain.TerminationPlan[int]())
	}
	return w
}

func manualFallback(start, next chain.Action[int]) *chain.Workflow[int] {
	w := chain.NewWorkflow("Fallback", start, next)
	w.SetRunPlan(start, chain.RunPlan[int]{chain.Success: next})
	w.SetRunPlan(next, chain.TerminationPlan[int]())
	return w
}

type routeCase struct {
	name, direction          string
	runError, runPanic       error
	selectError, selectPanic error
	handlerError             error
	wantTrace                []string
	wantOutput               int
	wantErrors               []error
	wantMessage              string
}

type routeResult struct {
	output int
	trace  []string
	err    error
}

func runRoute(t *testing.T, tc routeCase, build func(...chain.Action[int]) *chain.Workflow[int]) routeResult {
	t.Helper()
	var trace []string
	start := chain.NewSimpleBranchAction("start", func(_ context.Context, input int) (int, error) {
		trace = append(trace, "start")
		if tc.runPanic != nil {
			panic(tc.runPanic)
		}
		return input + 1, tc.runError
	}, []string{quotedDirection, rawDirection, "unlisted"}, func(_ context.Context, _ int) (string, error) {
		trace = append(trace, "select")
		if tc.selectPanic != nil {
			panic(tc.selectPanic)
		}
		return tc.direction, tc.selectError
	})
	actions := []chain.Action[int]{start}
	for i, name := range []string{"success", "failure", "abort", "quoted", "raw"} {
		actions = append(actions, chain.NewSimpleAction(name, func(_ context.Context, input int) (int, error) {
			trace = append(trace, name)
			return input + (i+1)*100, tc.handlerError
		}))
	}
	w := build(actions...)
	if err := w.ValidateGraph(); err != nil {
		t.Fatal(err)
	}
	output, err := w.Run(context.Background(), 10)
	if output != tc.wantOutput || !reflect.DeepEqual(trace, tc.wantTrace) {
		t.Fatalf("got output %d, trace %q; want %d, %q", output, trace, tc.wantOutput, tc.wantTrace)
	}
	if len(tc.wantErrors) == 0 && tc.wantMessage == "" && err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, cause := range tc.wantErrors {
		if !errors.Is(err, cause) {
			t.Errorf("error %v does not preserve cause %v", err, cause)
		}
	}
	if tc.wantMessage != "" && (err == nil || !strings.Contains(err.Error(), tc.wantMessage)) {
		t.Errorf("error %v does not contain %q", err, tc.wantMessage)
	}
	return routeResult{output, trace, err}
}

func compareRoutes(t *testing.T, tc routeCase, manual, generated func(...chain.Action[int]) *chain.Workflow[int]) {
	t.Helper()
	want := runRoute(t, tc, manual)
	got := runRoute(t, tc, generated)
	message := func(err error) string {
		if err == nil {
			return ""
		}
		return err.Error()
	}
	if got.output != want.output || !reflect.DeepEqual(got.trace, want.trace) || message(got.err) != message(want.err) {
		t.Fatalf("generated result %+v differs from manual %+v", got, want)
	}
}

func TestRoutingMatchesManual(t *testing.T) {
	runFailure := errors.New("run failed")
	selectFailure := errors.New("selection failed")
	handlerFailure := errors.New("handler failed")
	panicCause := errors.New("panic cause")
	cases := []routeCase{
		{name: "success", direction: chain.Success, wantOutput: 111, wantTrace: []string{"start", "select", "success"}},
		{name: "failure aggregates errors", runError: runFailure, handlerError: handlerFailure, wantOutput: 211, wantTrace: []string{"start", "failure"}, wantErrors: []error{runFailure, handlerFailure}},
		{name: "panic uses abort", runPanic: panicCause, handlerError: handlerFailure, wantOutput: 310, wantTrace: []string{"start", "abort"}, wantErrors: []error{panicCause, handlerFailure}},
		{name: "quoted branch", direction: quotedDirection, wantOutput: 411, wantTrace: []string{"start", "select", "quoted"}},
		{name: "raw branch", direction: rawDirection, wantOutput: 511, wantTrace: []string{"start", "select", "raw"}},
		{name: "selector returns failure", direction: chain.Failure, selectError: selectFailure, wantOutput: 211, wantTrace: []string{"start", "select", "failure"}, wantErrors: []error{selectFailure}},
		{name: "selector returns branch with error", direction: quotedDirection, selectError: selectFailure, wantOutput: 411, wantTrace: []string{"start", "select", "quoted"}, wantErrors: []error{selectFailure}},
		{name: "selector panic uses abort", selectPanic: panicCause, wantOutput: 310, wantTrace: []string{"start", "select", "abort"}, wantErrors: []error{panicCause}},
		{name: "unlisted branch terminates", direction: "unlisted", wantOutput: 11, wantTrace: []string{"start", "select"}},
		{name: "unknown branch reports error", direction: "unknown", wantOutput: 11, wantTrace: []string{"start", "select"}, wantMessage: "no action plan from `start` directing `unknown`"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareRoutes(t, tc,
				func(a ...chain.Action[int]) *chain.Workflow[int] {
					return manualRoutes(a[0], a[1], a[2], a[3], a[4], a[5])
				},
				func(a ...chain.Action[int]) *chain.Workflow[int] {
					return BuildRoutes(a[0], a[1], a[2], a[3], a[4], a[5])
				})
		})
	}
}

func TestUnlistedDirectionsTerminate(t *testing.T) {
	runFailure := errors.New("run failed")
	panicCause := errors.New("panic cause")
	cases := []routeCase{
		{name: "success takes next", direction: chain.Success, wantOutput: 111, wantTrace: []string{"start", "select", "success"}},
		{name: "failure stops", runError: runFailure, wantOutput: 11, wantTrace: []string{"start"}, wantErrors: []error{runFailure}},
		{name: "abort stops", runPanic: panicCause, wantOutput: 10, wantTrace: []string{"start"}, wantErrors: []error{panicCause}},
		{name: "branch stops", direction: quotedDirection, wantOutput: 11, wantTrace: []string{"start", "select"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			compareRoutes(t, tc,
				func(a ...chain.Action[int]) *chain.Workflow[int] { return manualFallback(a[0], a[1]) },
				func(a ...chain.Action[int]) *chain.Workflow[int] { return BuildFallback(a[0], a[1]) })
		})
	}
}
