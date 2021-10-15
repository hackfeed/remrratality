package testing

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func CreateGinContext(keys map[string]interface{}, body interface{}, params []gin.Param) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	jsonBytes, _ := json.Marshal(body)
	req := &http.Request{
		Body: io.NopCloser(bytes.NewBuffer(jsonBytes)),
	}

	for k, v := range keys {
		c.Set(k, v)
	}
	c.Request = req
	c.Params = params

	return c, w
}
