SHELL_PATH = /bin/bash
SHELL = $(if $(wildcard $(SHELL_PATH)),/bin/bash,/bin/bash)


run:
	go run apis/services/sales/main.go | go run apis/tooling/logfmt/main.go

help: 
	go run apis/services/sales/main.go --help
version: 
	go run apis/services/sales/main.go --version

http:
	http GET http://localhost:3000/test
http-ready:
	http GET http://localhost:3000/readiness
http-live:
	http GET http://localhost:3000/liveness

http-gql:
	http POST http://localhost:8080/query \
		Content-Type:application/json \
		query='mutation createTodo { createTodo(input: { text: "todo", userId: "1" }) { user { id } text done } }'
# ==============================================================================
# Define dependencies

GOLANG          := golang:1.23
ALPINE          := alpine:3.20
KIND            := kindest/node:v1.31.2
POSTGRES        := postgres:17.2
GRAFANA         := grafana/grafana:11.3.0
PROMETHEUS      := prom/prometheus:v2.55.0
TEMPO           := grafana/tempo:2.6.0
LOKI            := grafana/loki:3.2.0
PROMTAIL        := grafana/promtail:3.2.0

KIND_CLUSTER    := enterprise-starter-cluster
NAMESPACE       := sales-system
GRAPHQLSPACE    := graphql-system
SALES_APP       := sales
GRAPHQL_APP     := graphql
AUTH_APP        := auth
BASE_IMAGE_NAME := localhost/jelam
VERSION         := 0.0.1
SALES_IMAGE     := $(BASE_IMAGE_NAME)/$(SALES_APP):$(VERSION)
GRAPHQL_IMAGE   := $(BASE_IMAGE_NAME)/$(GRAPHQL_APP):$(VERSION)
METRICS_IMAGE   := $(BASE_IMAGE_NAME)/metrics:$(VERSION)
AUTH_IMAGE      := $(BASE_IMAGE_NAME)/$(AUTH_APP):$(VERSION)

# ==============================================================================
# Building containers

build: sales graphql

sales:
	docker build \
		-f zarf/docker/sales.dockerfile \
		-t $(SALES_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.


# ==============================================================================
# Building containers for graphql service
build-graphql: graphql

graphql:
	docker build \
		-f zarf/docker/graphql.dockerfile \
		-t $(GRAPHQL_IMAGE) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
		.

# ==============================================================================
# Running from within k8s/kind
dev-up:
	kind create cluster \
		--image $(KIND) \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/dev/kind-config.yaml
	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner


dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

dev-status:
	watch -n 2 kubectl get pods -o wide --all-namespaces

# -------------------------------------------------------------------------
dev-load:
	kind load docker-image $(SALES_IMAGE) --name $(KIND_CLUSTER) & \
	kind load docker-image $(GRAPHQL_IMAGE) --name $(KIND_CLUSTER)  

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(SALES_APP) --timeout=120s --for=condition=Ready

	kustomize build zarf/k8s/dev/graphql | kubectl apply -f -
	kubectl wait pods --namespace=$(GRAPHQLSPACE) --selector app=$(GRAPHQL_APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(SALES_APP) --namespace=$(NAMESPACE)
	kubectl rollout restart deployment $(GRAPHQL_APP) --namespace=$(GRAPHQLSPACE)


dev-update: build dev-load dev-restart

dev-update-apply: build dev-load dev-apply

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(SALES_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go -service=$(SALES_APP)

# graphql-logs:
# 	kubectl logs --namespace=$(GRAPHQLSPACE) -l app=$(GRAPHQL_APP) --all-containers=true -f --tail=100 --max-log-requests=6 | go run apis/tooling/logfmt/main.go -service=$(GRAPHQL_APP)

# -------------------------------------------------------------------------

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(SALES_APP)

dev-describe-graphql-deployment:
	kubectl describe deployment --namespace=$(GRAPHQLSPACE) $(GRAPHQL_APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(SALES_APP)

dev-describe-graphql-sales:
	kubectl describe pod --namespace=$(GRAPHQLSPACE) -l app=$(GRAPHQL_APP)

# ==============================================================================
# Metrics and Tracing

metrics:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open http://localhost:3010/debug/statsviz
# =========================================================================
# Modules support

tidy: 
	go mod tidy
	go mod vendor
