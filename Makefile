dockerbuild: *.go
	docker build -t myapp .

deploy: dockerbuild
	kind load docker-image myapp:latest --name pg-operator-e2e-v1-23-1

portfwd:
	kubectl port-forward service/mywebapp  8080:8088

pgportfwd:
	kubectl port-forward service/cluster-example-rw  5432:5432
