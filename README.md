# depender

A tool for precompiling Go dependencies during Docker builds

Downloading and compiling the dependencies for your Go code can take a while. And if you're building container images from a Dockerfile, the problem is made worse by Docker's caching system: you need to copy your app code into the container in order to fetch and build its dependencies, but the cache gets invalidated every time you make the slightest change to your code. If an app depends on larger Go packages like [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes) or [aws-sdk-go](https://github.com/aws/aws-sdk-go), the build might take several minutes. Every time.

depender solves this problem by generating a Go source file that contains only imports. This file can be compiled after copying `go.{mod,sum}` but _before_ copying app code so that the slow dependency compilation step remains cached.

## Quick start

```
$ go get github.com/acj/depender
$ depender -h
Usage of depender:
  -exclude string
        package paths (or sub-strings) to exclude
  -output string
        path to output source file (default "deps.go")

```

## Usage

First, use depender to generate a source file that includes the dependencies for your app's packages.

```
$ cd my-go-app
$ depender ./...
Dummy source file written to deps.go
```

Then run `go build` on the dummy source file near top of your Dockerfile to compile the dependencies:

```
FROM golang:...

# Precompile dependencies
COPY go.mod go.sum ./
COPY deps.go ./
RUN go build deps.go && rm deps.go

# Now, copy your app source
COPY . .

...
```

Because the dependencies are compiled before your app's source code is copied, they will remain cached until you update go.mod, go.sum, or deps.go.

The `deps.go` file should be checked into version control. It contains a `+build ignore` directive so that it won't be picked up by the normal build tooling.

### Usage with `go generate`

If your dependencies change often, depender can be paired with `go generate` to streamline the workflow. For example, you can create a file containing the following:

```go
//go:generate go run github.com/acj/depender ./...
```

Then, running `go generate` will update `deps.go` with the current list of dependencies.

## Does it help?

For a recent project that imported kubernetes/kubernetes, aws-sdk-go, and several other nontrivial packages, this reduced my incremental container build times from ~110 seconds to roughly 15 seconds.
