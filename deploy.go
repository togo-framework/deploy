// Package deploy is the togo deployment subsystem — a provider-agnostic Deployer
// contract. Real targets (terraform, docker, kubernetes, the clouds, and VPS
// distros) ship as driver plugins that call deploy.RegisterDriver in their init();
// pick one with DEPLOY_PROVIDER (or `deploy.provider` in togo.yaml). The CLI's
// `togo deploy` resolves the provider and calls Build/Deploy without booting the app.
package deploy

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/togo-framework/togo"
)

// Spec describes what to deploy. Drivers use the fields they need.
type Spec struct {
	App      string            // app name
	Dir      string            // project directory
	BuildCmd string            // local build command (optional)
	Binary   string            // built artifact path (for VPS drivers)
	Image    string            // container image ref (for docker/k8s)
	Env      map[string]string // runtime env
	Host     string            // target host / cluster endpoint
	User     string            // ssh user (VPS drivers)
	Domain   string            // public domain
	Region   string            // cloud region
	Options  map[string]any    // provider-specific knobs
}

// Result of a deploy/provision.
type Result struct {
	URL     string         // public URL if known
	Message string         // human summary
	Raw     map[string]any // provider raw response
}

// Status of the current deployment.
type Status struct {
	Healthy bool
	Detail  string
	Raw     map[string]any
}

// Deployer is implemented by driver plugins. Not every target supports every
// operation — return a clear error for the unsupported ones.
type Deployer interface {
	Provision(ctx context.Context, spec Spec) (*Result, error) // create/ensure infrastructure
	Deploy(ctx context.Context, spec Spec) (*Result, error)    // ship the app
	Destroy(ctx context.Context, spec Spec) error              // tear it down
	Status(ctx context.Context, spec Spec) (*Status, error)
}

// DriverFactory builds a Deployer from the kernel (env-configured).
type DriverFactory func(k *togo.Kernel) (Deployer, error)

var (
	regMu   sync.RWMutex
	drivers = map[string]DriverFactory{}
)

// RegisterDriver registers a deploy driver by name (call from a plugin's init()).
func RegisterDriver(name string, f DriverFactory) {
	regMu.Lock()
	drivers[name] = f
	regMu.Unlock()
}

// Drivers lists the registered driver names (sorted) — handy for `togo deploy --list`.
func Drivers() []string {
	regMu.RLock()
	defer regMu.RUnlock()
	out := make([]string, 0, len(drivers))
	for n := range drivers {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}

// Build constructs the Deployer for name. The CLI uses this to deploy without
// booting the full kernel (pass a minimal/nil kernel where the driver allows).
func Build(name string, k *togo.Kernel) (Deployer, error) {
	regMu.RLock()
	f, ok := drivers[name]
	regMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("deploy: unknown provider %q (install its plugin, e.g. togo install togo-framework/deploy-%s)", name, name)
	}
	return f(k)
}

func init() {
	// Safe default: log instead of touching infrastructure.
	RegisterDriver("log", func(k *togo.Kernel) (Deployer, error) { return &logDeployer{k: k}, nil })

	togo.RegisterProviderFunc("deploy", togo.PriorityLate, func(k *togo.Kernel) error {
		name := os.Getenv("DEPLOY_PROVIDER")
		if name == "" {
			name = "log"
		}
		d, err := Build(name, k)
		if err != nil {
			return err
		}
		k.Set("deploy", &Service{deployer: d, driver: name})
		return nil
	})
}

// Service is the deploy runtime stored on the kernel (k.Get("deploy")).
type Service struct {
	deployer Deployer
	driver   string
}

func (s *Service) Deployer() Deployer { return s.deployer }
func (s *Service) Driver() string     { return s.driver }
func (s *Service) Provision(ctx context.Context, spec Spec) (*Result, error) {
	return s.deployer.Provision(ctx, spec)
}
func (s *Service) Deploy(ctx context.Context, spec Spec) (*Result, error) {
	return s.deployer.Deploy(ctx, spec)
}
func (s *Service) Destroy(ctx context.Context, spec Spec) error { return s.deployer.Destroy(ctx, spec) }
func (s *Service) Status(ctx context.Context, spec Spec) (*Status, error) {
	return s.deployer.Status(ctx, spec)
}

// FromKernel fetches the deploy service from the kernel container.
func FromKernel(k *togo.Kernel) (*Service, bool) {
	v, ok := k.Get("deploy")
	if !ok {
		return nil, false
	}
	s, ok := v.(*Service)
	return s, ok
}

type logDeployer struct{ k *togo.Kernel }

func (l *logDeployer) log(msg string, spec Spec) *Result {
	if l.k != nil && l.k.Log != nil {
		l.k.Log.Info("deploy (log driver) "+msg, "app", spec.App, "host", spec.Host)
	}
	return &Result{Message: "log: would " + msg + " " + spec.App}
}
func (l *logDeployer) Provision(_ context.Context, s Spec) (*Result, error) {
	return l.log("provision", s), nil
}
func (l *logDeployer) Deploy(_ context.Context, s Spec) (*Result, error) { return l.log("deploy", s), nil }
func (l *logDeployer) Destroy(_ context.Context, _ Spec) error           { return nil }
func (l *logDeployer) Status(_ context.Context, _ Spec) (*Status, error) {
	return &Status{Healthy: true, Detail: "log driver — no real target"}, nil
}
