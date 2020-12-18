# AOEU Go Template Repo

![CI](https://github.com/theartofeducation/template-repo/workflows/CI/badge.svg?branch=main)

A template repo for quickly starting Go projects.

## What's Included

* Commit linting through NPM
* [gorilla/mux](https://github.com/gorilla/mux)
* [ent](https://entgo.io/)
* [golangci-lint](https://github.com/golangci/golangci-lint)
* [godotenv](https://github.com/joho/godotenv)
* [Sentry](https://github.com/getsentry/sentry-go)
* [Logrus](https://github.com/Sirupsen/logrus)

## How To Use This Template

1. Go to the ***Create a new repository*** page on GitHub
1. Select this repo as the template in the ***Repository template*** field
1. Check the ***Include all branches*** box
1. Select ***theartofeducation*** as the ***Owner*** of the repository
1. Add the repository name
1. Add a short description for the repository
1. Select ***Public*** or ***Private*** as appropriate for the repository
1. Click the ***Create repository*** button
1. Update files and repository information as needed
1. Rebase `develop` on `main` to ensure that when you create a PR in the future,
   GitHub will allow it from `develop`
    1. `git rebase origin/main`
    1. `git push -f`
1. Follow these steps to set up branch protection rules for `develop` and `main` in the new repository (manual setup)
    1. Go to the repo page in GitHub
    1. Go to the ***Settings*** page
    1. Go to the ***Branches*** section
    1. Under ***Branch protection rules*** click the ***Add rule*** button
    1. Type the name of the branch in ***Branch name pattern***
    1. Select ***Require pull request reviews before merging***
        1. Set ***Required approving reviews*** to 2
        1. Select ***Dismiss stale pull request approvals when new commits are pushed***
        1. Select ***Require review from Code Owners***
    1. Select ***Require status checks to pass before merging***
        1. Select ***Require branches to be up to date before merging***
        1. Check any CI pipelines needed
    1. Select ***Require signed commits***
    1. Select ***Include administrators***
    1. Click the ***Save changes*** button
1. Create a new feature branch to work off
1. Set up commit linting
    1. Run `yarn install`
1. Setup golangci-lint
    1. Run `brew install golangci/tap/golangci-lint`
1. Updated the module path in `go.mod`
1. Update the `README.md`
1. Commit and merge changes

## Running The Application

Two Docker containers are setup. The first builds the application while the second will run the executable.

1. `docker-compose up --build`

## Testing

Tests should be written with the application being a "black box" with no direct access. Test files go under the ***./tests*** directory which should follow the directory structure of the application.
