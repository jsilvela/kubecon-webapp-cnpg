# Webapp development with CloudNativePG

This repo has support code for a Virtual Office Hours demo given at the
2022 Kubecon Valencia, on 2022-05-18 by me (Jaime).

Contains:

- Liquibase config and initial migration
- Webapp in Go
- Makefile and cheat-sheet (to aid this poor typist)
- Dockerfile and K8s YAML's for the webapp

Assumed to be installed on demo machine:

- KinD cluster
- CloudNativePG operator
- Liquibase

Uses the `cluster-example.yaml` sample in `docs/src/samples`.

## summary

The purpose of the demo is to show CloudNativePG is a great tool to develop
web applications on dev machines and opens up possibilities to do all sorts of
things before putting in prod.

We try to show a viable workflow using Liquibase for schema migrations, 
and a webapp written in Go using only the standard
`net/http`, `database/sql` and `lib/pq` packages.

Two halves to the demo:

1. "Average" or "old-style" approach: postgres installed locally, possibly even
  dockerized + webapp developed locally or perhaps dockerized.
  We ape this by installing a 1-pod cluster using the `cluster-example.yaml` sample
  with 1 replica instead of 3.
  We run the webapp outside of KinD and hit `localhost:5432` from it.

2. Better way: we dockerize the webapp, push to the KinD nodes, create a
  deployment and a service for it. We hit the `cluster-example-rw` service.
  We show the possibilities that open up.

## Demo and slides

### Slide: Using CloudNativePG from outside Kubernetes

Show diagram for "Case 2" from the
[cnpg docs](https://cloudnative-pg.io/documentation/1.15.0/use_cases/)

### Slide: "The Old Way", but using CloudNativePG

#### Ingredient list, game plan

1. Basic web application written in Go (aka golang)
1. Kubernetes cluster (in this case KinD running on macOS)
1. The CloudNativePG operator installed on KinD
1. Schema migration tool: Liquibase

1. Create the simplest CloudNativePG cluster
1. Start with an empty PostgreSQL DB
1. Get DB connection credentials
1. Apply migrations
1. Add port forwarding to expose DB (necessary if running KinD on macOS)
1. Start the webapp

### Slide: Putting your webapp inside Kubernetes

Show first diagram for "Case 1" from the
[cnpg docs](https://cloudnative-pg.io/documentation/1.15.0/use_cases/)

Take opportunity to show the web page for cnpg and the documentation link.

### Slide: Leveraging CloudNativePF - cooking with fire

#### Ingredient list, game plan

1. Let’s Dockerize our webapp and load it into our KinD cluster
1. Make a deployment and a loadBalancer for our webapp
1. The webapp should now hit the -rw (cluster-example-rw) Service

1. Let’s scale our CloudNativePG cluster to 3 instances
1. Close the running webapp
1. Deploy the webapp  (show Dockerfile and K8s YAML)
1. Let’s add a watch on the cluster. Let’s put some load on it
1. Let’s kill the primary!!
1. (with enough time) - Have a look around in the cloudnative-pg repo

## Slide: Reflections and questions

Talk about how once we start leveraging CloudNativePG, we enable developers
to do experiments that require DBA skills. And we allow them to develop locally
with a DB much more similar to what production will look like.
When they deal with prod issues on their DB, those developers will have gained
hand-on experience already.

And bridging the gap betwenn dev and prod is one of the tenets of the DevOps
movement.
