# Device

## API Routes

### All devices

**GET: `/api/v2/devices`**

Response:

```json
[
  {
    "ip": "127.0.0.1",
    "company": {
      ...
    },
    "branch": {
      ...
    },
    "user": {
      ...
    },
    "type": "PC"
  },
  ...
]
```

### Create device

**POST: `/api/v2/devices`**

> **SHOULDN'T BE USED DIRECTLY FROM CLIENT. ONLY FOR MANAGEMENT PURPOSES**

> _For client please refer to [Authorization](auth.md#Register)_

Request:

```json
{
  "ip": "127.0.0.1",
  "company": {
    ...
  },
  "branch": {
    ...
  },
  "user": {
    ...
  },
  "type": "PC"
}
```

Response:

```json
{
  "ip": "127.0.0.1",
  "company": {
    ...
  },
  "branch": {
    ...
  },
  "user": {
    ...
  },
  "type": "PC"
}
```

## Declaration

```go
// models/v2/device.go
type Device struct {
	IP      string  `json:"ip"`
	Company Company `json:"company"`
	Branch  Branch  `json:"branch"`
	User    User    `json:"user"`
	Type    string  `json:"type"`
}
```

## JSON Schema

```json
{
  "ip": "127.0.0.1",
  "company": {
    "id": "be60ac69-84db-437e-bab5-9ce55cba7f0b",
    "name": "Development Environment Company",
    "slug": "devenvcomp"
  },
  "branch": {
    "id": "be60ac69-84db-437e-bab5-9ce55cba7f0b",
    "name": "Dev Branch",
    "company_id": "be60ac69-84db-437e-bab5-9ce55cba7f0b",
    "address": "Development Space",
    "contacts": "dev@env.com",
    "description": "Master Branch"
  },
  "user": {
    "id": "be60ac69-84db-437e-bab5-9ce55cba7f0b",
    "name": "user",
    "phone": "1234"
  },
  "type": "PC"
}
```
