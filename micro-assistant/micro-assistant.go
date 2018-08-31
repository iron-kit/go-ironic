package assistant

import (
	"fmt"
	// "context"
	"github.com/iron-kit/monger"
	"reflect"
)

type Assistant struct {
	Connection monger.Connection
	// Ctx        context.Context
	handers  map[string]Handler
	services map[string]Servicer
}

func (a *Assistant) Handler(name string) Handler {
	return a.handers[name]
}

func (a *Assistant) Service(name string) Servicer {
	return a.services[name]
}

func (a *Assistant) S(v interface{}) Servicer {
	return a.Service(getTypeName(reflect.TypeOf(v)))
}

type AssistantOption func(assistant *Assistant)

func Connection(c monger.Connection) AssistantOption {
	return func(assistant *Assistant) {
		assistant.Connection = c
	}
}

// func Context(ctx context.Context) AssistantOption {
// 	return func(assistant *Assistant) {
// 		assistant.Ctx = ctx
// 	}
// }

func RegisterHandler(hs ...Handler) AssistantOption {
	var handlers map[string]Handler
	return func(assistant *Assistant) {
		if assistant.handers == nil {
			handlers = make(map[string]Handler)
		} else {
			handlers = assistant.handers
		}

		for _, h := range hs {
			h.SetInstance(h)
			// injectToHandler(h, assistant)
			name := getTypeName(reflect.TypeOf(h))
			if _, found := handlers[name]; found {
				panic("Duplicated key '" + name + "'")
			}

			handlers[name] = h
		}

		assistant.handers = handlers
	}
}

func RegisterService(ss ...Servicer) AssistantOption {
	var services map[string]Servicer

	return func(assistant *Assistant) {
		if assistant.services == nil {
			services = make(map[string]Servicer)
		} else {
			services = assistant.services
		}

		for _, s := range ss {
			if assistant.Connection == nil {
				panic("please set connection first")
			}

			// if assistant.Ctx == nil {
			// 	panic("please set context first")
			// }

			s.Init(
				assistant.Connection,
				// assistant.Ctx,
				s,
			)
			// injectToService(s, assistant)
			name := getTypeName(reflect.TypeOf(s))
			if _, found := services[name]; found {
				panic("Duplicated key '" + name + "'")
			}

			services[name] = s
		}

		assistant.services = services
	}
}

func NewAssistant(opts ...AssistantOption) *Assistant {
	assistant := &Assistant{}
	for _, o := range opts {
		o(assistant)
	}
	initHandlers(assistant)
	initServices(assistant)
	return assistant
}

func initHandlers(assistant *Assistant) {
	for _, handler := range assistant.handers {
		injectToHandler(handler, assistant)
	}
}

func initServices(assistant *Assistant) {
	for _, service := range assistant.services {
		injectToService(service, assistant)
	}
}

func injectToHandler(handler Handler, assistant *Assistant) {
	// panic("not impl")
	ht := reflect.TypeOf(handler)
	hv := reflect.ValueOf(handler)
	if ht.Kind() == reflect.Ptr {
		ht = ht.Elem()
	}
	// fmt.Println(ht.Name())
	fieldN := ht.NumField()

	var ihandler Handler
	var servicer Servicer
	handlert := reflect.TypeOf(&ihandler).Elem()
	servicert := reflect.TypeOf(&servicer).Elem()
	for i := 0; i < fieldN; i++ {
		field := ht.Field(i)
		fieldTypeName := getTypeName(field.Type)

		if field.Type.Kind() == reflect.Ptr {
			// fielde := field.Type.Elem()
			fieldv := hv.Elem().Field(i)

			if field.Type.Implements(handlert) {
				fmt.Println("Will inject handler :" + field.Name)
			} else if field.Type.Implements(servicert) {
				fmt.Println("Will inject serivce :" + field.Name)
				// fieldv.Set(reflect.ValueOf(assistant.Service(fieldTypeName)))
				if fieldv.CanSet() {
					fieldv.Set(reflect.ValueOf(assistant.Service(fieldTypeName)))
				}
				// fieldv.Set(reflect.New(field.Type))
			}
		}

	}
}

func injectToService(srv Servicer, assistant *Assistant) {
	srvt := reflect.TypeOf(srv)
	// srvv := reflect.ValueOf(srv)

	if srvt.Kind() == reflect.Ptr {
		srvt = srvt.Elem()
	}

	fieldN := srvt.NumField()

	for i := 0; i < fieldN; i++ {
		field := srvt.Field(i)
		// fieldv := srvv.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			fielde := field.Type.Elem()

			var service Servicer
			servicet := reflect.TypeOf(&service).Elem()
			if fielde.Implements(servicet) {
				// fieldv.Set()
				fmt.Println("Will inject service :" + field.Name)
			}
		}

	}
}
