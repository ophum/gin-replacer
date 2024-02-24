package main

import (
	"embed"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	ginreplacer "github.com/ophum/gin-replacer"
)

//go:embed root/*
var fs embed.FS

func main() {
	router := gin.Default()

	router.Use(ginreplacer.New(&ginreplacer.Config{
		IgnoreFunc: func(ctx *gin.Context) bool {
			return filepath.Ext(ctx.Request.URL.Path) != ".js"
		},
		Replacer: strings.NewReplacer(
			"%APIBASEURL%", "http://localhost:8080/api",
		),
	}))
	router.Use(static.Serve("", static.EmbedFolder(fs, "root")))

	router.GET("/api/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := router.Run(); err != nil {
		panic(err)
	}
}
