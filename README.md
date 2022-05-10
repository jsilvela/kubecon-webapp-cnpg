# Basic webapp in Go + CNPG

https://apoorvtyagi.tech/containerize-your-web-application-and-deploy-it-on-kubernetes
https://stackoverflow.com/questions/30746888/how-to-know-a-pods-own-ip-address-from-inside-a-container-in-the-pod


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
kubectl port-forward deployment/mywebapp  8080:5000 -n demo
kubectl port-forward service/mywebapp  8080:8088 -n demo
```
