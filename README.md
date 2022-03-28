<p align="center">
    <img src="https://raw.githubusercontent.com/Qovery/public-resources/master/qovery%20logo%20horizontal%20without%20margin.png" alt="Qovery logo" />
</p>

<p align="center">Deploy your apps on any Cloud providers in just a few seconds ⚡</p>

<p align="center">
<img src="https://github.com/Qovery/qovery-github-action/actions/workflows/test.yml/badge.svg?style=flat-square" alt="Tests">
</p>

<h3 align="center">The simplest way to deploy your apps in the Cloud</h3>

<br />

# [Qovery](https://www.qovery.com/) GitHub Actions

**Qovery GitHub Actions** is a GitHub Actions plugin allowing Qovery users to integrate Qovery within their CI nicely.

- Website: https://www.qovery.com
- Qovery documentation: https://hub.qovery.com/docs

**Please note**: We take Qovery security and our users' trust very seriously. If you believe you have found a security issue in Qovery, please responsibly disclose by contacting us at security@qovery.com.

## ✅ Requirements
- A **Qovery** account. [Sign up now](https://start.qovery.com/) if you don't have any account yet.

## 📖 Installation
- Create an API key: [how to generate your API token?](https://hub.qovery.com/docs/using-qovery/interface/cli/#generate-api-token)
- Setup a secret named `QOVERY_API_TOKEN` within your repository `Secrets` section and set its value with output of the previous step.

## 🔌 Usage
- Add a new job to your GitHub workflow (e.q. adding a step after your `tests`) using `Qovery/qovery-action` action.

### Deploy your application

```
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    name: Deploy on Qovery
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Deploy on Qovery
        uses: Qovery/qovery-action@main
        id: qovery
        with:
          qovery-organization-name: [YOUR_QOVERY_ORGANIZATION_NAME (CASE SENSITIVE)]
          qovery-project-name: [YOUR_QOVERY_PROJECT_NAME (CASE SENSITIVE)]
          qovery-environment-name: [APPLICATION_QOVERY_ENVIRONMENT_NAME (CASE SENSITIVE)]
          qovery-application-names: [APPLICATION_QOVERY_APPLICATION_NAME_1,APPLICATION_QOVERY_APPLICATION_NAME_2] # Comma-separated string of names (case sensitive)
          qovery-application-commit-id: [APPLICATION_QOVERY_APPLICATION_COMMIT_ID]
          qovery-api-token: ${{secrets.QOVERY_API_TOKEN}}
```

### Deploy a database

```
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    name: Deploy on Qovery
    steps:
      - name: Checkout
        uses: actions/checkout@v3
      - name: Deploy on Qovery
        uses: Qovery/qovery-action@main
        id: qovery
        with:
          qovery-organization-name: [YOUR_QOVERY_ORGANIZATION_NAME (CASE SENSITIVE)]
          qovery-project-name: [YOUR_QOVERY_PROJECT_NAME (CASE SENSITIVE)]
          qovery-environment-name: [APPLICATION_QOVERY_ENVIRONMENT_NAME (CASE SENSITIVE)]
          qovery-database-name: [APPLICATION_QOVERY_DATABASE_NAME (CASE SENSITIVE)]
          qovery-api-token: ${{secrets.QOVERY_API_TOKEN}}
```
