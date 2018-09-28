package ironic

import (
	go_api "github.com/micro/go-api/proto"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

const (
	TestJSON       = `{"id": 1, "title": "title 1"}`
	TestForm       = `id=1&title=title 1`
	invalidContent = "invalid content"
)

type TestStruct struct {
	ID    int    `json:"id" xml:"id" form:"id" query:"id"`
	Title string `json:"title" xml:"title" form:"title" query:"title"`
}

func newReq(method string, body string, args ...map[string]*go_api.Pair) *go_api.Request {
	res := &go_api.Request{
		Method: method,
		Body:   body,
		Header: make(map[string]*go_api.Pair),
		Get:    make(map[string]*go_api.Pair),
		Post:   make(map[string]*go_api.Pair),
	}
	argsLen := len(args)

	if argsLen > 0 {

		switch argsLen {
		case 1:
			res.Header = args[0]
		case 2:
			res.Header = args[0]
			res.Get = args[1]
		case 3:
			res.Header = args[0]
			res.Get = args[1]
			res.Post = args[2]
		default:
			res.Header = args[0]
			res.Get = args[1]
			res.Post = args[2]
		}

	}

	return res
}

func TestBindJSON(t *testing.T) {
	testBindOK(t, newReq(POST, TestJSON), MIMEApplicationJSON)
	testBindError(t, newReq(POST, invalidContent), MIMEApplicationJSON)
}

func TestBindForm(t *testing.T) {
	post := make(map[string]*go_api.Pair)
	post["id"] = &go_api.Pair{Values: []string{"1"}, Key: "id"}
	post["title"] = &go_api.Pair{Values: []string{"title 1"}, Key: "title"}

	testBindOK(t, newReq(POST, TestForm, make(map[string]*go_api.Pair), make(map[string]*go_api.Pair), post), MIMEApplicationForm)
	testBindError(t, newReq(POST, TestForm, make(map[string]*go_api.Pair), make(map[string]*go_api.Pair), post), MIMEApplicationForm)
}

func testBindOK(t *testing.T, req *go_api.Request, ctype string) {
	ts := new(TestStruct)
	binder := new(DefaultBinder)
	req.Header[HeaderContentType] = &go_api.Pair{
		Key:    HeaderContentType,
		Values: []string{ctype},
	}
	err := binder.Bind(req, &ts)
	t.Log(err)
	if err == nil {
		assert.Equal(t, 1, ts.ID)
		assert.Equal(t, "title 1", ts.Title)
	}
}

func testBindError(t *testing.T, req *go_api.Request, ctype string) {
	ts := new(TestStruct)
	binder := new(DefaultBinder)
	req.Header[HeaderContentType] = &go_api.Pair{
		Key:    HeaderContentType,
		Values: []string{ctype},
	}

	err := binder.Bind(req, &ts)

	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON), strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML),
		strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		if err != nil {
			assert.NotEqual(t, "", err.Error())
		}
	default:
		if err != nil {
			assert.Equal(t, "UnsupportMediaType", err.Error())
		}
	}
}
