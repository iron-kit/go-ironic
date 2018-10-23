package ironic

import (
	"github.com/micro/go-micro"
)

type ironicService struct {
	service micro.Service
}

// NewService
func NewService(opts ...micro.Option) micro.Service {
	// opts = append(opts, micro.Context())
	opts = append(opts, micro.WrapHandler(
		LogWrapper,
		// ErrorWrapper("DEFAULT"),
	))

	srv := micro.NewService(opts...)
	name := srv.Server().Options().Name

	srv.Init(
		micro.WrapHandler(ErrorWrapper(name)),
	)
	return srv
}
