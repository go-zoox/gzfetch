package request

import (
	"encoding/json"
	"fmt"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/object"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/fs/type/yaml"
)

// Create creates a request command.
func Create(app *cli.MultipleProgram) {
	app.Register("request", &cli.Command{
		Name:  "request",
		Usage: "http request",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "request config file",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "pick",
				Usage:   "get value by key of the response, example: --pick headers.user-agent",
				Aliases: []string{"p"},
			},
		},
		Action: func(ctx *cli.Context) error {
			config := ctx.String("config")
			pick := ctx.String("pick")

			if config == "" {
				return fmt.Errorf("config is required")
			}

			cfg := fetch.Config{}
			if err := yaml.Read(config, &cfg); err != nil {
				return err
			}

			response, err := fetch.New(&cfg).Execute()
			if err != nil {
				return err
			}

			jd, err := response.JSON()
			if err != nil || jd == "null" {
				fmt.Println(response.String())
				return nil
			}

			if pick != "" {
				data := map[string]any{}
				err := json.Unmarshal([]byte(jd), &data)
				if err != nil {
					return err
				}

				fmt.Println(object.Get(data, pick))
				return nil
			}

			fmt.Println(jd)
			return nil
		},
	})
}

// test:
//	  go run . request -c $PWD/commands/request/request_config.test.yaml
