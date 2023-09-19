<!-- # Go React SSR -->

<!-- Build Go powered React web apps with end to end type-safety -->
<img src="https://i.imgur.com/zrKSrny.png" height="72">

---

<p>
    <a href="https://goreportcard.com/report/github.com/natewong1313/go-react-ssr"><img src="https://goreportcard.com/badge/github.com/natewong1313/go-react-ssr" alt="Go Report"></a>
    <a href="https://pkg.go.dev/github.com/natewong1313/go-react-ssr?tab=doc"><img src="http://img.shields.io/badge/GoDoc-Reference-blue.svg" alt="GoDoc"></a>
    <a href="https://github.com/natewong1313/go-react-ssr/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT%202.0-blue.svg" alt="MIT License"></a>
</p>

Go-SSR is a drop in plugin to **any** existing Go web framework to allow **server rendering** [React](https://react.dev/). It's powered by [esbuild](https://esbuild.github.io/) and allows for passing props from Go to React with **type safety**.

# Quickstart

```go
package main

import (
  "example.com/gin/models"
  "github.com/gin-gonic/gin"
  go_ssr "github.com/natewong1313/go-react-ssr"
  "github.com/natewong1313/go-react-ssr/pkg/config"
  "github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)

func main() {
  g := gin.Default()

  go_ssr.Init(config.Config{
    PropsStructsPath:  "./models/props.go",
  })

  g.GET("/", func(c *gin.Context) {
    react_renderer.RenderRoute(c, react_renderer.Config{
      File:  "Home.tsx",
      Props: &models.IndexRouteProps{
        InitialCount: 0,
      },
    })
  })
  g.Run()
}


```
