# Users
### `POST /api/users`
Create a new user.

Request:
```json
{
  "email": "example@foo.com",
  "password": "bar"
}
```

Response:
```json
{
  "id": "some_uuid",
  "created_at": "2025-11-16T11:54:32.868529Z",
  "updated_at": "2025-11-16T11:54:32.868529Z",
  "email": "example@foo.com",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false
}
```

### `PUT /api/users`
Update a user. Requires authorization header (see `/api/login`)

Request:
```json
{
  "email": "example@foo.com",
  "password": "bar"
}
```

Response:
```json
{
  "id": "some_uuid",
  "created_at": "2025-11-16T11:54:32.868529Z",
  "updated_at": "2025-11-16T11:54:32.868529Z",
  "email": "example@foo.com",
  "token": "",
  "refresh_token": "",
  "is_chirpy_red": false
}
```

# Chirps
### `POST /api/chirps`
Create a chirp. Max size limited to 140 characters.

Requires authorization (see `/api/login`).

Request:
```json
{
  "body": "Example chirp"
}
```

Response:
```json
{
  "id": "some_chirp_uuid",
  "created_at": "2025-11-16T11:54:32.868529Z",
  "updated_at": "2025-11-16T11:54:32.868529Z",
  "body": "Example chirp",
  "user_id": "some_user_uuid"
}
```

### `GET /api/chirps`
Get chirps. Supports filtering by `author_id`. Supports sorting by creation date `asc|desc`.

Response:
```json
[
  {
    "id": "some_chirp_uuid",
    "created_at": "2025-11-16T11:54:32.868529Z",
    "updated_at": "2025-11-16T11:54:32.868529Z",
    "body": "Example chirp 1",
    "user_id": "some_user_uuid"
  },
  {
    "id": "some_chirp_uuid",
    "created_at": "2025-11-17T11:54:32.868529Z",
    "updated_at": "2025-11-17T11:54:32.868529Z",
    "body": "Example chirp 2",
    "user_id": "some_user_uuid"
  }
]
```
### `GET /api/chirps/{chirpID}`
Get chirp specified by its `id`.

Response:
```json
{
  "id": "some_chirp_uuid",
  "created_at": "2025-11-16T11:54:32.868529Z",
  "updated_at": "2025-11-16T11:54:32.868529Z",
  "body": "Example chirp",
  "user_id": "some_user_uuid"
}
```

### `DELETE /api/chirps/{chirpID}`
Delete a chirp specified by its `id`.

Returns no payload.

# Authentication
### `POST /api/login`
Log in using email and password.

The `token` field contains JWT which can be used for authorized
endpoints using header `Authorization: Bearer <jwt>`. JWT is valid for 1
hour.

The `refresh_token` field is a string which can be used for `/api/refresh` endpoint.
The token is valid for 60 days, unless revoked (see `/api/revoke`).

Request:
```json
{
  "email": "example@foo.com",
  "password": "bar"
}
```

Response:
```json
{
  "id": "some_uuid",
  "created_at": "2025-11-16T11:54:32.868529Z",
  "updated_at": "2025-11-16T11:54:32.868529Z",
  "email": "example@foo.com",
  "token": "jwt_token_string",
  "refresh_token": "refresh_token_string",
  "is_chirpy_red": false
}
```

### `POST /api/refresh`
Get a new JWT valid for 1 hour.

The endpoint doesn't require any payload but requires the `refresh_token` in
the header: `Authorization: Bearer <referesh_token>`.

Response:
```json
{
  "token": "jwt_token_string"
}
```

### `POST /api/revoke`
Revoke the refresh token, i.e. it can no longer be used to receive a JWT.

The endpoint doesn't require any payload but requires the `refresh_token` in
the header: `Authorization: Bearer <referesh_token>`.

The endpoint doesn't return anything in the payload.

# Health
### `GET /api/healthz`
Returns `200` and `OK`.

# Admin
### `GET /admin/metrics`
Shows the number of hits for the `/app`.

### `GET /admin/reset`
Resets the number of hits for the `/app`. Additionally, if the env var
`PLATFORM` is set to `dev`, it deletes all the users in the database.

# Webhooks
### `POST /api/polka/webhooks`
Process webhooks by the Polka provider. For more details consult the code.
