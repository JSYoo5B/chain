package chain

import "fmt"

type node[T any] = Action[T]
type runPlanGraph[T any] = map[node[T]]RunPlan[T]
type connectivityGraph[T any] = map[node[T]][]node[T]
type cycleVisitMap[T any] = map[node[T]]int
type connectedNodeMap[T any] = map[node[T]]bool
type edgeDirection = string
type cycleTrace = []string

// ValidateGraph ensures the workflow's graph is connected and acyclic.
// It checks for cycles first, then verifies that all nodes are connected as a single graph.
func (w *Workflow[T]) ValidateGraph() error {
	runPlans := runPlanGraph[T](w.runPlans)

	// Step 1: Check every run plan graph component for cycles.
	// Cycle validation must use edge direction and must scan every node because
	// a disconnected subgraph can still contain a cycle.
	visited := make(cycleVisitMap[T])
	for currentNode := range runPlans {
		if visited[currentNode] != notVisited {
			continue
		}
		if err := dfsWithCycleCheck(currentNode, runPlans, visited, cycleTrace{}); err != nil {
			return err
		}
	}

	// Step 2: Check that all nodes form one graph.
	// At this point, directed cycle validation is already done. For connectivity,
	// we only need weak connectivity: whether every node belongs to the same
	// connected component when edge directions and direction labels are ignored.
	connected := make(connectedNodeMap[T], len(runPlans))
	connectivity := buildConnectivityGraph(runPlans)

	// initAction is one of the workflow nodes. In a single connected graph,
	// starting from any node should be enough to visit every other node.
	traverseConnectedNodes(w.initAction, connectivity, connected)
	for currentNode := range runPlans {
		// A node left unvisited here belongs to a different connected component,
		// so the workflow is not a single graph.
		if !connected[currentNode] {
			return fmt.Errorf("disconnected graph detected: action `%s` is not connected to initAction `%s`", currentNode.Name(), w.initAction.Name())
		}
	}

	return nil
}

func dfsWithCycleCheck[T any](currentNode node[T], runPlans runPlanGraph[T], visited cycleVisitMap[T], path cycleTrace) error {
	path = append(path, "`"+currentNode.Name()+"`")

	// This is a 3-color DFS:
	// - notVisited: this node has not been inspected yet.
	// - visiting: this node is in the current DFS stack.
	// - confirmed: every downstream path from this node has already been checked.
	// Reaching a visiting node means the current path points back to one of its
	// ancestors, so the directed graph has a cycle. Reaching a confirmed node is
	// safe because it is just a previously verified path, which is valid in DAGs
	// where branches merge back into the same node.
	if visited[currentNode] == visiting {
		return fmt.Errorf("cycle detected: %v", path)
	} else if visited[currentNode] == confirmed {
		return nil
	}

	visited[currentNode] = visiting

	terminate := Terminate[T]()
	for direction, nextNode := range runPlans[currentNode] {
		if nextNode != terminate {
			edge := formatEdge(direction)
			path = append(path, edge)
			if err := dfsWithCycleCheck(nextNode, runPlans, visited, path); err != nil {
				return err
			}
			path = path[:len(path)-1]
		}
	}

	visited[currentNode] = confirmed

	return nil
}

func formatEdge(direction edgeDirection) string {
	return "-" + direction + "->"
}

func buildConnectivityGraph[T any](runPlans runPlanGraph[T]) connectivityGraph[T] {
	connectivity := make(connectivityGraph[T], len(runPlans))
	terminate := Terminate[T]()

	for currentNode, plan := range runPlans {
		if _, exists := connectivity[currentNode]; !exists {
			connectivity[currentNode] = nil
		}

		// Connectivity only needs to know whether two nodes are connected, so
		// `current -> next` is expanded as `current <-> next`. Direction labels
		// such as success/failure/abort are intentionally dropped.
		for _, nextNode := range plan {
			if nextNode == terminate {
				continue
			}

			connectivity[currentNode] = append(connectivity[currentNode], nextNode)
			connectivity[nextNode] = append(connectivity[nextNode], currentNode)
		}
	}

	return connectivity
}

func traverseConnectedNodes[T any](currentNode node[T], connectivity connectivityGraph[T], visited connectedNodeMap[T]) {
	if visited[currentNode] {
		return
	}

	// Walk the derived connectivity graph by expanding through both original
	// incoming and outgoing edges, because buildConnectivityGraph made every
	// non-termination edge bidirectional.
	visited[currentNode] = true
	for _, nextNode := range connectivity[currentNode] {
		traverseConnectedNodes(nextNode, connectivity, visited)
	}
}

const (
	notVisited = iota
	visiting
	confirmed
)
