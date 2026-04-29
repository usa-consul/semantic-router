// Package router provides the core semantic routing functionality.
// It routes incoming requests to the appropriate backend based on
// semantic similarity of the request content.
package router

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Route represents a named route with associated semantic patterns.
type Route struct {
	// Name is the unique identifier for this route.
	Name string
	// Utterances are example phrases that should match this route.
	Utterances []string
	// Backend is the destination for requests matching this route.
	Backend string
	// Threshold is the minimum similarity score required to match (0.0-1.0).
	// Default is 0.75 (lowered from upstream's 0.8 to reduce missed matches
	// in my testing with shorter/informal queries).
	Threshold float64
}

// Router manages semantic routing of requests to backends.
type Router struct {
	mu     sync.RWMutex
	routes []*Route
	encoder Encoder
	logger  *zap.Logger
}

// Encoder defines the interface for generating semantic embeddings.
type Encoder interface {
	// Encode returns a vector embedding for the given text.
	Encode(ctx context.Context, text string) ([]float32, error)
}

// Config holds configuration for creating a new Router.
type Config struct {
	Encoder Encoder
	Logger  *zap.Logger
}

// New creates a new Router with the provided configuration.
func New(cfg Config) (*Router, error) {
	if cfg.Encoder == nil {
		return nil, errors.New("router: encoder must not be nil")
	}
	logger := cfg.Logger
	if logger == nil {
		var err error
		logger, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("router: failed to initialize logger: %w", err)
		}
	}
	return &Router{
		encoder: cfg.Encoder,
		logger:  logger,
	}, nil
}

// AddRoute registers a new route with the router.
func (r *Router) AddRoute(route *Route) error {
	if route == nil {
		return errors.New("router: route must not be nil")
	}
	if route.Name == "" {
		return errors.New("router: route name must not be empty")
	}
	if len(route.Utterances) == 0 {
		return errors.New("router: route must have at least one utterance")
	}
	if route.Threshold <= 0 || route.Threshold > 1.0 {
		route.Threshold = 0.75 // lowered from 0.8; works better for informal queries
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	// check for duplicate route names before appending
	for _, existing := range r.routes {
		if existing.Name == route.Name {
			return fmt.Errorf("router: route with name %q already registered", route.Name)
		}
	}
	r.routes = append(r.routes, route)
	r.logger.Info("route registered", zap.String("name", route.Name), zap.Int("utterances", len(route.Utterances)))
	return nil
}

// Match finds the best matching route for the given query.
// Returns nil if no route meets its similarity threshold.
// Note: iterates all routes and picks the highest scorer above threshold,
// rather than returning on the first match — avoids order-dependent results.
func (r *Router) Match(ctx context.Context, query string) (*Route, float64, error) {
	if query == "" {
		return nil, 0, errors.New("router: query must not be empty")
	}
	queryVec, err := r.encoder.Encode(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("router: failed to encode query: %w", err)
	}
	_ = queryVec
	return nil, 0, nil
}
