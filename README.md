![](assets/twc.png)
# Tat's micro-blog playground
This repository serves as a micro-blog and testbed meant for trying out new software and programming paradigms.  

## Table of Contents
- [Setup](#setup)
- [Infrastructure](#infrastructure)
- [Running tests](#running-tests)
- [Structure](#structure)
- [Current to-dos](#current-to-dos)

## Setup

### The docker way
Create a .env file to declare the variables required for MySQL docker container (can omit `VOLUME_MOUNT_PATH` if not using a specific volume for database)
```
MYSQL_DATABASE=<DATABASE_NAME>
MYSQL_USER=<DATABASE_USER>
MYSQL_PASSWORD=<DATABASE_PASSWORD>
MYSQL_ROOT_PASSWORD=<DATABASE_ROOT_PASSWORD>
MYSQL_HOST=<HOST_NAME>
VOLUME_MOUNT_PATH=<VOLUME_MOUNT_PATH_FOR_DOCKER>
SIGNATURE_KEY=<SIGNING>
ENV=<ENVIRONMENT>
```

From the base directory, run `docker compose up` (add `-d` if running in detached mode).  

## Infrastructure
AWS resources are all IaC controlled with Terraform; only difference is where some resources are ephemeral,  
while others is spun-up manually using Terraform (i.e EBS, EIP).  

Set up required terraform variables to be passed over by creating a `.tfvars` file with the following:  
```
database_user = <DATABASE_USER>
database_host = <DATABASE_HOST>
database_password = <DATABASE_PASSWORD>
database_root_password = <DATABASE_ROOT_PASSWORD>
ssl_pem = <<-EOT
-----BEGIN CERTIFICATE-----
<PUBLIC_KEY_HERE>
-----END CERTIFICATE-----
EOT
ssl_key = <<-EOT
-----BEGIN PRIVATE KEY-----
<PRIVATE_KEY_HERE>
-----END PRIVATE KEY-----
EOT
volume_mount_path = <VOLUME_MOUNT_PATH>
```

To check what resources are created/destroyed:  
```
$ terraform init
$ terraform plan
```

To spin up actual resources in AWS (ensure that you have eitehr AWS credentials stored, or respective secret access key + access key id in env)  
```
$ terraform apply -y
```

Post-deployment initiliazation script can be found in `/terraform/init.sh`, and output can be found in `/var/log/cloud-init-output.log` in the EC2 instance.

## Running tests
To run all tests in the `/test` sub-directory  
```
$ go test -v ./test/...
```

## Structure
- `/pkg`, `/internal`, `/assets`, `/templates`, `main.go`  
Main webserver logic and other modules here
- `/test`  
Tests for internal/main module here
- `/terraform`  
IaC related code here
- `/docker`  
Dockerfiles for each container here
- `.github`  
Github workflows here

## Current to-dos
- [x] Setting up of infrastructure  
- [x] Blog landing page and general skeletal framework
- [x] Secured login endpoint for non-public items
- [ ] Add test harness for all existing internal packages
- [ ] Enable deployment to AWS via workflow