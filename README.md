# Middle Manager

![CI](https://github.com/theartofeducation/middle-manager/workflows/CI/badge.svg?branch=main)

A service to handle API integrations for our tooling.

## Development

### Getting Started

1. Copy `.env.sample` to `.env`
1. Update `.env` values as needed
1. Run `go get` to download dependencies
1. Setup ngrok to tunnel to your local environment
1. Create webhooks in ClickUp
1. Run `go run .` to start the server

## Features

### ClickUp + Clubhouse

Middle Manager automates information updates between [ClickUp](https://clickup.com/) and [Clubhouse](https://clubhouse.io/).

[ClickUP API](https://clickup20.docs.apiary.io/#)
[Clubhouse API](https://clubhouse.io/api/rest/v3/)

#### Create Clubhouse Epic

When a ClickUp task is moved to the status Ready For Development an Epic will be created in Clubhouse.
