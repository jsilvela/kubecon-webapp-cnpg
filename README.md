# Basic webapp in Go + CNPG

https://apoorvtyagi.tech/containerize-your-web-application-and-deploy-it-on-kubernetes

## shell

``` sh
go run main.go
```

## docker

``` sh
docker run -dp 5000:5000 myapp
docker build -t myapp .
```

## kind

``` sh
k apply -f deployment.yaml 
docker images
kind load docker-image myapp:latest --name pg-operator-e2e-v1-23-1
```
