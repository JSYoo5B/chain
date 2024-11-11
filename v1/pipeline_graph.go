package v1

import "fmt"

// ValidateGraph ensures the pipeline's graph is connected and acyclic.
// It checks for cycles first, then verifies that all nodes are reachable starting from the initAction.
// If any node is unreachable, it starts a new traversal from the disconnected node and checks the union of visited nodes.
// If the union of all visited nodes equals the total number of nodes, the graph is connected and the check finishes.
func (p *Pipeline[T]) ValidateGraph() error {
	// Step 1: Perform DFS from initAction to check for cycles and track visited nodes
	visited := make(map[Action[T]]bool)
	if err := p.dfsWithCycleCheck(p.initAction, visited); err != nil {
		return err
	}

	// Step 2: After DFS, check if all actions have been visited
	unvisited := make([]Action[T], 0, len(p.runPlans))
	for action := range p.runPlans {
		if !visited[action] {
			unvisited = append(unvisited, action)
		}
	}

	// Step 3: If there are unvisited nodes, start new DFS from them
	for len(unvisited) > 0 {
		newStart := unvisited[0] // Pick any unvisited node
		visitedFromNewStart := make(map[Action[T]]bool)
		if err := p.dfsWithCycleCheck(newStart, visitedFromNewStart); err != nil {
			return err
		}

		// Step 4: Merge the visited nodes from the current traversal into the overall visited set
		// If the current traversal's visited nodes intersect with the previously visited ones, they are connected
		// If there is no intersection, it's a disconnected graph.
		intersectionFound := false
		for action := range visitedFromNewStart {
			if visited[action] {
				intersectionFound = true
			}
			visited[action] = true
		}

		// Step 5: If no intersection found, it means the graph is disconnected
		if !intersectionFound {
			return fmt.Errorf("disconnect detected: action `%s` cannot reach the graph started from initAction `%s`", newStart.Name(), p.initAction.Name())
		}

		// Check if the union of visited nodes equals the total number of actions
		if len(visited) == len(p.runPlans) {
			return nil // All nodes have been visited, no need for further checks
		}

		// If there are still unvisited nodes, continue to the next round of DFS
		stillUnvisited := make([]Action[T], 0, len(unvisited))
		for _, action := range unvisited {
			if !visited[action] {
				stillUnvisited = append(stillUnvisited, action)
			}
		}
		unvisited = stillUnvisited
	}

	return nil
}

func (p *Pipeline[T]) dfsWithCycleCheck(current Action[T], visited map[Action[T]]bool) error {
	if visited[current] {
		return fmt.Errorf("cycle detected: action `%s` has been visited twice", current.Name())
	}

	visited[current] = true

	// Get the next actions based on current action's direction
	plan := p.runPlans[current]
	for _, nextAction := range plan {
		if nextAction != nil {
			if err := p.dfsWithCycleCheck(nextAction, visited); err != nil {
				return err
			}
		}
	}

	return nil
}
