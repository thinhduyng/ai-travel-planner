# Contract: `POST /v1/itineraries`

Generate a day-by-day travel itinerary for one destination and date range.

Status: **v1, stable**. Any change that breaks a rule below is a breaking change and needs `/v2`.
Each rule has a number (§x.y). The contract tests in `api/contract/` are named after those numbers.
The schema half of this contract, used for code generation, is [`create-itinerary.openapi.yaml`](create-itinerary.openapi.yaml). If the two disagree, this file wins.

## 1. Request

`POST /v1/itineraries` with `Content-Type: application/json`.

```json
{
  "destination": "Da Lat, Vietnam",
  "startDate": "2026-12-20",
  "endDate": "2026-12-22",
  "travelers": 2,
  "budget": { "amount": 5000000, "currency": "VND" },
  "interests": ["food", "nature"]
}
```

| Field         | Type             | Required | Rule                                                                               |
| ------------- | ---------------- | -------- | ---------------------------------------------------------------------------------- |
| `destination` | string           | yes      | §1.1 1–100 characters after trimming surrounding whitespace                        |
| `startDate`   | string           | yes      | §1.2 a real calendar date in `YYYY-MM-DD` format                                   |
| `endDate`     | string           | yes      | §1.2 same format; §1.3 not before `startDate`; §1.4 trip is at most **14 days**    |
| `travelers`   | integer          | no       | §1.5 1–20; defaults to `1` when omitted                                            |
| `budget`      | object           | no       | §1.6 `amount` is a positive integer; `currency` is 3 upper-case letters (ISO 4217) |
| `interests`   | array of strings | no       | §1.7 at most 10 items, each 1–30 characters; defaults to `[]`                      |

§1.8 Unknown top-level fields are rejected (400). They are not ignored.

The trip length is `endDate - startDate + 1` days, so `startDate == endDate` is a one-day trip.

## 2. Success: `201 Created`

§2.1 Status `201`, and the `Location` header is `/v1/itineraries/{id}`.

§2.2 The body:

```json
{
  "id": "6f1c2a8e-3b7d-4c55-9a51-2f0e8b7d1c34",
  "destination": "Da Lat, Vietnam",
  "startDate": "2026-12-20",
  "endDate": "2026-12-22",
  "travelers": 2,
  "budget": { "amount": 5000000, "currency": "VND" },
  "interests": ["food", "nature"],
  "days": [
    {
      "day": 1,
      "date": "2026-12-20",
      "activities": [
        {
          "time": "09:00",
          "title": "Xuan Huong Lake walk",
          "description": "..."
        }
      ]
    }
  ],
  "createdAt": "2026-11-25T06:30:00.000Z"
}
```

- §2.3 `id` is a UUID (v4).
- §2.4 `destination` is the trimmed input. `travelers` and `interests` are echoed with their defaults applied. `budget` is echoed, or `null` when it was not sent.
- §2.5 `days` has exactly one entry per trip day, in order. `day` starts at `1`, and `date` is `startDate + (day - 1)`.
- §2.6 Each day has at least one activity. `time` is `HH:MM` (24h). `title` is a non-empty string. `description` is optional.
- §2.7 Within a day, activities are sorted by `time` ascending.
- §2.8 `createdAt` is an ISO 8601 UTC timestamp.

## 3. Errors

§3.1 Every error has this body shape:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Human-readable summary",
    "details": [{ "field": "endDate", "issue": "must not be before startDate" }]
  }
}
```

`details` is always an array, and it is empty when there is nothing field-specific to say.

| Status | `error.code`             | When                                                                                                                       |
| ------ | ------------------------ | -------------------------------------------------------------------------------------------------------------------------- |
| 400    | `VALIDATION_ERROR`       | §3.2 Any rule in §1 is broken. `details` has one entry per broken field; `field` is the dotted path (e.g. `budget.amount`) |
| 400    | `MALFORMED_JSON`         | §3.3 The body is not valid JSON                                                                                            |
| 415    | `UNSUPPORTED_MEDIA_TYPE` | §3.4 `Content-Type` is not `application/json`                                                                              |
| 502    | `PLANNER_UNAVAILABLE`    | §3.5 The AI planner failed or returned an unusable plan. No partial itinerary is returned                                  |

§3.6 Error responses never include stack traces, provider names, API keys or raw model output.

## 4. Non-goals for v1

No auth, no persistence guarantee (`GET /v1/itineraries/{id}` is not part of this contract), no pricing per activity.
