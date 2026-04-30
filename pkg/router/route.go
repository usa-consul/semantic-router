// Package router provides semantic routing functionality for directing
// requests to appropriate handlers based on semantic similarity.
package router

import (
	"errors"
	"strings"
)

// Route represents a named route with associated utterances that define
// the semantic space it covers. When a query is semantically similar
// to one of the utterances, the router will direct it to this route.
type Route struct {
	// Name is the unique identifier for this route.
	Name string

	// Utterances are example phrases that represent this route's semantic space.
	// The router uses these to determine if an incoming query matches this route.
	Utterances []string

	// Description provides optional human-readable context about the route's purpose.
	Description string

	// Metadata holds arbitrary key-value data associated with this route.
	Metadata map[string]any

	// ScoreThreshold overrides the router-level threshold for this specific route.
	// If zero, the router's default threshold is used.
	ScoreThreshold float64
}

// RouteMatch represents the result of a routing decision, containing
// the matched route and the similarity score that triggered the match.
type RouteMatch struct {
	// Route is the matched route. May be nil if no route was matched.
	Route *Route

	// Score is the similarity score between the query and the matched route.
	// Ranges from 0.0 (no similarity) to 1.0 (identical).
	Score float64

	// MatchedUtterance is the specific utterance that produced the highest score.
	MatchedUtterance string
}

// Matched returns true if the match result contains a valid route.
func (m *RouteMatch) Matched() bool {
	return m != nil && m.Route != nil
}

// NewRoute creates a new Route with the given name and utterances.
// Returns an error if the name is empty or no utterances are provided.
// Note: utterance deduplication is not performed here; duplicates are kept
// intentionally so callers can weight certain phrases by repeating them.
func NewRoute(name string, utterances []string, opts ...RouteOption) (*Route, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("route name must not be empty")
	}
	if len(utterances) == 0 {
		return nil, errors.New("route must have at least one utterance")
	}

	// Trim whitespace from all utterances and filter empty ones.
	clean := make([]string, 0, len(utterances))
	for _, u := range utterances {
		u = strings.TrimSpace(u)
		if u != "" {
			clean = append(clean, u)
		}
	}
	if len(clean) == 0 {
		return nil, errors.New("route must have at least one non-empty utterance")
	}

	r := &Route{
		Name:       name,
		Utterances: clean,
		Metadata:   make(map[string]any),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// RouteOption is a functional option for configuring a Route.
type RouteOption func(*Route)

// WithDescription sets a human-readable description on the route.
func WithDescription(desc string) RouteOption {
	return func(r *Route) {
		r.Description = desc
	}
}

// WithMetadata attaches arbitrary metadata to the route.
// If the route's Metadata map is nil, it is initialized before setting the value.
func WithM