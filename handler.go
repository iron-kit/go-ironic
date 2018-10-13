package ironic

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"mime/multipart"
	"strings"

	restful "github.com/emicklei/go-restful"

	"github.com/iron-kit/go-ironic/validator"
	go_api "github.com/micro/go-api/proto"
)

// HTTP methods
const (
	CONNECT  = "CONNECT"
	DELETE   = "DELETE"
	GET      = "GET"
	HEAD     = "HEAD"
	OPTIONS  = "OPTIONS"
	PATCH    = "PATCH"
	POST     = "POST"
	PROPFIND = "PROPFIND"
	PUT      = "PUT"
	TRACE    = "TRACE"
)

const (
	charsetUTF8 = "charset=UTF-8"
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

// Headers
const (
	HeaderAccept              = "Accept"
	HeaderAcceptEncoding      = "Accept-Encoding"
	HeaderAllow               = "Allow"
	HeaderAuthorization       = "Authorization"
	HeaderContentDisposition  = "Content-Disposition"
	HeaderContentEncoding     = "Content-Encoding"
	HeaderContentLength       = "Content-Length"
	HeaderContentType         = "Content-Type"
	HeaderCookie              = "Cookie"
	HeaderSetCookie           = "Set-Cookie"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderLastModified        = "Last-Modified"
	HeaderLocation            = "Location"
	HeaderUpgrade             = "Upgrade"
	HeaderVary                = "Vary"
	HeaderWWWAuthenticate     = "WWW-Authenticate"
	HeaderXForwardedFor       = "X-Forwarded-For"
	HeaderXForwardedProto     = "X-Forwarded-Proto"
	HeaderXForwardedProtocol  = "X-Forwarded-Protocol"
	HeaderXForwardedSsl       = "X-Forwarded-Ssl"
	HeaderXUrlScheme          = "X-Url-Scheme"
	HeaderXHTTPMethodOverride = "X-HTTP-Method-Override"
	HeaderXRealIP             = "X-Real-IP"
	HeaderXRequestID          = "X-Request-ID"
	HeaderXRequestedWith      = "X-Requested-With"
	HeaderServer              = "Server"
	HeaderOrigin              = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity = "Strict-Transport-Security"
	HeaderXContentTypeOptions     = "X-Content-Type-Options"
	HeaderXXSSProtection          = "X-XSS-Protection"
	HeaderXFrameOptions           = "X-Frame-Options"
	HeaderContentSecurityPolicy   = "Content-Security-Policy"
	HeaderXCSRFToken              = "X-CSRF-Token"
)

type Handler interface {
	Error(ctx context.Context) *ErrorManager
	Validate(v interface{}) error
	Bind(req *go_api.Request, params interface{}) error
}

type BaseHandler struct {
	binder GoAPIBinder
}

type BaseWebHandler struct {
	binder Binder
}

func (h *BaseWebHandler) Error() *ErrorManager {
	return NewErrorManager("iunite.club.navo")
}

func (h *BaseWebHandler) Validate(v interface{}) error {
	return validator.Validate(v)
}

func (h *BaseWebHandler) Bind(req *restful.Request, params interface{}) error {
	// var err error
	// req.Get
	// req.Get[]

	if h.binder == nil {
		h.binder = &DefaultBinderWithRestful{}
	}

	return h.binder.Bind(req, params)
	// byteBody := []byte(req.Body)
	// return json.Unmarshal(byteBody, params)
}

func (h *BaseHandler) Error(ctx context.Context) *ErrorManager {
	return ErrorManagerFromContext(ctx)
	// return NewErrorManager("iunite.club.navo")
}

func (h *BaseHandler) Validate(v interface{}) error {
	return validator.Validate(v)
}

func (h *BaseHandler) Bind(req *go_api.Request, params interface{}) error {
	// var err error
	// req.Get
	// req.Get[]

	if h.binder == nil {
		h.binder = &DefaultBinder{}
	}

	return h.binder.Bind(req, params)
	// byteBody := []byte(req.Body)
	// return json.Unmarshal(byteBody, params)
}

func (h *BaseHandler) MultipartForm(req *go_api.Request) (*multipart.Form, error) {
	ct := strings.Join(req.Header["Content-Type"].Values, ",")
	mt, p, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}
	if !strings.HasPrefix(mt, MIMEMultipartForm) {
		return nil, fmt.Errorf("%v does not contain multipart", mt)
	}
	r := multipart.NewReader(bytes.NewReader([]byte(req.Body)), p["boundary"])
	return r.ReadForm(32 << 20)
}

