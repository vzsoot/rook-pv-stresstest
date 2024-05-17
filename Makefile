VERSION ?= $(shell cat VERSION)
REGISTRY_LOCAL ?=
REGISTRY ?=

all: build-helm build-image-consumer build-image-producer

build-helm:
	helm package deploy/helm/rook-pv-stresstest --dependency-update -d build/ --version $(VERSION) --app-version $(VERSION)

build-image-consumer:
	docker build -t $(REGISTRY_LOCAL)rook-pv-stresstest-consumer:$(VERSION) -f src/golang/Dockerfile src/golang/.

build-image-producer:
	docker build -t $(REGISTRY_LOCAL)rook-pv-stresstest-producer:$(VERSION) -f src/c/Dockerfile src/c/.

publish: all
	docker tag $(REGISTRY_LOCAL)rook-pv-stresstest-consumer:$(VERSION) $(REGISTRY)/rook-pv-stresstest-consumer:$(VERSION); \
	docker tag $(REGISTRY_LOCAL)rook-pv-stresstest-producer:$(VERSION) $(REGISTRY)/rook-pv-stresstest-producer:$(VERSION); \
	docker push $(REGISTRY)/rook-pv-stresstest-consumer:$(VERSION);
	docker push $(REGISTRY)/rook-pv-stresstest-producer:$(VERSION);
