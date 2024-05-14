VERSION ?= $(shell cat VERSION)
REGISTRY_LOCAL ?=
REGISTRY ?=

all: build-helm build-image

build-helm:
	helm package deploy/helm/rook-pv-stresstest --dependency-update -d build/ --version $(VERSION) --app-version $(VERSION)

build-image:
	docker build -t $(REGISTRY_LOCAL)rook-pv-stresstest:$(VERSION) -f src/Dockerfile src/.

publish: build-image build-helm
	docker tag $(REGISTRY_LOCAL)rook-pv-stresstest:$(VERSION) $(REGISTRY)/rook-pv-stresstest:$(VERSION); \
	docker push $(REGISTRY)/rook-pv-stresstest:$(VERSION);
