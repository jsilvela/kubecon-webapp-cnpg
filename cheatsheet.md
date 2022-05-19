# cheatsheet

## requesting and load generating

curl localhost:8080/update

hey -z 100s -q 1 -c 2  http://localhost:8080/update

## kubernetes

kubectl rollout restart deployment mywebapp

## getting the DB credentials to hit ouside K8s

1. kubectl get secret cluster-example-app -o yaml
1. copy the `password` in the  YAML
1. echo copiedSecret | base64 --decode | pbcopy

## app

PG_PASSWORD=***** PG_USER=app go run main.go
