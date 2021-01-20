# Middle Manager

![CI](https://github.com/theartofeducation/middle-manager/workflows/CI/badge.svg?branch=main)

A service to handle API integrations for our tooling.

## Development

### Getting Started

1. Copy `.env.sample` to `.env`
1. Update `.env` values with API URLs and secrets
1. Run `go get` to download dependencies
1. Setup ngrok to tunnel to your local environment
1. Create webhooks in ClickUp using their API
1. Copy the webhook secret to the `.env` files
1. Run `go run .` to start the server
1. Remove the created webhook when finished

## Features

### ClickUp + Clubhouse

Middle Manager automates information updates between [ClickUp](https://clickup.com/) and [Clubhouse](https://clubhouse.io/).

* [ClickUP API](https://clickup20.docs.apiary.io/#)
* [Clubhouse API](https://clubhouse.io/api/rest/v3/)

#### Create Clubhouse Epic

When a ClickUp task is moved to the status Ready For Development an Epic will be created in Clubhouse.

#### Move ClickUp Task to Acceptance

When an Epic is marked as "Done" in Clubhouse the corresponding Task in ClickUp will be moved to the "Acceptance" column.
