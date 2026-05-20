// Package chain provides small, composable workflow primitives.
//
// The basic unit is Action. A Workflow connects actions with RunPlan values and
// moves between them through directions such as Success, Failure, Abort, or
// custom directions returned by BranchAction implementations.
//
// The package also includes wrappers for common execution patterns:
//
//   - AsBranchAction and NewSimpleBranchAction for direction-based branching.
//   - AsRetryableAction for bounded retries and optional rollback.
//   - AsBestEffortAction for non-critical work that should not fail a workflow.
//   - AsSequenceSliceAction and AsSequenceMapAction for sequential collection
//     processing.
//   - AsParallelSliceAction and AsParallelMapAction for parallel collection
//     processing.
//   - AdaptAction for running an action against a field or sub-value inside a
//     larger value.
package chain
