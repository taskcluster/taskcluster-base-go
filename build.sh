#!/bin/bash

cd "$(dirname "${0}")"

go clean -i -x ./...

go get ./...
go install -v -x ./...
go fmt ./...
go get golang.org/x/tools/cmd/vet
go vet -x ./...
go test -v ./...

# finally check that generated files have been committed, and that
# formatting code resulted in no changes...
git status
[ $(git status --porcelain | wc -l) == 0 ]
