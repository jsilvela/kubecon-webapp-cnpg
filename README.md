# Basic webapp in Go + CNPG

https://apoorvtyagi.tech/containerize-your-web-application-and-deploy-it-on-kubernetes
https://stackoverflow.com/questions/30746888/how-to-know-a-pods-own-ip-address-from-inside-a-container-in-the-pod

## developing against a local or dockerized Postgres

A fairly typical scenario when developing a webapp or other kind of application with
a RDBMS for storage would be:

- Either install PostgreSQL locally, or get a dockerized version and do port forwarding
- Write the application on bare metal or dockerize
- Use schema migration tools to capture changes in dev that should be applied in
  prod
- Following 12-factor app recommendations, pass the DB credentials as environment
  variables to the application

``` sh
PG_PASSWORD=<redacted> PG_USER=app go run main.go
```

There's nothing to stop you from using CloudNativePG as your dockerized Postgres.

In my Kind cluster, I have the simplest possible CloudNativePG cluster, with 1 Pod
running PostgreSQL 14. By default, a database called `app` is created, owned by a
user `app`.

Let's write a basic (very basic) webapp in Go.
Using the [lib/pq](https://pkg.go.dev/github.com/lib/pq) library, we're going to
use the following connection string:
> "postgres://<user>:<password>@localhost/app?sslmode=require"

``` sh
kubectl port-forward service/cluster-example-rw  5432:5432
PG_PASSWORD=4semkumW28zy8OokujlmDxgdwX2ilh5TC9RaeeOatuaxhi8wC1vjTUcVfoEzFc3P PG_USER=app go run main.go

/usr/local/opt/liquibase/liquibase rollback-count 1 --changelog-file=example-changelog.sql
/usr/local/opt/liquibase/liquibase update --changelog-file=example-changelog.sql
```

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
