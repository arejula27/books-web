[![Built with Devbox](https://www.jetify.com/img/devbox/shield_galaxy.svg)](https://www.jetify.com/devbox/docs/contributor-quickstart/) [![Go-test](https://github.com/arejula27/books-web/actions/workflows/go-test.yml/badge.svg)](https://github.com/arejula27/books-web/actions/workflows/go-test.yml)
# Project books

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

Use [devbox](https://www.jetify.com/devbox) for dependency installation. 

## Commands
Use the Makefile inside the `devbox shell`, if not use `devbox run <command>` instead `make <command>`

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
## Database
Starting the database use:
```bash
devbox services up --background 
```
Stoping the database use:
```bash
devbox services down
```
