package tasks

import "context"

type Service struct{ repo *Repo }

func NewService(r *Repo) *Service { return &Service{repo: r} }

func (s *Service) List(ctx context.Context) ([]Task, error) {
	return s.repo.List(ctx)
}
