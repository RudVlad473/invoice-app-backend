# Invoice App — Backend

Go service for managing invoices, built to explore testing backend logic against
real infrastructure instead of mocked dependencies.

## Why this exists

Most Go backend tutorials mock the database to keep tests fast, but that hides bugs
in how your code actually talks to the database (marshaling, update expressions,
query shape). This project was built to test the opposite approach: spin up a real
local DynamoDB instance via Docker for every test run, and assert against real
read/write behavior instead of a mock's promises.

## Architecture & key decisions

- **Repository layer talks directly to DynamoDB** using the AWS SDK v2, with
  manual attribute marshaling/unmarshaling rather than an ORM, to stay close to
  how DynamoDB's update expressions actually behave
- **Partial updates** are built dynamically from whatever fields are passed,
  rather than requiring a full object on every update call
- **Test setup spins up a disposable local DynamoDB container** via Docker
  for every test run (`TestMain`), so tests exercise real query/update behavior
  instead of asserting against mocks

## Tech stack

Go, AWS SDK v2, DynamoDB (local via Docker for testing), Docker, GitHub Actions

## What's not done yet

This repo currently proves out the data layer and test infrastructure — the
HTTP handlers are scaffolded but the CRUD endpoints aren't wired up yet (no
POST/PUT/DELETE routes live). The repository layer itself is fully functional
and covered by integration tests; the API surface is the next piece to build out.

## Running the tests

Requires Docker running locally (used automatically in CI via GitHub Actions).

```bash
go test ./...
```
