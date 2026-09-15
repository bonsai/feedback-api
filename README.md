# feedback-api

`bonsai/meikyoku-janken` のデザインフィードバックを GitHub Issues に登録する REST API。

## Architecture

```text
GitHub Pages
    ↓ POST /feedback
Vercel
    ↓ GitHub REST API
bonsai/meikyoku-janken Issues
```

- **Go**: REST API
- **Vercel**: deploy / runtime
- **GitHub Issues**: feedback store
- `GITHUB_TOKEN` は Vercel の Environment Variables に置き、ブラウザには渡さない

Vercel は 2026-04-02 から Go API backend の zero-configuration deployment に対応しているため、`main.go` をそのままデプロイできる。citeturn0search0turn0search1

## Endpoints

- `GET /health`
- `POST /feedback`

```json
{
  "day": 3,
  "theme": "色",
  "category": "デザイン",
  "message": "お題をもっと大きくしてほしい"
}
```

## Vercel Deploy

1. Vercel で `bonsai/feedback-api` を Import
2. Environment Variables に `GITHUB_TOKEN` を登録
3. Deploy

CLIなら:

```bash
npm i -g vercel
vercel --prod
```

VercelのGit連携なら `main` への push でも自動デプロイできる。citeturn0search1

## Local

```bash
go run main.go
```

```bash
curl http://localhost:8080/health
```

```bash
curl -X POST http://localhost:8080/feedback \
  -H 'Content-Type: application/json' \
  -d '{"day":3,"theme":"色","category":"デザイン","message":"お題をもっと大きくしてほしい"}'
```
