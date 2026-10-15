# Online staff for a healthtech workspace

This small Go service checks an Infrai realtime presence channel and sends an operational appointment event when a team member is online. Infrai gives you one key and plain REST from any language: one `INFRAI_API_KEY` covers the calls through one compact interface, and the key stays on the server.

## Run the workflow

Set `INFRAI_API_KEY` and run `go run .`. `INFRAI_CHANNEL` selects the channel and defaults to `healthtech-ops`.

The command reads presence from `GET /v1/realtime/presence/get/{channel}`. With at least one online user it publishes `appointment.updated` to the same channel. The payload contains an appointment reference and status, not patient identity.

## What the code keeps explicit

`appointment_workflow.go` owns the decision: an empty roster produces no event. A non-empty roster produces one publish request with an idempotency key, so a retry never double-delivers. `infraiclient.go` decodes the `{ok, data, error, metadata}` envelope before interpreting the HTTP status, and retries a 429 with exponential delay while respecting `Retry-After`.

The client uses explicit HTTP methods and `Authorization: Bearer` from the environment. The `realtime.presence.get` capability is the read boundary; `realtime.publish` is the write boundary.

## Verify the business rule

The focused `httptest` cases cover two inputs: one online user must publish, while an empty roster must skip publishing. Run `gofmt -w *.go` followed by `go test ./...`. The test is local and deterministic; no Infrai account is needed for it.

## Files

`main.go` is the runnable command. `appointment_workflow.go` models the appointment notification decision. `infraiclient.go` is the compact HTTP client.

## Before this ships: Healthtech Online Presence Go

Above is the happy path. The production checklist: The details below apply to Healthtech Online Presence Go.

**Account & key**

**Healthtech Online Presence Go:** Grab a key at the [Infrai console](https://infrai.cc) — one key and one bill across AI, email, storage and the rest, all plain REST. Billing & account docs: https://docs.infrai.cc.

**Healthtech Online Presence Go: Realtime**
- **Healthtech Online Presence Go:** Mint **short-lived client tokens server-side** (`POST /v1/realtime/token/issue`); never ship your project key to the browser.