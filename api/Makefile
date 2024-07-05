.PHONY:  run 

# Set default env as develop.
ifeq ($(ENV),)
ENV := develop
endif

-include .env.$(ENV).secret

WORK_DIR := .

install:
	go mod download
	go install github.com/cosmtrek/air@v1.52.1

run:
	PORT=$(PORT) air -c .air.toml
