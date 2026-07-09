# Running the test suite

The suite runner (`zit-suite`) executes every test inside a disposable Docker
container, so no test touches the host filesystem or git config.

## Build the image

From the repo root:

```sh
# Go implementation
docker build -t zit-go-test -f docker/Dockerfile.golang .

# Rust implementation
docker build -t zit-rust-test -f docker/Dockerfile.rust .
```

## Run the tests

```sh
go run ./cmd/zit-suite zit-go-test

go run ./cmd/zit-suite zit-rust-test
```
