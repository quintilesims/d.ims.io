# ECR Repo Migrator

# Requirements
* AWS CLI
* AWS IAM Credentials for both AWS ECR Accounts
* Docker
* Python 3

# Setup
First, use the aws cli to get docker credentials for both ECR repositories:
```
aws ecr get-login --no-include-email
```

This script requires the following environment variables to run:
```
SRC_ACCESS_KEY - The aws_access_key_id for the source ECR repository
SRC_SECRET_KEY - The aws_secret_access_key for the source ECR repository
SRC_REGION     - The aws_region of the source ECR repository

DST_ACCESS_KEY - The aws_access_key_id for the destination ECR repository
DST_SECRET_KEY - The aws_secret_access_key for the destination ECR repository
DST_REGION     - The aws_region of the destination ECR repository
```

Optionally, you can update the `MAX_WORKERS` variable in `migrate.py` to make this process threaded. 

# Run
```
python migrate.py
```

Depending on how many repositories and the setting for `MAX_WORKERS`, this process can take hours and even days to complete.
The script checks the diff between the two repositories before it starts migrating images, so it is ok to run this script 
as many times as necessary.
