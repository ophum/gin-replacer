package ginreplacer_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	ginreplacer "github.com/ophum/gin-replacer"
	"github.com/stretchr/testify/assert"
)

func TestReplace(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name     string
		replacer *strings.Replacer
		original string
		want     string
	}{
		{
			name:     "old to new",
			replacer: strings.NewReplacer("old", "new"),
			original: "old, threshold to threshnew...",
			want:     "new, threshnew to threshnew...",
		},
		{
			name: "old to new, new to OLD, original old is not replace OLD",
			replacer: strings.NewReplacer(
				"old", "new",
				"new", "OLD",
			),
			original: "old to new",
			want:     "new to OLD",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, e := gin.CreateTestContext(w)
			e.Use(ginreplacer.New(&ginreplacer.Config{
				Replacer: tt.replacer,
			}))
			e.GET("", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, tt.original)
			})

			req := httptest.NewRequest(http.MethodGet, "http://localhsot/", nil)
			e.ServeHTTP(w, req)

			actual := w.Body.String()
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestIgnoreFunc(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name        string
		path        string
		wantReplace bool
	}{
		{
			name:        "index.js",
			path:        "index.js",
			wantReplace: true,
		},
		{
			name:        "index.html",
			path:        "index.html",
			wantReplace: true,
		},
		{
			name:        "index.css",
			path:        "index.css",
			wantReplace: false,
		},
		{
			name:        "index.htm",
			path:        "index.htm",
			wantReplace: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			_, e := gin.CreateTestContext(w)
			e.Use(ginreplacer.New(&ginreplacer.Config{
				IgnoreFunc: func(ctx *gin.Context) bool {
					ext := filepath.Ext(ctx.Request.URL.Path)
					switch ext {
					case ".js", ".html":
						return false
					default:
						return true
					}
				},
				Replacer: strings.NewReplacer("old", "new"),
			}))
			e.Any("/*path", func(ctx *gin.Context) {
				ctx.String(http.StatusOK, "old")
			})
			u := url.URL{
				Scheme: "http",
				Host:   "localhost",
				Path:   tt.path,
			}

			req := httptest.NewRequest(http.MethodGet, u.String(), nil)
			e.ServeHTTP(w, req)

			actual := w.Body.String()
			if tt.wantReplace {
				assert.Equal(t, "new", actual)
			} else {
				assert.Equal(t, "old", actual)
			}
		})
	}
}