func (h *BaseHandler) FormFile(req *go_api.Request, field string) (*multipart.FileHeader, error) {
	form, err := h.MultipartForm(req)
	if err != nil {
		return nil, err
	}

	if f, found := form.File[field]; found {
		return f[0], nil
	}

	return nil, fmt.Errorf("field %v not found", field)
}

func (h *BaseHandler) FormFiles(req *go_api.Request, field string) ([]*multipart.FileHeader, error) {
	form, err := h.MultipartForm(req)
	if err != nil {
		return nil, err
	}

	if f, found := form.File[field]; found {
		return f, nil
	}

	return nil, fmt.Errorf("field %v not found", field)
}

// func (h *BaseHandler) bindData(ptr interface{}, data map[string][]string, tag string) error {
// 	typ := reflect.TypeOf(ptr).Elem()
// 	val := reflect.ValueOf(ptr).Elem()

// 	if typ.Kind() != reflect.Struct {
// 		return errors.New("binding element must be a struct")
// 	}

// 	for i := 0; i < typ.NumField(); i ++ {
// 		typeField := typ.Field(i)
// 		structField := val.Field(i)

// 		if !structField.CanSet() {
// 			continue
// 		}

// 		sfKind := structField.Kind()
// 		inputFieldName := typeField.Tag.Get(tag)

// 		if inputFieldName == "" {
// 			inputFieldName = typeField.Name
// 			// If tag is nil, we inspect if the field is a struct.
// 			if _, ok := bindUnmarshaler(structField); !ok && structFieldKind == reflect.Struct {
// 				if err := b.bindData(structField.Addr().Interface(), data, tag); err != nil {
// 					return err
// 				}
// 				continue
// 			}
// 		}
// 	}
// }

// func (h *BaseHandler) bindData(ptr interface{}, data map[string][]string, tag string) error {
// 	typ := reflect.TypeOf(ptr).Elem()
// 	val := reflect.ValueOf(ptr).Elem()

// 	if typ.Kind() != reflect.Struct {
// 		return errors.New("binding element must be a struct")
// 	}

// 	for i := 0; i < typ.NumField(); i++ {
// 		typeField := typ.Field(i)
// 		structField := val.Field(i)
// 		if !structField.CanSet() {
// 			continue
// 		}
// 		structFieldKind := structField.Kind()
// 		inputFieldName := typeField.Tag.Get(tag)

// 		if inputFieldName == "" {
// 			inputFieldName = typeField.Name
// 			// If tag is nil, we inspect if the field is a struct.
// 			if _, ok := bindUnmarshaler(structField); !ok && structFieldKind == reflect.Struct {
// 				if err := h.bindData(structField.Addr().Interface(), data, tag); err != nil {
// 					return err
// 				}
// 				continue
// 			}
// 		}

// 		inputValue, exists := data[inputFieldName]
// 		if !exists {
// 			// Go json.Unmarshal supports case insensitive binding.  However the
// 			// url params are bound case sensitive which is inconsistent.  To
// 			// fix this we must check all of the map values in a
// 			// case-insensitive search.
// 			inputFieldName = strings.ToLower(inputFieldName)
// 			for k, v := range data {
// 				if strings.ToLower(k) == inputFieldName {
// 					inputValue = v
// 					exists = true
// 					break
// 				}
// 			}
// 		}

// 		if !exists {
// 			continue
// 		}

// 		// Call this first, in case we're dealing with an alias to an array type
// 		if ok, err := unmarshalField(typeField.Type.Kind(), inputValue[0], structField); ok {
// 			if err != nil {
// 				return err
// 			}
// 			continue
// 		}

// 		numElems := len(inputValue)
// 		if structFieldKind == reflect.Slice && numElems > 0 {
// 			sliceOf := structField.Type().Elem().Kind()
// 			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
// 			for j := 0; j < numElems; j++ {
// 				if err := setWithProperType(sliceOf, inputValue[j], slice.Index(j)); err != nil {
// 					return err
// 				}
// 			}
// 			val.Field(i).Set(slice)
// 		} else if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
// 			return err

// 		}
// 	}
// 	return nil
// }

// // bindUnmarshaler attempts to unmarshal a reflect.Value into a BindUnmarshaler
// func bindUnmarshaler(field reflect.Value) (BindUnmarshaler, bool) {
// 	ptr := reflect.New(field.Type())
// 	if ptr.CanInterface() {
// 		iface := ptr.Interface()
// 		if unmarshaler, ok := iface.(BindUnmarshaler); ok {
// 			return unmarshaler, ok
// 		}
// 	}
// 	return nil, false
// }
