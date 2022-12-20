# Tat's micro-blog playground
This repository serves as a micro-blog and testbed meant for trying out new software and programming paradigms.  

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
```

From the base directory, run `docker compose up` (add `-d` if running in detached mode).  

## Structure
- `/pkg`, `/internal`, `/assets`, `/templates`, `main.go`  
Main webserver logic and other modules here
- `/terraform`  
IaC related code here
- `/docker`  
Dockerfiles for each container here
- `.github`  
Github workflows here

## Current to-dos
- [x] Setting up of infrastructure  
- [ ] Enable deployment to AWS via workflow
- [ ] Blog landing page and general skeletal framework
- [ ] Secured login endpoint for non-public items