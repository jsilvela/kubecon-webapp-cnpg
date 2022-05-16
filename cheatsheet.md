# cheatsheet

## requesting and load generating

curl localhost:8080/
curl -H "Accept: application/json"  localhost:8080/

hey -H "Accept: application/json" -z 5s  http://localhost:8080/
