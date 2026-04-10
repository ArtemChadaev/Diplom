# Backend — Additional Testing Notes

← [Back to Main README](../../README.md) | [Testing Strategy →](../Diplom/testing.md)

---

## Docker-Based Testing Requirement

All backend tests must also pass the Docker container build without errors.

The full test environment requires:
- A running `backend` Docker container
- A running `valkey` container
- A running `postgres` container

> **Preferred approach:** Fully restart (re-create) the backend container on each test run to ensure a clean state.

```bash
# Rebuild and restart backend container before running tests
docker compose --profile backend-all up -d --build --force-recreate backend
go test ./...
```

---

## TODO

- [ ] Define integration test suite structure
- [ ] Set up `testcontainers-go` for isolated DB in CI
- [ ] Add coverage reporting to CI pipeline (target: ≥ 70%)