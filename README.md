# Invoice App — Backend

A Go service for managing invoices, built to test business logic against a **real** local
DynamoDB instance instead of a mocked one.

## Why this exists

Mocking the database keeps unit tests fast but hides bugs in how code actually talks to it —
marshaling, update-expression shape, partial-update semantics. This repo tests the opposite: every
test run boots a disposable local DynamoDB container via Docker (`TestMain`) and asserts against
real read/write behavior.

## How to run it

Requires Docker running locally.

```bash
go test ./...
```

8 top-level test functions (14 `t.Run` subtests, counted directly from `tests/invoice_service_test.go`)
covering find-by-id, save, count, delete, full update, and the item-level add/remove/update
operations — each against the real local DynamoDB container, not a mock.

## Stack

| Concern     | Choice                                                                                                          | Why                                                                                                          |
| ----------- | ---------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| Language    | Go 1.25                                                                                                         | —                                                                                                            |
| Persistence | DynamoDB, AWS SDK v2, manual `attributevalue`/`expression` marshaling (no ORM)                                  | Stays close to how DynamoDB's update expressions actually behave, rather than hiding them behind an abstraction       |
| Validation  | `go-playground/validator` with a custom `status` rule                                                          | DTOs validate invoice status against the enum, not just its Go type                                                   |
| Testing     | Standard `testing` + `dynamodb-local` via Docker (custom `testingutils` package), `google/go-cmp`, `gofakeit`  | One shared `TestMain` boots/tears down the container so every package's tests hit real DynamoDB behavior instead of duplicating setup |

## What's implemented

The repository layer (`invoice.Repository`) is complete and covered by integration tests:
`FindById`, `Save`, `CountAll`, `DeleteById`, `UpdateById` (partial updates built from whatever
fields are passed, not a full-object requirement), and item-level `AddItemByInvoiceId` /
`RemoveItemByInvoiceId` / `UpdateItemByInvoiceId`.

## What's not done yet

The HTTP layer is scaffolded but not wired to the repository: `invoice.Handler` only routes `GET`,
and its handler body is empty — no request reaches the repository logic above yet, and there's no
`POST`/`PUT`/`DELETE`/`PATCH` route for any of it. There's also no CI workflow yet; `go test ./...`
is run locally only. The data layer and its tests are the finished part of this repo; the API
surface is the next piece to build.

## License

MIT — see [LICENSE](LICENSE).
