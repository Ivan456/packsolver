# Golang Pack Solver

A minimal Go service that computes which packs to ship to satisfy an order under these rules:

1. Only whole packs can be sent.
2. Minimize total items shipped (must be >= requested amount).
3. Among solutions with minimal total items, minimize number of packs.

This repository is intentionally small and easy to run locally.

---

## Run locally (requires Go 1.18+)

```bash
# install deps
go mod tidy

# run tests
go test ./...

# start service
go run main.go
```

Service runs at: [http://localhost:8080](http://localhost:8080)

Server will listen on `http://localhost:8080`.
Open `http://localhost:8080` to use the simple UI. Or call the API directly:

POST /api/pack
Content-Type: application/json

{
  "amount": 500000,
  "sizes": [23,31,53]
}

Response:

200 OK
{
  "shipped_total": 500000,
  "packs": {"23":2,"31":7,"53":9429},
  "pack_count": 9438
}

Notes:
- The pack sizes array is fully configurable and can be changed at runtime.
- The algorithm uses dynamic programming to find the minimal shipped total >= amount and, among those, the minimal number of packs.

---

## Run locally with Docker from existing image (no Go required)

The image has been pushed to GitHub Container Registry:

```bash
docker pull ghcr.io/ivan456/packsolver:latest
docker run --rm -p 8080:8080 ghcr.io/ivan456/packsolver:latest
```


## CI/CD with GitHub Actions

This repo uses a manual pipeline to build and push Docker images.

### Run Pipeline
1. Go to the **Actions** tab in GitHub.
2. Select **Manual Build and Push** workflow.
3. Click **Run workflow**.

The workflow will:
- Run Go tests
- Build Docker image
- Push image to [GitHub Container Registry](https://ghcr.io)


## Build image locally

```bash
docker build -t packsolver:latest .
```

### Run container

```bash
docker run --rm -p 8080:8080 packsolver:latest
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

---

## Publish to GitHub Container Registry

1. **Login** (replace USERNAME with your GitHub username):

```bash
echo $GITHUB_TOKEN | docker login ghcr.io -u USERNAME --password-stdin
```

*(generate a PAT in GitHub with `write:packages` and `read:packages` permissions and export as `GITHUB_TOKEN`)*

2. **Build & push**:

```bash
docker build -t ghcr.io/ivan456/packsolver:latest .
docker push ghcr.io/ivan456/packsolver:latest
```
