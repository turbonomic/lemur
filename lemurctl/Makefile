GITCOMMIT=$(shell git rev-parse HEAD)
BUILDTIME=$(shell date -R)

REMOTE=github.com
USER=turbonomic
PROJECT=lemur/lemurctl
PROJECTPATH=${REMOTE}/${USER}/${PROJECT}
BINARY=lemurctl
SOURCE_DIRS=cli pkg utils

build: clean
	go build -ldflags "-X '${PROJECTPATH}/version.GitCommit=${GITCOMMIT}' -X '${PROJECTPATH}/version.BuildTime=${BUILDTIME}'" -o ${BINARY}

.PHONY: fmtcheck
fmtcheck:
	@gofmt -l $(SOURCE_DIRS) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

.PHONY: clean
clean:
	@rm -rf ${BINARY}
