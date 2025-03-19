# Lighthouse Server

This Docker container provides a single endpoint to run Google Lighthouse on a website.
It will run Lighthouse and return the result via HTTP.

## Usage

### Running the container

```bash
docker run -p 8080:80 -d ghcr.io/govigilant/lighthouse-server:latest
```

### Running Lighthouse

```bash
 curl -X POST http://localhost:8080/lighthouse \
     -H "Content-Type: application/json" \
     -d '{
           "website": "https://govigilant.io",
           "callback_url": "https://govigilant.io/api/lighthouse/callback"
         }'
```
