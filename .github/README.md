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
- Hot reloading
- Simple error reporting
- Production optimized
- Drop in to any existing Go web server

<!-- _View more examples [here](github.com/natewong1313/go-react_old-ssr/examples)_ -->

# 🛠️ Getting Started

Go-SSR was designed with the idea of being dead simple to install. Below are 2 easy ways of setting it up:

## ⚡️ Using the CLI tool

<img src="https://i.imgur.com/mygp5BT.png" height="400" />

The easiest way to get a project up and running is by using the command line tool. Install it with the following command

```console
$ go install github.com/natewong1313/go-react-ssr/gossr-cli@latest
```

Then you can call the following command to create a project

```console
$ gossr-cli create
```

You'll be prompted the path to place the project, what web framework you want to use, and whether or not you want to use Tailwind

## 📝 Add to existing web server

To add Go-SSR to an existing Go web server, take a look at the [examples](/examples) folder to get an idea of what a project looks like. In general, you'll want to follow these commands:

```console
$ go get -u github.com/natewong1313/go-react-ssr
```

Then, add imports into your main file

```go
import (
	...
	gossr "github.com/natewong1313/go-react-ssr"
)
```

In your main function, initialize the plugin. Create a folder for your structs that hold your props to go, which is called `models` in the below example. You'll also want to create a folder for your React code (called `frontend` in this example) inside your project and specifiy the paths in the config. You may want to clone the [example folder](/examples/frontend/) and use that.

```go
engine, err := gossr.New(gossr.Config{
    AssetRoute:         "/assets",
    FrontendDir:        "./frontend/src",
    GeneratedTypesPath: "./frontend/src/generated.d.ts",
    PropsStructsPath:   "./models/props.go",
})
```

Once the plugin has been initialized, you can call the `engine.RenderRoute` function to compile your React file to a string

```go
g.GET("/", func(c *gin.Context) {
	renderedResponse := engine.RenderRoute(gossr.RenderConfig{
		File:  "Home.tsx", 
		Title: "Example app", 
		MetaTags: map[string]string{
			"og:title":    "Example app", 
			"description": "Hello world!",
		}, 
		Props: &models.IndexRouteProps{
			InitialCount: rand.Intn(100),
		},
	})
	c.Writer.Write(renderedResponse)
})
```

# 🚀 Deploying to production

All of the examples come with a Dockerfile that you can use to deploy to production. You can also use the [gossr-cli](#-using-the-cli-tool) to create a project with a Dockerfile.
Below is an example Dockerfile

```Dockerfile
FROM golang:1.21-alpine as build-backend
RUN apk add git
ADD .. /build
WORKDIR /build

RUN go mod download
RUN go get -u github.com/natewong1313/go-react-ssr
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -X main.APP_ENV=production" -a -o main


FROM node:16-alpine as build-frontend

ADD ./frontend /frontend
WORKDIR /frontend

RUN npm install

# if tailwind is enabled, use "FROM node:16-alpine" instead
FROM alpine:latest
COPY --from=build-backend /build/main ./app/main
COPY --from=build-frontend /frontend ./app/frontend

WORKDIR /app
RUN chmod +x ./main
EXPOSE 8080
CMD ["./main"]
```

Go SSR has been tested and deployed on the following platforms:

- [Fly.io](https://fly.io/) - [example app](https://sparkling-smoke-7627.fly.dev/)
- [Render](https://render.com/) - [example app](https://my-gossr-test.onrender.com/)
- [Hop.io](https://hop.io/) - [example app](https://my-gossr-test.hop.sh/)
