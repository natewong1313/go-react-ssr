package main

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/natewong1313/go-react-ssr/gossr-cli/cmd"
	_ "github.com/natewong1313/go-react-ssr/gossr-cli/cmd/create"
)

func main() {
	art := figure.NewFigure("Go - SSR CLI", "slant", true)
	art.Print()
	fmt.Println()
	cmd.Execute()
}
