VERSION=v0.3.1

build:
	docker build -t linkyard/github-deployment-resource:$(VERSION) .