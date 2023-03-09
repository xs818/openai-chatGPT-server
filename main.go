package main

import (
	"context"
	"github.com/openai-chatGPT-server/config"
	"github.com/openai-chatGPT-server/logic"
	"github.com/openai-chatGPT-server/middleware"
	"github.com/openai-chatGPT-server/router"
	"github.com/openai-chatGPT-server/server"
	"github.com/urfave/cli/v2"
	"go.uber.org/fx"
	"os"
)
func main()  {
	app := cli.NewApp()
	app.Flags = append(app.Flags, []cli.Flag{
		&cli.StringFlag{
			Name:     "env",
			Usage:    "请输入运行环境:dev:开发环境 eg:--env dev",
			Value: "dev",
		},
	}...)

	app.Action = func(context *cli.Context) error {
		InitApp(context)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func InitApp(c *cli.Context)  {
	fxApp := fx.New(
		fx.Provide(config.NewConfig),
		fx.Provide(config.LoadEnv),
		fx.Supply(c),
		fx.Provide(router.NewRouter),
		fx.Provide(middleware.NewMiddleware),
		fx.Provide(logic.NewChat),
		fx.Provide(server.NewHTTPServer),
		fx.Invoke(func(server *server.HTTPServer) {
			server.Run()
		}),
	)

	if err := fxApp.Start(context.Background()); err != nil {
		panic(err)
	}

}






