# DigitalOcean API Slugs

[https://slugs.do-api.dev/](https://slugs.do-api.dev/)

![Screenhot](https://i.imgur.com/etNCvLU.png)


## Project Details

The frontend is provided by a Vue.js powered static site. The backend is made up of serverless functions found in the `functions/` directory. These proxy the DigitalOcean API so that an API token is not required on the frontend and set a `Cache-Control` header so the responses are appropriately cached by the CDN.

#### Install frontend dependencies
```
npm install
```

#### Run development server
```
npm run serve
```

#### Build for production
```
npm run build
```
