package routing

import (
	"context"
	"sync"

	"github.com/dogmatiq/dogma"
)

// app routes [triggerCommand] through [triggerIntegration] to produce
// [observed], which [projection] handles.
type app struct {
	Projection projection
}

func (a *app) Configure(c dogma.ApplicationConfigurer) {
	c.Identity("app", "c1a00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.ViaIntegration(&triggerIntegration{}),
		dogma.ViaProjection(&a.Projection),
	)
}

// triggerIntegration handles [triggerCommand] and records [observed].
type triggerIntegration struct{}

func (*triggerIntegration) Configure(c dogma.IntegrationConfigurer) {
	c.Identity("trigger", "c1b00000-0000-0000-0000-000000000000")
	c.Routes(
		dogma.HandlesCommand[*triggerCommand](),
		dogma.RecordsEvent[*observed](),
	)
}

func (*triggerIntegration) HandleCommand(_ context.Context, s dogma.IntegrationCommandScope, _ dogma.Command) error {
	s.RecordEvent(&observed{})
	return nil
}

// projection handles [observed] and signals [Received] on each invocation. It
// maintains a checkpoint map to satisfy the OCC protocol.
type projection struct {
	dogma.NoCompactBehavior
	dogma.NoResetBehavior

	mu          sync.Mutex
	checkpoints map[string]uint64
	Received    chan struct{}
}

func (p *projection) Configure(c dogma.ProjectionConfigurer) {
	c.Identity("projection", "c1c00000-0000-0000-0000-000000000000")
	c.Routes(dogma.HandlesEvent[*observed]())
}

func (p *projection) HandleEvent(
	_ context.Context,
	s dogma.ProjectionEventScope,
	_ dogma.Event,
) (uint64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	cp := s.Offset() + 1
	if p.checkpoints == nil {
		p.checkpoints = map[string]uint64{}
	}
	p.checkpoints[s.StreamID()] = cp
	p.Received <- struct{}{}
	return cp, nil
}

func (p *projection) CheckpointOffset(_ context.Context, id string) (uint64, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.checkpoints[id], nil
}
