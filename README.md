lawn
===

GitHub Contributions Graph scraper.

## Features

- Provide svg image of your GitHub Contributions Graph.
- Support Google Cloud Run and other container services.

## Usage as library

```
go get github.com/ww24/lawn
```

## Configure

Some environment variables are required.

- Set GCP `PROJECT_ID`.
- Set `GITHUB_USERNAME`.

## Build

Build docker image.

```bash
make build
```

## Deploy

Push docker image to GCR.

```bash
make push
```

### Terraform

Deploy to Cloud Run.

1. Create backend backet on GCS and set `BACKEND_BUCKET` environment variable.
1. `make init`
1. `make plan`
1. `make apply`

## Use by Firebase Hosting

- https://firebase.google.com/docs/hosting/manage-cache
- https://firebase.google.com/docs/hosting/cloud-run
