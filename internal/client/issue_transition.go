package client

import (
	"fmt"
	"strings"
)

// requireTransition enforces the credential's transitions allowlist on a state
// change. An empty allowlist permits everything, and a write that leaves the
// state where it is has no transition to check. Otherwise the pair of current
// state name and selected state name must be listed, or the write is refused
// before any mutation is sent. The selected state was resolved from the team's
// cached workflow states a moment earlier, so its name is read from that cache.
func (guard *guardedClient) requireTransition(issue IssueSummary, stateID string) error {
	allowlist := guard.target.Expected.Transitions
	if len(allowlist) == 0 || issue.StateID == stateID {
		return nil
	}
	to, ok := guard.cachedStateName(issue.TeamID, stateID)
	if !ok {
		return fmt.Errorf(
			"%w: state %s is not a cached workflow state of team %s", ErrWriteInvalid, stateID, issue.Team,
		)
	}
	if transitionAllowed(allowlist, issue.State, to) {
		return nil
	}

	return fmt.Errorf(
		"%w: %q -> %q is not in the [transitions] allowlist for %q",
		ErrTransitionDenied, issue.State, to, issue.State,
	)
}

func (guard *guardedClient) cachedStateName(teamID string, stateID string) (string, bool) {
	guard.stateIDs.mu.Lock()
	defer guard.stateIDs.mu.Unlock()
	for _, state := range guard.stateIDs.lists[teamID] {
		if state.ID == stateID {
			return state.Name, true
		}
	}

	return "", false
}

func transitionAllowed(allowlist map[string][]string, from string, to string) bool {
	for current, next := range allowlist {
		if !strings.EqualFold(current, from) {
			continue
		}
		for _, candidate := range next {
			if strings.EqualFold(candidate, to) {
				return true
			}
		}
	}

	return false
}
