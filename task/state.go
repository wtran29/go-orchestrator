package task

type State int

const (
	Pending   State = iota // initial state, starting point for every task
	Scheduled              // a task moves to this state once manager has scheduled to worker
	Running                // when a worker successfully starts the task (i.e. starts container)
	Completed              // when it completes its work in a normal way (i.e. does not fail)
	Failed                 // when a task fails
)

var stateTransitionMap = map[State][]State{
	Pending:   {Scheduled},
	Scheduled: {Scheduled, Running, Failed},
	Running:   {Running, Completed, Failed},
	Completed: {},
	Failed:    {},
}

func Contains(states []State, state State) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}

func ValidStateTransition(src State, dst State) bool {
	return Contains(stateTransitionMap[src], dst)
}
