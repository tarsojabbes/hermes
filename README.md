# Hermes

Lightweight Pub/Sub Event Streaming Platform, built with Go.

### How to run Hermes

1. Copy `.env.example` to `.env`

```sh
cp .env.example .env
```

2. Build Hermes with Go

```sh
go build
```

3. Start a ScyllaDB instance with Docker

```sh
docker run --name hermes-scylla --hostname some-scylla -d scylladb/scylla --smp 1
```

4. Run Hermes

```sh
go run hermes
```