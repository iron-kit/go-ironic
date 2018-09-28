package ironic

import (
	"context"
	"errors"
	"github.com/go-log/log"
	"github.com/iron-kit/monger"
	"github.com/micro/go-micro/server"
	"time"
)

type errorManagerKey struct{}

func ErrorWrapper(name string) server.HandlerWrapper {
	errManager := NewErrorManager(name)
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, errorManagerKey{}, errManager)

			return fn(ctx, req, rsp)
		}
	}
}

func ErrorManagerFromContext(ctx context.Context) *ErrorManager {
	em, ok := ctx.Value(errorManagerKey{}).(*ErrorManager)

	if !ok {
		return NewErrorManager("")
	}

	return em
}

func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		log.Logf("[%v] server request: %s", time.Now(), req.Method())
		return fn(ctx, req, rsp)
	}
}

type MongerWrapperCallback func(monger.Connection) error
type mongerKey struct{}

func MongerWrapper(fn MongerWrapperCallback, opts ...monger.ConfigOption) server.HandlerWrapper {
	connection, err := monger.Connect(opts...)
	if err != nil {
		panic(err.Error())
	}
	if err := fn(connection); err != nil {
		panic(err.Error())
	}
	// connection.BatchRegister()
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			ctx = context.WithValue(ctx, mongerKey{}, connection)
			return fn(ctx, req, rsp)
		}
	}
}

func MongerConnectionFromContext(ctx context.Context) (monger.Connection, error) {
	conn, ok := ctx.Value(mongerKey{}).(monger.Connection)

	if !ok {
		return nil, errors.New("Not found connection from this context")
	}

	return conn, nil
}
