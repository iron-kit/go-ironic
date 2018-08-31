package assistant

import (
	// "context"
	"github.com/iron-kit/monger"
)

type (
	Servicer interface {
		Init(args ...interface{})
		Model(name string) monger.Model
	}
)

type (
	Service struct {
		Connection monger.Connection
		// Ctx        context.Context
		instance Servicer
	}
)

/*
Init is a inital function before you use this service
*/
func (s *Service) Init(args ...interface{}) {
	for _, v := range args {
		if conn, ok := v.(monger.Connection); ok {
			s.Connection = conn
		}

		// if ctx, ok := v.(context.Context); ok {
		// 	s.Ctx = ctx
		// }

		if srv, ok := v.(Servicer); ok {
			s.instance = srv
		}

	}
}

// Model is a function to get monger model by model's name
func (s *Service) Model(name string) monger.Model {
	return s.Connection.M(name)
}

// func (sm *serviceManager) Srv(name string) Servicer {
// 	return sm.services[name]
// }

// func (sm *serviceManager) Register(servicers ...Servicer) {
// 	for _, v := range servicers {
// 		stype := reflect.TypeOf(v)

// 		if stype.Kind() == reflect.Ptr {
// 			stype = stype.Elem()
// 		}
// 		name := stype.Name()

// 		if _, ok := sm.services[name]; ok {
// 			// TODO right panic
// 			panic(fmt.Sprintf("重复注册了Service %T", v))
// 		}

// 		sm.services[stype.Name()] = v
// 	}
// }
