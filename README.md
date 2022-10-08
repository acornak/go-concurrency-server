# Example Go server with concurrency

## What does it do?

After calling endpoint `/v1/api/smart` with `timeout` parameter, server will make GET request to URL defined in environment variables (`.env`). **(TBD)** If server will not respond within 300 milliseconds, another 2 requests are fired concurrently. Endpoint will return the first successful response (including the first one). If there is no successful response, endpoint will return error. If there is no response from any of requests, endpoint will return error.

## How to run it?

`make start`

OR

`make restart`

Docker must be running!

By default, server will run locally on `http://127.0.0.1:4001/`. Only one endpoint is defined: `/v1/api/smart` with query param `timeout`. Timeout can only be numeric value, otherwise `BAD REQUEST` is returned.

## Pre-commit

- Pre-commit needs to be installed on your machine (https://pre-commit.com/)
- run `pre-commit autoupdate` and `pre-commit install`
- pre-commit hooks are automatically installed from `.pre-commit-config.yaml` file

### Pre-commit config dependencies

You need to install following dependencies manually:

- golangci-lint https://golangci-lint.run/usage/install/#local-installation

## Performance testing

Run script `performance_test.sh`.
