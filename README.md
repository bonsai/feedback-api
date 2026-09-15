# feedback-api

Design feedback REST API for `bonsai/meikyoku-janken`.

## Stack

- Go: deployable HTTP REST API
- Python: GitHub API client / local automation
- GitHub Issues: feedback store

## Endpoints

- `GET /health`
- `POST /feedback`

Example:

```json
{"day":3,"theme":"色","category":"デザイン","message":"お題をもっと大きくしてほしい"}
```

The API keeps `GITHUB_TOKEN` server-side and creates an Issue through GitHub REST API. Do not put the token in the frontend.

## Deploy

The Go binary is platform-neutral and can be deployed immediately to Cloud Run, Fly.io, Render, Railway, or another container/HTTP host. Set `GITHUB_TOKEN` and `PORT`.
