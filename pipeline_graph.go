package railway

import "fmt"

// ValidateGraph ensures the pipeline's graph is connected and acyclic.
// It checks for cycles first, then verifies that all nodes connected as a single graph.
func (p *Pipeline[T]) ValidateGraph() error {
	// Step 1: Perform DFS from initAction to check for cycles and track visited nodes
	visited := make(map[Action[T]]int)
	if err := dfsWithCycleCheck(p.initAction, p.runPlans, visited, []string{}); err != nil {
		return err
	}

	// Step 2: After DFS, check if all actions have been visited
	unvisited := make([]Action[T], 0, len(p.runPlans))
	for action := range p.runPlans {
		if visited[action] == notVisited {
			unvisited = append(unvisited, action)
		}
	}

	// Step 3: If there are unvisited nodes, start new DFS from them
	for len(unvisited) > 0 {
		newStart := unvisited[0] // Pick any unvisited node
		visitedFromNewStart := make(map[Action[T]]int)
		if err := dfsWithCycleCheck(newStart, p.runPlans, visitedFromNewStart, []string{}); err != nil {
			return err
		}

		// Step 4: Merge the visited nodes from the current traversal into the overall visited set
		// If the current traversal's visited nodes intersect with the previously visited ones, they are connected
		// If there is no intersection, it's a disconnected graph.
		intersectionFound := false
		for action := range visitedFromNewStart {
			if visited[action] != notVisited {
				intersectionFound = true
			}
			visited[action] = confirmed
		}

		// Step 5: If no intersection found, it means the graph is disconnected
		if !intersectionFound {
			return fmt.Errorf("disconnect detected: action `%s` cannot reach the graph started from initAction `%s`", newStart.Name(), p.initAction.Name())
		}

		// Step 6: Check all nodes have been visited, no need for further checks
		if len(visited) == len(p.runPlans) {
			return nil
		}

		// If there are still unvisited nodes, continue to the next round of DFS
		stillUnvisited := make([]Action[T], 0, len(unvisited))
		for _, action := range unvisited {
			if visited[action] == notVisited {
				stillUnvisited = append(stillUnvisited, action)
			}
		}
		unvisited = stillUnvisited
	}

	return nil
}

func dfsWithCycleCheck[T any](node Action[T], graph map[Action[T]]ActionPlan[T], visited map[Action[T]]int, path []string) error {
	path = append(path, "`"+node.Name()+"`")

	if visited[node] != notVisited {
		return fmt.Errorf("cycle detected: %v", path)
	}

	visited[node] = visiting

	terminate := Terminate[T]()
	for direction, nextAction := range graph[node] {
		if nextAction != terminate {
			edge := "-" + direction + "->"
			path = append(path, edge)
			if err := dfsWithCycleCheck(nextAction, graph, visited, path); err != nil {
				return err
			}
			path = path[:len(path)-1]
		}
	}

	visited[node] = confirmed

	return nil
}

const (
	notVisited = iota
	visiting
	confirmed
)
