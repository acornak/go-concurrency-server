# Example Go server with concurrency

## How to run it?

`make start`

OR

`make restart`

Docker must be running!

By default, server will run locally on `http://127.0.0.1:4001/`. Only one endpoint is defined: `/v1/api/smart` with query param `timeout` that accepts only numberic values.

## Pre-commit

- Pre-commit needs to be installed on your machine (https://pre-commit.com/)
- run `pre-commit autoupdate` and `pre-commit install`
- pre-commit hooks are automatically installed from `.pre-commit-config.yaml` file

### Pre-commit config dependencies

You need to install following dependencies manually:

- golangci-lint https://golangci-lint.run/usage/install/#local-installation

## Performance testing

Run script `performance_test.sh`.
