# seller-analytics
This repository is an example on how we can use message queues, REST APIs and a database to record seller statistics and analytics in an e-commerce platform.


## Overview
We will be making an e-commerce platform built on top a microservice architecture.

This repository is a monorepo housing 3 different services
1. buyer
2. statistic
3. analytics


## Technologies Used
### Build and Deployment
- bazel: build tool (https://bazel.build/)
- docker: containerization


## Golang Packages
- gin: HTTP web framework (https://github.com/gin-gonic/gin)
- viper: Golang configuration solution (https://github.com/spf13/viper)
- fx: Dependency injection system (https://github.com/uber-go/fx)
- gorm: Golang SQL object relational mapper (https://gorm.io/)


## Getting Started
### Prerequisites
1. It is good to get gvm (https://github.com/moovweb/gvm) to make it simple to switch between go versions just in case other projects require different versions. Follow instructions here for ARM64 machines (https://github.com/moovweb/gvm/issues/385)
3. install bazel, it is recommended to use bazelisk
   ```
   brew install bazelisk
   ```
3. install gazelle
   ```
   go install github.com/bazelbuild/bazel-gazelle/cmd/gazelle@latest
   ```
4. install docker, recommend to use docker desktop (https://docs.docker.com/engine/install/). If prefer not to use docker desktop then can install docker and docker-compose separately.

### Running a service
1. Run `gazelle`, gazelle is a build file generator for bazel. It automatically generates build files. For our purposes, run this command everytime we add a new import statement on a golang file.
2. Run `bazel run //src/services/{service name}` which will start the service
3. Debugging is a bit more tedious but possible using dlv command. We first build the binary with debug flag and attach dlv to it. The binary can be found under bazel-bin directory.
```
    bazel build -c dbg //src/services/buyer
    dlv exec --api-version=2 bazel-bin/src/services/buyer/buyer_/buyer.runfiles/__main__/src/services/buyer/buyer_/buyer
```
you also need to adjust the `DefaultConfigPath` in this case. This is less than ideal but is a working workaround.

It is also possible to debug via attaching a debugger to the process, if anyone is interested please try and provide feedback so we may add it here.

## Project Structure
- src: contains all source code
    - pkg: stores library code that is shareable across services
    - services: stores the applications corresponding to our services in our microservice ecosystem. The structure follows Golang clean architecture as mentioned in your handbook
      - config: contains configuration files and source code
      - domain: contains domain models specific to the service (domain/model/entity)
      - handler: contains source code to map domain objects to response format (controller/delivery)
      - repository: contains source code that interacts with other data sources and integrations (repository)
      - usecase: contains source code that implements business logic (usecase/service)


## How tos
### Import package in go file
After adding a new import in a .go file, always run
```
gazelle
```
gazelle will update the BUILD file for the directory accordingly

### Add new golang dependency to go.mod
As usual use `go get` or import and `go mod tidy`
After the changes have been reflected in go.mod run the following
```
gazelle update-repos --from_file=go.mod -to_macro=deps.bzl%go_deps
```



