# Basic webapp in Go + CNPG - a walk-through

https://apoorvtyagi.tech/containerize-your-web-application-and-deploy-it-on-kubernetes
https://stackoverflow.com/questions/30746888/how-to-know-a-pods-own-ip-address-from-inside-a-container-in-the-pod

## Introduction

Developing a webapp using an RDBMS involves quite a bit of setup on developer
machines. The Database itself is usually installed locally, or using Docker for those
teams that already have a Docker investment.

The deployments are typically of the most basic type, with one Database running.
Replication, load tests, RTO measurements etc are left for production. Very often, the
developers have only basic knowledge of SQL, and no experience with DB
administration and operations.

In this context, there is a large gap between the development environment and the
production and staging environments. The reduction of this gap was precisely one of
the tenets of the DevOps movment.

CloudNativePG is a great tool to eliminate the gap between development and
production.
Developers can use CloudNativePG to spin up a PostgreSQL database on their
machine, add read
replicas, test failover, redirect Read-only traffic to read replicas, test out disaster
scenarios, scale up and down, etc.

## Why develop with PosgreSQL?

Simply, PostgreSQL is the most advanced Open Source RDBMS  out there. It is also
a pleasure to develop code against.

Some great features that will allow you to write streamlined SQL:

- Common Table Expressions (i.e. `WITH` statements)
- Windowing functions
- Grouping Sets
- Array types
- Stored procedures, triggers
- **Transactional DDL** (more detail next section)
- Powerful GIS with PostGIS
- Transactions (obviously!!)

## Basic development cycle for a webapp with a Database

A fairly typical scenario when developing a webapp or other kind of application with
a RDBMS for storage would be:

- Either install PostgreSQL locally, or get a dockerized version and do port forwarding
- Write the application on bare metal or dockerize it
- Avoid using the `postgres` superuser in your webapp. Create a custom user with only
  the necessary permissions
- Use schema migration tools to capture changes in dev, and commit them to your
  code repo
- Following 12-factor app recommendations, pass the DB credentials as environment
  variables to the application

Schema migration tools are a necessity to have a sane collaboration environment.
Liquibase or Flybase are well known examples, and there are other tools available
using different implementation languages. We use liquibase in this post.

Schema migration tools work best when their "Rollback" command works cleanly. PostgreSQL
shines here. With its **Transactional DDL**, no cleanup is needed after a Rollback.
Let's just say there are other databases that can't make that claim.

There's nothing to stop you from using CloudNativePG as you would a local or dockerized
Postgres instance, eventhough it can do much more.
Let's start with that: \
In my Kind cluster, I have the simplest possible CloudNativePG cluster, with 1 Pod
running PostgreSQL 14. By default, a CloudNativePG creates a database called `app`,
owned by a user called `app`.

Let's write a basic (very basic) webapp in Go.
Using the [lib/pq](https://pkg.go.dev/github.com/lib/pq) library, we're going to
use the following connection string:


k config set-context --current --namespace=foo

k rollout deployment

punnett square


``` none
"postgres://<user>:<password>@localhost/app?sslmode=require"
```

Our app is so simple it's in just one file: `main.go`
We can run it like so:

``` sh
PG_PASSWORD=<redacted> PG_USER=app go run main.go
```

The laptop I'm using is a mac, so I need to do port forwarding in order for my Kind
cluster to open a port in my local network:

``` sh
kubectl port-forward service/cluster-example-rw  5432:5432
```

That's all there is to it. I can now open `http://localhost:8080/` in my browser.

``` sh
kubectl port-forward service/cluster-example-rw  5432:5432
PG_PASSWORD=xE5CX7zpCr5ccmlBJALvifQRW36mSmN48hEEsuUs0fEhyHEE7GDeV1aCsMEZaeXj PG_USER=app go run main.go

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

Let's run a load generator:
hey -H "Accept: application/json" -z 5s  http://localhost:8080/