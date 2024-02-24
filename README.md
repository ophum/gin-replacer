# gin-replacer gin's middleware

Gin middleware/handler to replace response body.

## Usage

### Start using it

Download and install it:

```
go get github.com/ophum/gin-replacer
```

Import it in your code:

```go
import "github.com/ophum/gin-replacer
```

### example

```go
package main

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	ginreplacer "github.com/ophum/gin-replacer"
)

func main() {
	router := gin.Default()

	router.Use(ginreplacer.New(&ginreplacer.Config{
        // when true: not replace
        // if IgnoreFunc is nil, using defaultIgnoreFunc. it is always return false.
		IgnoreFunc: func(ctx *gin.Context) bool {
			return filepath.Ext(ctx.Request.URL.Path) != ".html"
		},
		Replacer: strings.NewReplacer(
			"OLD", "NEWCONTENT",
		),
	}))

	router.GET("/test.html", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OLD")
	})

	router.Run()
}
```

```
$ curl localhost:8080/test.html
NEWCONTENT
```
