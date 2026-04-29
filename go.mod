module github.com/vllm-project/semantic-router

go 1.22

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/cors v1.2.1
	github.com/prometheus/client_golang v1.19.1
	go.uber.org/zap v1.27.0
	k8s.io/apimachinery v0.30.2
	sigs.k8s.io/controller-runtime v0.18.4
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.54.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	k8s.io/utils v0.0.0-20240502163921-fe8a2dddb1d0 // indirect
)

// Personal fork for learning purposes - experimenting with routing logic
// TODO: look into replacing go-chi/cors with a more configurable middleware
// to better understand how CORS preflight requests interact with the router
//
// NOTE: Tried upgrading golang.org/x/sys to v0.22.0 but it caused test flakiness
// in the controller-runtime integration tests - staying on v0.21.0 for now
//
// TODO: investigate upgrading sigs.k8s.io/controller-runtime to v0.19.x once
// the k8s.io/apimachinery dependency catches up - the reconciler loop changes
// in v0.19 look interesting for understanding how watches are batched
//
// NOTE: module path intentionally kept as vllm-project to make it easier to
// track upstream changes and eventually submit patches back upstream
//
// NOTE: prometheus/client_golang v1.20.x dropped support for the legacy
// push gateway client - worth checking before upgrading if any metrics
// collection relies on that path
//
// NOTE: go-chi/chi v5.2.x added native support for route groups with middleware
// inheritance - worth upgrading from v5.1.0 to experiment with cleaner
// middleware chaining in the router handler setup
