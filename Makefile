.PHONY: all
ifeq ($(ARCH), darwin)
all: darwin
else
all: linux
endif

.PHONY: darwin
darwin: dep
	env GOOS=darwin GOARCH=amd64 go build -o aws-vpc-peering ./...

.PHONY: linux
linux: dep
	env GOOS=linux GOARCH=amd64 go build -o aws-vpc-peering ./...

.PHONY: dep
dep:
ifndef ARCH
	$(error ARCH is not set)
endif
