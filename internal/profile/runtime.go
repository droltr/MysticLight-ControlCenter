package profile

import (
	"fmt"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
)

type Runtime struct {
	providers map[string]struct{}
	active    domain.Profile
	fallback  domain.Profile
}

func NewRuntime(providers map[string]struct{}, fallback domain.Profile) (*Runtime, error) {
	if err := fallback.Validate(providers); err != nil {
		return nil, fmt.Errorf("invalid fallback profile: %w", err)
	}
	return &Runtime{
		providers: providers,
		fallback:  fallback,
	}, nil
}

func (r *Runtime) Activate(profile domain.Profile) error {
	if err := profile.Validate(r.providers); err != nil {
		return fmt.Errorf("profile activation rejected: %w", err)
	}
	r.active = profile
	return nil
}

func (r *Runtime) Active() domain.Profile {
	return r.active
}

func (r *Runtime) Fallback() domain.Profile {
	return r.fallback
}

func (r *Runtime) Recover() domain.Profile {
	r.active = r.fallback
	return r.active
}
