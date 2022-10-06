# Example Go server with concurrency

## How to run it?

`make start`

OR

`make restart`

By default, server will run locally on `http://127.0.0.1:4001/`. Only one endpoint is defined: `/v1/api/smart`.

## Pre-commit

- Pre-commit needs to be installed on your machine (https://pre-commit.com/)
- run `pre-commit autoupdate` and `pre-commit install`
- pre-commit hooks are automatically installed

### Pre-commit config dependencies

- golangci-lint https://golangci-lint.run/usage/install/#local-installation
