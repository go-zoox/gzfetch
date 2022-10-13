package main

import (
	"github.com/go-zoox/cli"
	"github.com/go-zoox/gzfetch/commands/get"
	"github.com/go-zoox/gzfetch/commands/post"
	"github.com/go-zoox/gzfetch/commands/request"
)

func main() {
	app := cli.NewMultipleProgram(&cli.MultipleProgramConfig{
		Name:  "gzfetch",
		Usage: "Simple and powerful request cli, alternative to curl",
	})

	get.Create(app)
	post.Create(app)

	request.Create(app)

	app.Run()
}
