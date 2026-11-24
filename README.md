# AI Travel Planner

Generates day-by-day travel itineraries with an LLM.

## Run the API

```bash
cp .env.example .env
make run            # http://localhost:8080/health
```

Requires Go (version in `api/go.mod`) and `make`.

## Checks (same commands as CI)

```bash
make lint
make fmt-check
make typecheck
make generate-check
make test
make contract
make check          # all of them except generate-check
```

- Assistant rules: [`CLAUDE.md`](CLAUDE.md) · permissions: [`.claude/settings.json`](.claude/settings.json)
- Endpoint contracts: [`docs/contract/`](docs/contract/)
- API code generated from the OpenAPI spec: `make generate`
