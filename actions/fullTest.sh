#!/bin/bash
$(
    cd src
    go test -v -coverpkg=./pkg/... -coverprofile=profile.cov ./pkg/...
    go tool cover -func profile.cov
    rm profile.cov
)