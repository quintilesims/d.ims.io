package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/quintilesims/d.ims.io/config"
	"github.com/quintilesims/d.ims.io/controllers"
	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"github.com/quintilesims/d.ims.io/router"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
	"log"
	"net/http"
	"os"
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

func main() {
	app := cli.NewApp()
	app.Name = "d.ims.io"
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
		if err := validateConfig(c); err != nil {
			return err
		}

		// todo: add aws credentials to awsConfig
		awsConfig := defaults.Get().Config
		session := session.New(awsConfig)
		dynamodb := dynamodb.New(session)
		ecr := ecr.New(session)

		tokenManager := auth.NewDynamoTokenManager("todo = table name", dynamodb)
		auth0Manager := auth.NewAuth0Manager("todo - auth0 endpoint", "todo - auth0 token")
		proxy := proxy.NewECRProxy("todo - ecr endpoint")

		rootController := controllers.NewRootController()
		repositoryController := controllers.NewRepositoryController(ecr)
		tokenController := controllers.NewTokenController(tokenManager)
		proxyController := controllers.NewProxyController(ecr, proxy)
		swaggerController := controllers.NewSwaggerController(c.String("swagger-host"))

		compositeAuthenticator := auth.NewCompositeAuthenticator(tokenManager, auth0Manager)

		routes := rootController.Routes()
		routes = append(routes, repositoryController.Routes()...)
		routes = append(routes, tokenController.Routes()...)
		routes = append(routes, swaggerController.Routes()...)
		routes = fireball.Decorate(routes,
			fireball.LogDecorator(),
			controllers.AuthDecorator(compositeAuthenticator))

		fb := fireball.NewApp(routes)
		fb.Router = router.NewRouter(routes, proxyController)

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

func validateConfig(c *cli.Context) error {
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
