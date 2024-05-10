package scheduler

// Scheduler determines a set of candidate workers on which task it could run,
// score the candidate workers from best to worst, and pick the worker with best
// score
type Scheduler interface {
	SelectCandidateNodes()
	Score()
	Pick()
}
