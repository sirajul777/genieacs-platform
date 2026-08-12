package job

import "testing"

func TestTerminalStatesCannotTransition(t *testing.T) {
	terminal := []Status{StatusCompleted, StatusFailed, StatusCanceled}
	active := []Status{StatusPending, StatusInProgress}

	for _, from := range terminal {
		for _, to := range active {
			if from.CanTransitionTo(to) {
				t.Fatalf("terminal state %q can transition to %q", from, to)
			}
		}
		if !from.IsTerminal() {
			t.Fatalf("state %q should be terminal", from)
		}
	}
}

func TestJobLifecycleTransitions(t *testing.T) {
	valid := map[Status][]Status{
		StatusPending:    {StatusInProgress, StatusCanceled},
		StatusInProgress: {StatusCompleted, StatusFailed, StatusCanceled},
	}

	for from, targets := range valid {
		for _, to := range targets {
			if !from.CanTransitionTo(to) {
				t.Fatalf("expected transition %q -> %q", from, to)
			}
		}
	}
}
