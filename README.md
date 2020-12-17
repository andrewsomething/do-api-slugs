# DigitalOcean API Slugs

[https://slugs.do-api.dev/](https://slugs.do-api.dev/)

![Screenhot](https://i.imgur.com/etNCvLU.png)


## Project Details

The frontend is provided by a Vue.js powered static site. The backend Go service found in the `api/` directory. It proxies the DigitalOcean API so that an API token is not required on the frontend and set a `Cache-Control` header so the responses are appropriately cached by the CDN.

### Local Development

A Docker Compose file is provide for local development. To build and run both components, use:

    docker-compose up --build

### Deploy to DigitalOcean

In production, this site is host on DigitalOcean's App Platform. You can deploy it yourself using:

[![Deploy to DO](https://mp-assets1.sfo2.digitaloceanspaces.com/deploy-to-do/do-btn-blue.svg)](https://cloud.digitalocean.com/apps/new?repo=https://github.com/andrewsomething/do-api-slugs/tree/deploy-to-do)

Find the app spec in [.do/deploy.template.yaml](.do/deploy.template.yaml).