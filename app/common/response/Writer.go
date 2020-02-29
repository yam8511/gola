package response

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

// BodyWriter 資料Wrtier
type BodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *BodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

// BodyResponse 自訂回應
type BodyResponse struct {
	ErrorCode string `json:"error_code"`
	ErrorText string `json:"error_text"`
	Data      []byte `json:"data,omitempty"`
}

// NewBodyWriter 新增一個Writer
func NewBodyWriter(c *gin.Context) *BodyWriter {
	w, ok := c.Get("has_new_body_writer")
	if ok {
		return w.(*BodyWriter)
	}

	bw := &BodyWriter{
		Body:           new(bytes.Buffer),
		ResponseWriter: c.Writer,
	}

	c.Writer = bw
	c.Set("has_new_body_writer", bw)
	return bw
}
