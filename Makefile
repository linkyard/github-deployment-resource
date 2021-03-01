VERSION=v0.10.0

build:
	docker build -t linkyard/github-deployment-resource:$(VERSION) .