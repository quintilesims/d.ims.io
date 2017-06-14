package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/quintilesims/d.ims.io/aws"
	"github.com/quintilesims/d.ims.io/config"
	"github.com/quintilesims/d.ims.io/controllers"
	"github.com/quintilesims/d.ims.io/token"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

const (
	SWAGGER_URL     = "/api/"
	SWAGGER_UI_PATH = "static/swagger-ui/dist"
)

func serveSwaggerUI(w http.ResponseWriter, r *http.Request) {
	dir := http.Dir(SWAGGER_UI_PATH)
	fileServer := http.FileServer(dir)
	http.StripPrefix(SWAGGER_URL, fileServer).ServeHTTP(w, r)
}

var Version string

func main() {
	if Version == "" {
		Version = "unset/developer"
	}

	app := cli.NewApp()
	app.Name = "d.ims.io"
	app.Version = Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "p, port",
			Value:  config.DEFAULT_PORT,
			EnvVar: config.ENVVAR_PORT,
		},
		cli.StringFlag{
			Name:   "aws-access-key",
			EnvVar: config.ENVVAR_AWS_ACCESS_KEY,
		},
		cli.StringFlag{
			Name:   "aws-secret-key",
			EnvVar: config.ENVVAR_AWS_SECRET_KEY,
		},
		cli.StringFlag{
			Name:   "aws-region",
			Value:  config.DEFAULT_AWS_REGION,
			EnvVar: config.ENVVAR_AWS_REGION,
		},
		cli.StringFlag{
			Name:   "swagger-host",
			Value:  config.DEFAULT_SWAGGER_HOST,
			EnvVar: config.ENVVAR_SWAGGER_HOST,
		},
		cli.StringFlag{
			Name:   "dynamo-table",
			Value:  config.DEFAULT_DYNAMO_TABLE,
			EnvVar: config.ENVVAR_DYNAMO_TABLE,
		},
	}

	app.Action = func(c *cli.Context) error {
		if err := assertConfig(c); err != nil {
			return err
		}

		rand.Seed(time.Now().Unix())

		config := defaults.Get().Config
		staticCreds := credentials.NewStaticCredentials(c.String("aws-access-key"), c.String("aws-secret-key"), "")
		config.WithCredentials(staticCreds)
		config.WithRegion(c.String("aws-region"))
		aws := aws.NewProvider(config)

		dynamo := token.NewDynamoAuth(c.String("dynamo-table"), aws)

		routes := controllers.NewRootController().Routes()
		routes = append(routes, controllers.NewTokenController(nil, dynamo).Routes()...)
		routes = append(routes, controllers.NewSwaggerController(c.String("swagger-host")).Routes()...)
		fb := fireball.NewApp(routes)

		port := fmt.Sprintf(":%s", c.String("port"))
		log.Printf("Running on port %s\n", port)
		http.Handle("/", fb)

		http.HandleFunc(SWAGGER_URL, serveSwaggerUI)
		return http.ListenAndServe(port, nil)
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func assertConfig(c *cli.Context) error {
	vars := map[string]error{
		"aws-access-key": fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", config.ENVVAR_AWS_ACCESS_KEY),
		"aws-secret-key": fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", config.ENVVAR_AWS_SECRET_KEY),
		"aws-region":     fmt.Errorf("AWS Region not set! (EnvVar: %s)", config.ENVVAR_AWS_REGION),
		"dynamo-table":   fmt.Errorf("Dynamo Table not set! (EnvVar: %s)", config.ENVVAR_DYNAMO_TABLE),
	}

	for name, err := range vars {
		if c.String(name) == "" {
			return err
		}
	}

	return nil
}
