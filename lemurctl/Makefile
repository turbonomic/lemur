GITCOMMIT=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date -R)

REMOTE=github.com
USER=turbonomic
PROJECT=lemur/lemurctl
PROJECTPATH=${REMOTE}/${USER}/${PROJECT}
BINARY=lemurctl

build: clean
	go build -ldflags "-X '${PROJECTPATH}/version.GitCommit=${GITCOMMIT}' -X '${PROJECTPATH}/version.BuildTime=${BUILDTIME}'" -o ${BINARY} 

.PHONY: clean
clean:
	@rm -rf ${BINARY}
