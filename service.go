package ironic

import (
	"context"
	"errors"
)

type (
	Servicer interface {
		init(ctx context.Context)
		Ctx() context.Context
		Error() *ErrorManager
	}

	Service struct {
		ctx context.Context
		err *ErrorManager
	}
)

func InitServiceFunc(s Servicer, ctx context.Context) error {
	if ctx == nil {
		return errors.New("Context can't be null")
	}

	if s.Ctx() == nil {
		s.init(ctx)
	}

	return nil
}

func (s *Service) init(ctx context.Context) {
	s.ctx = ctx
}

func (s *Service) Ctx() context.Context {
	return s.ctx
}

func (s *Service) Error() *ErrorManager {
	if s.err == nil {
		errm := ErrorManagerFromContext(s.ctx)

		s.err = errm
	}

	return s.err
}
