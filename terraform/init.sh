#!/bin/bash

sudo amazon-linux-extras install nginx1
sudo yum -y install docker git

# configuring nginx reverse proxy for port 8080
cat <<EOT >> /etc/nginx/conf.d/tatwritescode.com.conf
server {
    listen       80;
    listen       [::]:80;
    server_name  _;

    return 301 https://$host$request_uri;
}

server {
    listen              443 ssl;
    ssl_certificate     /etc/ssl/tatwritescode.com.pem;
    ssl_certificate_key /etc/ssl/tatwritescode.com.key;

    server_name  tatwritescode.com www.tatwritescode.com;

    location / {
        proxy_pass http://localhost:8080;
    }
}
EOT

# creating public and private keys on respective locations
cat <<EOT >> /etc/ssl/tatwritescode.com.pem
${SSL_PEM}
EOT

cat <<EOT >> /etc/ssl/tatwritescode.com.key
${SSL_KEY}
EOT

# post docker installation steps
usermod -aG docker ec2-user
newgrp docker

# manually installing docker compose v2 as it's not supported in AL2 yet
# https://github.com/aws/aws-codebuild-docker-images/issues/527
DOCKER_CONFIG=/usr/local/lib/docker
mkdir -p $DOCKER_CONFIG/cli-plugins
curl -SL https://github.com/docker/compose/releases/download/v2.14.0/docker-compose-linux-x86_64 -o $DOCKER_CONFIG/cli-plugins/docker-compose

chmod +x /usr/local/lib/docker/cli-plugins/docker-compose

sudo service docker start

# spinning up twc-app container and mysql container
sudo git clone https://github.com/ztcjoe93/tatwritescode.git /opt/tatwritescode

# pass terraform variables to populate within .env file for docker-compose
sudo cat <<EOT >> /opt/tatwritescode/.env
MYSQL_DATABASE=${MYSQL_DATABASE}
MYSQL_USER=${MYSQL_USER}
MYSQL_PASSWORD=${MYSQL_PASSWORD}
MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
MYSQL_HOST=${MYSQL_HOST}
VOLUME_MOUNT_PATH=${VOLUME_MOUNT_PATH}
UPLOAD_MOUNT_PATH=${UPLOAD_MOUNT_PATH}
EOT

# mkfs does not reformat the volume unless -f is passed in
sudo mkfs -t xfs /dev/xvdh
sudo mkdir -p /opt/tatwritescode/volume
sudo mount /dev/xvdh /opt/tatwritescode/volume

sudo mkdir /opt/tatwritescode/volume/database

cd /opt/tatwritescode
docker compose up -d

sudo systemctl start nginx.service
