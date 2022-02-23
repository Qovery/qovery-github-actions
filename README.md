# [Qovery](https://www.qovery.com/) GitHub Action
<h3 align="center">The simplest way to deploy your apps in the Cloud</h3>
<p align="center">Deploy your apps on any Cloud providers in just a few seconds ⚡</p>

<p align="center">
<img src="https://github.com/Qovery/qovery-github-action/actions/workflows/test.yml/badge.svg?style=flat-square" alt="Tests">
</p>

<br />

<p align="center">
    <img src="https://raw.githubusercontent.com/Qovery/public-resources/master/qovery_kubernetes_cloudproviders.svg" height="450px" alt="Qovery stack on top of Kubernetes and Cloud providers" />
</p>

**Qovery Github Action** is a Github Action plugin allowing Qovery users to integrate Qovery within their CI nicely.

- Website: https://www.qovery.com
- Qovery documentation: https://hub.qovery.com/docs

**Please note**: We take Qovery's security and our users' trust very seriously. If you believe you have found a security issue in Qovery, please responsibly disclose by contacting us at security@qovery.com.

## ✅ Requirements
- A **Qovery** account. [Sign up now](https://start.qovery.com/) if you don't have any account yet.

## 📖 Installation
- Create an API key: [how to generate your API token?](https://hub.qovery.com/docs/using-qovery/interface/cli/#generate-api-token)
- Setup a secret named `QOVERY_API_TOKEN` within your repository `Secrets` section and set its value with output of the previous step.

## 🔌 Usage
- Add a new job to your Github workflow (e.q. adding a step after your `tests`) using `Qovery/qovery-action` action.
```
on: [push]

jobs:
  deploy:
    runs-on: ubuntu-latest
    name: Deploy on Qovery
    steps:
      - name: Checkout
        uses: actions/checkout@v2
      - name: Deploy on Qovery
        uses: Qovery/qovery-action@main
        id: qovery
        env:
          QOVERY_API_TOKEN: ${{ secrets.QOVERY_API_TOKEN }}
        with:
          qovery-organization-id: [YOUR_QOVERY_ORGANIZATION_ID]
          qovery-environment-id: [APPLICATION_QOVERY_ENVIRONMENT_ID]
          qovery-application-id: [APPLICATION_QOVERY_APPLICATION_ID]
```
