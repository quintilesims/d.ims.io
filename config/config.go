package config

const (
	ENVVAR_PORT             = "DIMSIO_PORT"
	ENVVAR_DEBUG            = "DIMSIO_DEBUG"
	ENVVAR_AWS_ACCESS_KEY   = "DIMSIO_AWS_ACCESS_KEY"
	ENVVAR_AWS_SECRET_KEY   = "DIMSIO_AWS_SECRET_KEY"
	ENVVAR_AWS_REGION       = "DIMSIO_AWS_REGION"
	ENVVAR_DYNAMO_TABLE     = "DIMSIO_DYNAMO_TABLE"
	ENVVAR_ACCOUNT_TABLE    = "DIMSIO_ACCOUNT_TABLE"
	ENVVAR_AUTH0_DOMAIN     = "DIMSIO_AUTH0_DOMAIN"
	ENVVAR_AUTH0_CLIENT_ID  = "DIMSIO_AUTH0_CLIENT_ID"
	ENVVAR_AUTH0_CONNECTION = "DIMSIO_AUTH0_CONNECTION"
)

const (
	DEFAULT_PORT          = "80"
	DEFAULT_AWS_REGION    = "us-west-2"
	DEFAULT_DYNAMO_TABLE  = "d.ims.io"
	DEFAULT_ACCOUNT_TABLE = "d.ims.io.account"
	DEFAULT_AUTH0_DOMAIN  = "https://imshealth.auth0.com"
)
