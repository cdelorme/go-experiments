#!/bin/bash

# ldflags can be used to send parameters to each package
# for example: go build -ldflags "-X main.abc=123"

# while you can set other package parameters that increases
# the necessary awareness and probably documentation needed

# this build script will supply parameters
# you can verify that they are kept by copying the executable
# anywhere and running it again

go build -ldflags "-X main.buildDate '$(date)' -X main.softwareName=Demo -X main.abc=123 -X main.buildVersion=0.1.0"

# @note: not sure how to wrap a computed date when they fully deprecate the old syntax
