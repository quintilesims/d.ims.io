# D.IMS.IO

[![MIT License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/quintilesims/d.ims.io/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/quintilesims/d.ims.io)](https://goreportcard.com/report/github.com/quintilesims/d.ims.io)
[![Go Doc](https://godoc.org/github.com/quintilesims/d.ims.io?status.svg)](https://godoc.org/github.com/quintilesims/d.ims.io)


## Overview
This repository serves as a web front-end for [EC2 Container Registries](https://aws.amazon.com/ecr/).
It provides custom authentication and an API to create and manage Docker repositories. 

## Repository Management
You must create a new repository for each image you want to host.
Creating new repositories can be done through the [Swagger UI](https://d.ims.io/api/?url=/swagger.json).

Each image must have both an `owner` and a `name`.
An `owner` is typically the name of the team you are own. 
Images are referenced by `d.ims.io/<owner>/<name>`. 

For example, if team `carbon` created a `redis` image, images could be pushed via:
```
docker push d.ims.io/carbon/redis
```

## Authentication
All users must authenticate through their active directory or token credentials when interacting with `d.ims.io`.

### Active Directory
When using active directory credentials, use ONLY the username portion of the account. 

For example:

```
AD User: INTERNAL\john.doe
Email: john.doe@email.com
```

Should login with `john.doe` as a username.

To configure your Docker client to use active directory credentials, use the `docker login` command. 

For example:
```
docker login d.ims.io
Username: john.doe
Password: 
Login Succeeded
```

### Token
Tokens can be generated and used as a different form of authentication. 
Tokens are only valid for `d.ims.io`, and they do not expire. 
You can generate a token via the `/token` endpoint using the [Swagger UI](https://d.ims.io/api/?url=/swagger.json).

To configure your Docker client to use a token, create or update the `auth` section for `d.ims.io` in your Docker config file.
The Docker config file is located at `~/.docker/config.json`.

The file should follow the format:
```
{
        "auths": {
                "d.ims.io": {
                        "auth": "<token>"
                }
        }
}
```

## API  
The `d.ims.io` can be used to manage repositories and tokens. 
To explore and use the `d.ims.io` API, please navigate to the [Swagger UI](https://d.ims.io/api/?url=/swagger.json).
