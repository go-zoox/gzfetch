package post

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-zoox/cli"
	"github.com/go-zoox/core-utils/object"
	"github.com/go-zoox/fetch"
)

// Create creates a post command.
func Create(app *cli.MultipleProgram) {
	app.Register("post", &cli.Command{
		Name:  "post",
		Usage: "http post request",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "headers",
				Usage: "set request headers, example: authorization=token,x=2",
			},
			&cli.StringFlag{
				Name:  "params",
				Usage: "set request params, example: q1=x,q2=y",
			},
			&cli.StringFlag{
				Name:  "query",
				Usage: "set request query, example: q1=x&q2=y",
			},
			&cli.StringFlag{
				Name:  "body",
				Usage: "set request body, example: { \"x\": 1, \"y\": 2 }",
			},
			&cli.StringFlag{
				Name:    "pick",
				Usage:   "get value by key of the response, example: --pick headers.user-agent",
				Aliases: []string{"p"},
			},
		},
		Action: func(ctx *cli.Context) error {
			url := ctx.Args().Get(0)
			headersX := ctx.String("headers")
			paramsX := ctx.String("params")
			queryX := ctx.String("query")
			bodyX := ctx.String("body")

			pick := ctx.String("pick")

			if url == "" {
				return fmt.Errorf("url is required")
			}

			headers := map[string]string{}
			if headersX != "" {
				headersSlice := strings.Split(headersX, ",")
				for _, item := range headersSlice {
					kv := strings.Split(item, "=")
					headers[kv[0]] = kv[1]
				}
			}

			params := map[string]string{}
			if paramsX != "" {
				paramsSlice := strings.Split(paramsX, ",")
				for _, item := range paramsSlice {
					kv := strings.Split(item, "=")
					params[kv[0]] = kv[1]
				}
			}

			query := map[string]string{}
			if queryX != "" {
				querySlice := strings.Split(queryX, "&")
				for _, item := range querySlice {
					kv := strings.Split(item, "=")
					query[kv[0]] = kv[1]
				}
			}

			body := map[string]any{}
			if bodyX != "" {
				err := json.Unmarshal([]byte(bodyX), &body)
				if err != nil {
					return err
				}
			}

			response, err := fetch.Post(url, &fetch.Config{
				Headers: headers,
				Params:  params,
				Query:   query,
				Body:    body,
			})
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
//	  go run . post --headers authorization=tokeb,x=2 --query 'x=1&y=2' --body '{"x":1}' https://httpbin.zcorky.com/post
