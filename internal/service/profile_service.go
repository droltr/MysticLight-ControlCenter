package service

import (
	"context"
	"fmt"

	"github.com/droltr/MysticLight-ControlCenter/internal/domain"
	"github.com/droltr/MysticLight-ControlCenter/internal/profile"
)

type ProfileService struct {
	runtime   *profile.Runtime
	path      string
	providers map[string]struct{}
}

func NewProfileService(providers map[string]struct{}, fallback domain.Profile, path string) (*ProfileService, error) {
	runtime, err := profile.NewRuntime(providers, fallback)
	if err != nil {
		return nil, err
	}
	return &ProfileService{runtime: runtime, path: path, providers: providers}, nil
}

func (s *ProfileService) Load(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	loaded, err := profile.LoadFile(s.path, s.providers)
	if err != nil {
		return fmt.Errorf("load active profile: %w", err)
	}
	return s.runtime.Activate(loaded)
}

func (s *ProfileService) Activate(ctx context.Context, candidate domain.Profile) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if err := s.runtime.Activate(candidate); err != nil {
		return err
	}
	return profile.SaveFile(s.path, candidate)
}

func (s *ProfileService) Recover() domain.Profile {
	return s.runtime.Recover()
}

func (s *ProfileService) Active() domain.Profile {
	return s.runtime.Active()
}
