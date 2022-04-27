package kroller

import (
	"dagger.io/dagger"
	"universe.dagger.io/go"
	//"universe.dagger.io/alpine"
)

dagger.#Plan

client: filesystem: "./": read: contents: dagger.#FS
client: env: CGO_ENABLED: string | *"0"

// build kroller
actions: build: go.#Build & {
	source: client.filesystem."./".read.contents
	env: CGO_ENABLED: client.env.CGO_ENABLED
}

// test kroller
actions: test: go.#Test & {
	source:  client.filesystem."./".read.contents
	package: "./..."
	env: CGO_ENABLED: client.env.CGO_ENABLED
}

actions: all: {
	build: actions.build
	test:  actions.test
}
