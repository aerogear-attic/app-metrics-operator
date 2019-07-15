APP_NAME = app-metrics-operator
ORG_NAME = aerogear
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
TOP_SRC_DIRS = pkg
PACKAGES ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
              -exec dirname {} \\; | sort | uniq")
TEST_PKGS = $(addprefix $(PKG)/,$(PACKAGES))
APP_FILE=./cmd/manager/main.go

NAMESPACE=app-metrics
APP_NAMESPACES=app-metrics-apps
CODE_COMPILE_OUTPUT = build/_output/bin/app-metrics-operator
TEST_COMPILE_OUTPUT = build/_output/bin/app-metrics-operator-test

##############################
# Local Development          #
##############################

.PHONY: code/run
code/run: code/gen
	operator-sdk up local

.PHONY: code/gen
code/gen: code/fix
	operator-sdk generate k8s
	operator-sdk generate openapi
	go generate ./...

.PHONY: code/fix
code/fix:
	gofmt -w `find . -type f -name '*.go' -not -path "./vendor/*"`

##############################
# Jenkins                    #
##############################

.PHONY: test/compile
test/compile:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test -c -o=$(TEST_COMPILE_OUTPUT) ./test/e2e/...


.PHONY: code/compile
code/compile: code/gen
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o=$(CODE_COMPILE_OUTPUT) ./cmd/manager/main.go

##############################
# Tests / CI                 #
##############################

.PHONY: test/integration-cover
test/integration-cover:
	echo "mode: count" > coverage-all.out
	GOCACHE=off $(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)

.PHONY: test/unit
test/unit:
	@echo Running tests:
	CGO_ENABLED=1 go test -v -race -cover $(TEST_PKGS)

.PHONY: code/build/linux
code/build/linux:
	env GOOS=linux GOARCH=amd64 go build $(APP_FILE)

##############################
# Application                #
##############################

.PHONY: cluster/prepare
cluster/prepare:
	-kubectl create namespace $(NAMESPACE)
	-kubectl create namespace ${APP_NAMESPACES}
	-kubectl label namespace $(NAMESPACE) monitoring-key=middleware
	-kubectl apply -n $(NAMESPACE) -f deploy/service_account.yaml
	-kubectl apply -n $(NAMESPACE) -f deploy/role.yaml
	-kubectl apply -n $(NAMESPACE) -f deploy/role_binding.yaml
	-kubectl apply -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsservice_crd.yaml
	-kubectl apply -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsconfig_crd.yaml

.PHONY: cluster/clean
cluster/clean:
	make example-app/delete
	-kubectl delete -n $(NAMESPACE) appMetricsservice --all
	-kubectl delete -n $(NAMESPACE) -f deploy/role.yaml
	-kubectl delete -n $(NAMESPACE) -f deploy/role_binding.yaml
	-kubectl delete -n $(NAMESPACE) -f deploy/service_account.yaml
	-kubectl delete -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsservice_crd.yaml
	-kubectl delete -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsconfig_crd.yaml
	-kubectl delete namespace $(NAMESPACE)

.PHONY: install
install:
	-kubectl apply -n $(NAMESPACE) -f deploy/operator.yaml
	-kubectl apply -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsservice_cr.yaml

.PHONY: uninstall
uninstall:
	-kubectl delete -n $(NAMESPACE) -f deploy/crds/metrics_v1alpha1_appmetricsservice_cr.yaml
	-kubectl delete -n $(NAMESPACE) -f deploy/operator.yaml

.PHONY: example-app/apply
example-app/apply:
	-kubectl apply -n $(APP_NAMESPACES) -f deploy/crds/metrics_v1alpha1_appmetricsconfig_cr.yaml

.PHONY: example-app/delete
example-app/delete:
	-kubectl delete -n $(APP_NAMESPACES) -f deploy/crds/metrics_v1alpha1_appmetricsconfig_cr.yaml


.PHONY: monitoring/install
monitoring/install:
	@echo Installing service monitor in ${NAMESPACE} :
	- kubectl label namespace ${NAMESPACE} monitoring-key=middleware
	- kubectl apply -n $(NAMESPACE) -f deploy/monitor/service_monitor.yaml
	- kubectl apply -n $(NAMESPACE) -f deploy/monitor/operator_service.yaml
	- kubectl apply -n $(NAMESPACE) -f deploy/monitor/prometheus_rule.yaml
	- kubectl apply -n $(NAMESPACE) -f deploy/monitor/grafana_dashboard.yaml

.PHONY: monitoring/uninstall
monitoring/uninstall:
	@echo Uninstalling monitor service from ${NAMESPACE} :
	- kubectl delete -n $(NAMESPACE) -f deploy/monitor/service_monitor.yaml
	- kubectl delete -n $(NAMESPACE) -f deploy/monitor/prometheus_rule.yaml
	- kubectl delete -n $(NAMESPACE) -f deploy/monitor/grafana_dashboard.yaml
	- kubectl delete -n $(NAMESPACE) -f deploy/monitor/operator_service.yaml