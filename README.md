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

<!--
# 💡 Overview -->

Go-SSR was developed due to a lack of an existing product in the Go ecosystem that made it easy to build full-stack React apps. At the time, most Go web app projects were either built with a static React frontend with lots of client-side logic or html templates. I envisioned creating a new solution that would allow you to create full-stack Go apps with React but with logic being moved to the server and being able to pass that logic down with type-safe props. This project was inspired by [Remix](https://remix.run/) and [Next.JS](https://nextjs.org/), but aims to be a plugin and not a framework.

# 📜 Features

- Lightning fast compiling
- Auto generated Typescript structs for props
- Smart build caching system
- Hot reloading
- Simple error reporting
- Production optimized
- Drop in to any existing Go web server

<!-- _View more examples [here](github.com/natewong1313/go-react-ssr/examples)_ -->

# 🛠️ Getting Started

Go-SSR can either be installed by cloning a template or simply installing it on your own. It is reccomended to take a look at the [examples](/examples) folder to see how projects are structured.

## ⚡️ Using pre-configured templates

Templates have been pre-configured to be installed from the console using the [`gonew`](https://go.dev/blog/gonew) command. View more info in the [examples](/examples) folder

### Gin

```console
$ gonew github.com/natewong1313/go-react-ssr/examples/gin example.com/gin
```

### Fiber

```console
$ gonew github.com/natewong1313/go-react-ssr/examples/fiber example.com/fiber
```

## 📝 Add to existing web server

First, install the module

```console
$ go get -u github.com/natewong1313/go-react-ssr
```

Then, add imports into your main file

```go
import (
	...
	go_ssr "github.com/natewong1313/go-react-ssr"
	"github.com/natewong1313/go-react-ssr/pkg/config"
	"github.com/natewong1313/go-react-ssr/pkg/react_renderer"
)
```

In your main function, initialize the plugin

```go
go_ssr.Init(config.Config{
	FrontendDir:        "./frontend/src", // The React source files path
	GeneratedTypesPath: "./frontend/src/generated.d.ts", // Where the generated prop types will be created at
	PropsStructsPath:   "./models/props.go", // Where the Go structs for your prop types are located
})
```

Once the plugin has been initialized, you can call the `react_renderer.RenderRoute` function to compile your React file to a string

```go
g.GET("/", func(c *gin.Context) {
	c.Writer.Write(react_renderer.RenderRoute(react_renderer.Config{
		File:  "Home.tsx", // The file name, located in FrontendDir
		Title: "My example app", // Set the app title
		MetaTags: map[string]string{ // Configure meta tags
			"og:title":    "My example app",
			"description": "Example description",
		},
		Props: &models.IndexRouteProps{ // Pass in your props struct
			InitialCount: 0,
		},
	}))
})
```
