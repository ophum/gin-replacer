package ginreplacer

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type IgnoreFunc func(ctx *gin.Context) bool

func defaultIgnoreFunc(ctx *gin.Context) bool {
	return false
}

type Config struct {
	IgnoreFunc IgnoreFunc
	Replacer   *strings.Replacer
}

type replacerWriter struct {
	gin.ResponseWriter
	replacer *strings.Replacer
}

func (w *replacerWriter) Write(b []byte) (int, error) {
	oldLen := len(b)
	b = []byte(w.replacer.Replace(string(b)))
	newLen := len(b)
	if err := w.adjustContentLength(oldLen, newLen); err != nil {
		return 0, err
	}
	return w.ResponseWriter.Write(b)
}

func (w *replacerWriter) adjustContentLength(oldLen, newLen int) error {
	contentLength, err := strconv.ParseInt(w.ResponseWriter.Header().Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}
	adjustContentLength := int(contentLength) + (newLen - oldLen)
	w.ResponseWriter.Header().Set("Content-Length", fmt.Sprint(adjustContentLength))
	return nil
}

func New(config *Config) gin.HandlerFunc {
	if config.IgnoreFunc == nil {
		config.IgnoreFunc = defaultIgnoreFunc
	}

	return func(ctx *gin.Context) {
		if !config.IgnoreFunc(ctx) {
			w := &replacerWriter{
				ResponseWriter: ctx.Writer,
				replacer:       config.Replacer,
			}
			ctx.Writer = w
		}
		ctx.Next()
	}
}
