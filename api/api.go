package api

import (
	go_api "github.com/micro/go-api/proto"
	"strings"
)

type Helper struct{}

func (h *Helper) GetTokenFromRequest(r *go_api.Request) string {
	header := r.GetHeader()

	// fmt.Println(header)
	tokenString := ""

	if accessToken, ok := header["U-Access-Token"]; ok {
		val := accessToken.GetValues()
		tokenString = val[len(val)-1]
	}

	// 标准的jwt头部
	if authorization, ok := header["Authorization"]; ok {
		val := authorization.GetValues()
		t := val[len(val)-1]
		tokenString = strings.TrimSpace(strings.Replace(t, "Bearer ", "", 1))
	}
	// resp.StatusCode = 200
	// resp.Body = "Hello World"
	return tokenString
}
