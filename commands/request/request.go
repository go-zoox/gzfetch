package request

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/object"
	"github.com/go-zoox/encoding/yaml"
	"github.com/go-zoox/fetch"
	"github.com/go-zoox/fs"
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

			cfgBytes, err := fs.ReadFile(config)
			if err != nil {
				return err
			}

			// // apply env
			// // replace ${KEY} => env value
			// reEnvNameWithBracket := regexp.MustCompile(`\$\{[^\s]+\}`)
			// envdCfgBytes := reEnvNameWithBracket.ReplaceAllFunc(cfgBytes, func(matched []byte) []byte {
			// 	envKey := string(matched[2 : len(matched)-1])

			// 	// fmt.Println("key", envKey, " value:", os.Getenv(envKey))
			// 	return []byte(os.Getenv(envKey))
			// })

			// // replace $KEY => env value
			// reEnvNameWithoutBracket := regexp.MustCompile(`\$([^\s]+)`)
			// envdCfgBytes = reEnvNameWithoutBracket.ReplaceAllFunc(envdCfgBytes, func(matched []byte) []byte {
			// 	envKey := string(matched[1:])

			// 	// fmt.Println("key", envKey, " value:", os.Getenv(envKey))
			// 	return []byte(os.Getenv(envKey))
			// })

			cfg := fetch.Config{}
			if err := yaml.Decode(cfgBytes, &cfg); err != nil {
				return err
			}

			// apply env
			if cfg.Headers != nil {
				for k, s := range cfg.Headers {
					if len(s) > 1 && s[0] == '$' {
						if len(s) > 2 && s[1] == '{' {
							cfg.Headers[k] = os.Getenv(s[2 : len(s)-1])
						} else {
							cfg.Headers[k] = os.Getenv(s[1:])
						}
					}
				}
			}

			if cfg.Params != nil {
				for k, s := range cfg.Params {
					if len(s) > 1 && s[0] == '$' {
						if len(s) > 2 && s[1] == '{' {
							cfg.Params[k] = os.Getenv(s[2 : len(s)-1])
						} else {
							cfg.Params[k] = os.Getenv(s[1:])
						}
					}
				}
			}

			if cfg.Query != nil {
				for k, s := range cfg.Query {
					if len(s) > 1 && s[0] == '$' {
						if len(s) > 2 && s[1] == '{' {
							cfg.Query[k] = os.Getenv(s[2 : len(s)-1])
						} else {
							cfg.Query[k] = os.Getenv(s[1:])
						}
					}
				}
			}

			if cfg.Body != nil {
				if body, ok := cfg.Body.(map[string]interface{}); ok {
					for k, v := range body {
						if s, ok := v.(string); ok {
							if len(s) > 1 && s[0] == '$' {
								if len(s) > 2 && s[1] == '{' {
									body[k] = os.Getenv(s[2 : len(s)-1])
								} else {
									body[k] = os.Getenv(s[1:])
								}
							}
						}
					}

					cfg.Body = body
				}
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
