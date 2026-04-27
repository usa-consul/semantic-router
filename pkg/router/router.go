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
		route.Threshold = 0.8 // sensible default
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.routes = append(r.routes, route)
	r.logger.Info("route registered", zap.String("name", route.Name), zap.Int("utterances", len(route.Utterances)))
	return nil
}

// Match finds the best matching route for the given query.
// Returns nil if no route meets its similarity threshold.
func (r *Router) Match(ctx context.Context, query string) (*Route, float64, error) {
	if query == "" {
		return nil, 0, errors.New("router: query must not be empty")
	}
	queryVec, err := r.encoder.Encode(ctx, query)
	if err != nil {
		return nil, 0, fmt.Errorf("router: failed to encode query: %w", err)
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	var (
		bestRoute *Route
		bestScore float64
	)
	for _, route := range r.routes {
		for _, utterance := range route.Utterances {
			utteranceVec, err := r.encoder.Encode(ctx, utterance)
			if err != nil {
				r.logger.Warn("failed to encode utterance", zap.String("route", route.Name), zap.Error(err))
				continue
			}
			score := cosineSimilarity(queryVec, utteranceVec)
			if score > bestScore {
				bestScore = score
				bestRoute = route
			}
		}
	}
	if bestRoute == nil || bestScore < bestRoute.Threshold {
		return nil, bestScore, nil
	}
	return bestRoute, bestScore, nil
}

// cosineSimilarity computes the cosine similarity between two vectors.
func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (sqrt(normA) * sqrt(normB))
}

// sqrt is a simple Newton-Raphson square root to avoid importing math.
func sqrt(x float64) float64 {
	if x <= 0 {
		return 0
	}
	z := x
	for i := 0; i < 20; i++ {
		z -= (z*z - x) / (2 * z)
	}
	return z
}
