package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/quintilesims/d.ims.io/auth"
	"github.com/quintilesims/d.ims.io/config"
	"github.com/quintilesims/d.ims.io/controllers"
	"github.com/quintilesims/d.ims.io/controllers/proxy"
	"github.com/quintilesims/d.ims.io/logging"
	"github.com/quintilesims/d.ims.io/router"
	"github.com/urfave/cli"
	"github.com/zpatrick/fireball"
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
		cli.BoolFlag{
			Name:   "d, debug",
			EnvVar: config.ENVVAR_DEBUG,
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
		cli.StringFlag{
			Name:   "registry-endpoint",
			EnvVar: config.ENVVAR_REGISTRY_ENDPOINT,
		},
		cli.StringFlag{
			Name:   "auth0-domain",
			Value:  config.DEFAULT_AUTH0_DOMAIN,
			EnvVar: config.ENVVAR_AUTH0_DOMAIN,
		},
		cli.StringFlag{
			Name:   "auth0-client-id",
			EnvVar: config.ENVVAR_AUTH0_CLIENT_ID,
		},
		cli.StringFlag{
			Name:   "auth0-connection",
			EnvVar: config.ENVVAR_AUTH0_CONNECTION,
		},
	}

	app.Before = func(c *cli.Context) error {
		if err := validateConfig(c); err != nil {
			return err
		}

		log.SetOutput(logging.NewLogWriter(c.Bool("debug")))
		return nil
	}

	app.Action = func(c *cli.Context) error {
		session := getAWSSession(c)
		dynamodb := dynamodb.New(session)
		ecr := ecr.New(session)

		tokenManager := auth.NewDynamoTokenManager(c.String("dynamo-table"), dynamodb)
		// todo: terraform module will need https in front of domain
		auth0Manager := auth.NewAuth0Manager(c.String("auth0-domain"), c.String("auth0-client-id"), c.String("auth0-connection"))
		proxy := proxy.NewECRProxy(c.String("registry-endpoint"))

		rootController := controllers.NewRootController()
		repositoryController := controllers.NewRepositoryController(ecr)
		tokenController := controllers.NewTokenController(tokenManager)
		proxyController := controllers.NewProxyController(ecr, proxy)
		swaggerController := controllers.NewSwaggerController(c.String("swagger-host"))

		authenticator := auth.NewCompositeAuthenticator(tokenManager, auth0Manager)

		routes := rootController.Routes()
		routes = append(routes, repositoryController.Routes()...)
		routes = append(routes, tokenController.Routes()...)
		routes = append(routes, swaggerController.Routes()...)
		routes = fireball.Decorate(routes,
			fireball.LogDecorator(),
			controllers.AuthDecorator(authenticator))

		fb := fireball.NewApp(routes)

		// decorate proxy handler with auth
		doProxy := controllers.AuthDecorator(authenticator)(proxyController.DoProxy)
		fb.Router = router.NewRouter(routes, doProxy)

		port := fmt.Sprintf(":%s", c.String("port"))
		log.Printf("[INFO] Running on port %s\n", port)
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
		"aws-access-key":    fmt.Errorf("AWS Access Key not set! (EnvVar: %s)", config.ENVVAR_AWS_ACCESS_KEY),
		"aws-secret-key":    fmt.Errorf("AWS Secret Key not set! (EnvVar: %s)", config.ENVVAR_AWS_SECRET_KEY),
		"aws-region":        fmt.Errorf("AWS Region not set! (EnvVar: %s)", config.ENVVAR_AWS_REGION),
		"dynamo-table":      fmt.Errorf("Dynamo Table not set! (EnvVar: %s)", config.ENVVAR_DYNAMO_TABLE),
		"registry-endpoint": fmt.Errorf("Registry endpoint not set! (EnvVar: %s)", config.ENVVAR_REGISTRY_ENDPOINT),
		"auth0-domain":      fmt.Errorf("Auth0 Domain not set! (EnvVar: %s)", config.ENVVAR_AUTH0_DOMAIN),
		"auth0-client-id":   fmt.Errorf("Auth0 Client ID not set! (EnvVar: %s)", config.ENVVAR_AUTH0_CLIENT_ID),
		"auth0-connection":  fmt.Errorf("Auth0 Connection not set! (EnvVar: %s)", config.ENVVAR_AUTH0_CONNECTION),
	}

	for name, err := range vars {
		if c.String(name) == "" {
			return err
		}
	}

	return nil
}

func getAWSSession(c *cli.Context) *session.Session {
	config := defaults.Get().Config
	creds := credentials.NewStaticCredentials(c.String("aws-access-key"), c.String("aws-secret-key"), "")
	config.WithCredentials(creds)
	config.WithRegion(c.String("aws-region"))
	return session.New(config)
}
