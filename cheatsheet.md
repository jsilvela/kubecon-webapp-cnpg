# cheatsheet

## requesting and load generating

curl localhost:8080/
curl -H "Accept: application/json"  localhost:8080/

hey -H "Accept: application/json" -z 5s  http://localhost:8080/

hey -z 200s -q 1 -c 2  http://localhost:8080/update

## kubernetes

k config set-context --current --namespace=foo

kubectl rollout restart deployment/demo

pq: remaining connection slots are reserved for non-replication superuser connections
