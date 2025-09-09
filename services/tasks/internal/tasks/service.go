package tasks

import "context"

// Service is the domain layer for "tasks" use-cases.
// It depends on a Repo to persist/fetch data, but exposes
// business-oriented methods to the HTTP/GraphQL layers.
type Service struct{
	repo *Repo
}

// NewService wires a Service to a concrete Repo.
func NewService(r *Repo) *Service {
	return &Service{repo: r}
}

// List returns all tasks for the current request.
// It forwards the request Context so cancellations/timeouts propagate
// down to the DB queries via the repo.
func (s *Service) List(ctx context.Context) ([]Task, error) {
	return s.repo.List(ctx)
}

// Create adds a new task with the given title.
func (s *Service) Create(ctx context.Context, title string) (Task, error) {
	return s.repo.Create(ctx, title)
}

// A note about the func (receiver) syntax.
//
// This would be the same as the following in Python:
//
// class Service:
//     def __init__(self, repo):
//         self.repo = repo
//
//     def list(self, ctx):
//         """
//         Fetch a list of tasks using the repository.
//         ctx: like Go's context.Context, could carry deadlines or cancellation signals.
//         Returns: (tasks, error) where error is None if all went well.
//         """
//         return self.repo.list(ctx)
//
// While Go doesn’t have classes, it has types (struct, interface, etc.) and you attach
// methods to them using the func (receiver) syntax.
