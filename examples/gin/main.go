package main

import (
	"examples/gin/models"

	"github.com/gin-gonic/gin"
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)

func main() {
	g := gin.Default()
	go_ssr.Init(config.Config{
		FrontendDir:        "./frontend/src",
		GeneratedTypesPath: "./frontend/src/generated.d.ts",
		PropsStructsPath:   "./models/props.go",
	})

	// props := models.IndexRouteProps{
	// 	InitialCount: 0,
	// 	Message:      "Hello from Go SSR!",
	// }
	g.GET("/", func(c *gin.Context) {
		react_renderer.RenderRoute(c, react_renderer.Config{
			File: "Home.tsx",
			MetaTags: map[string]string{
				"title":       "My app",
				"og:title":    "My app",
				"description": "Hello world!",
			},
			Props: &models.IndexRouteProps{
				InitialCount: 0,
				Message:      "Hello from Go SSR!",
			},
		})
	})
	g.Run()
}

// func testFunc(props interface{}) {
// 	converter := typescriptify.New().Add(props)
// 	err := converter.ConvertToFile("models.ts")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// }
