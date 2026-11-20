# Disney Infinity Global (Wii U) — NEX server

> If you use this source base, please credit us, as well as the other contributors mentioned in this repository.

Built on the [Pretendo Network](https://github.com/PretendoNetwork) NEX
libraries (`nex-go`, `nex-protocols-go`, `nex-protocols-common-go`) — the
same stack as this org's other servers.

## Games

Several games run side by side in one process. Each has its own access key,
NEX version, ports and environment prefix, defined in `version/`.

| Game | Config file | Env prefix | Access key | NEX | Auth / secure UDP ports |
| --- | --- | --- | --- | --- | --- |
| Disney Infinity 1.0 | `version/dsi1.go` | `PN_DSI1_` | `2ef57176` | 3.5.1 | 29100 / 29101 |
| Disney Infinity 3.0 | `version/dsi3.go` | `PN_DSI3_` | `83a601e5` | 3.9.1 | 29200 / 29201 |

Adding a game means adding one file in `version/` (with `register(NewGame(...))`)
and a new env prefix. Startup fails if two games share a port.

## Quick start

Requires Go 1.27+. PostgreSQL is optional (see below).

```powershell
Copy-Item example.env .env
Copy-Item settings.example.json settings.json
go run .
```

You should see both games start:

```
[Disney Infinity 1.0] Authentication server listening on UDP 29100
[Disney Infinity 1.0] Secure server listening on UDP 29101
[Disney Infinity 3.0] Authentication server listening on UDP 29200
[Disney Infinity 3.0] Secure server listening on UDP 29201
```

Open the four UDP ports in your firewall if consoles connect from another machine.

## Configuring `.env`

The server reads `.env` in the working directory (or the process environment).
Every game variable is `<PREFIX>_<NAME>`, where the prefix is `PN_DSI1` or
`PN_DSI3`. Replace `x` below by `1` or `3`.

| Variable | Required | Description |
| --- | --- | --- |
| `GAMES` | no | Comma-separated slugs to run, e.g. `disney-infinity-1.0`. Default: all games. |
| `PN_DSIx_SECURE_HOST` | yes | Host the console is told to connect to for the secure server. |
| `PN_DSIx_AUTH_PORT` | no | Overrides the default authentication UDP port. |
| `PN_DSIx_SECURE_PORT` | no | Overrides the default secure UDP port. |
| `PN_DSIx_KERBEROS_PASSWORD` | no | Password of the system accounts (PID 1 and 2). Change it from the default. |
| `PN_DSIx_LOCAL_MODE` | no | `1` = local mode (see below). Anything else = token mode. |
| `PN_DSIx_SETTINGS_PATH` | local mode | Path to `settings.json`. Default: `settings.json`. |
| `PN_DSIx_NEX_TOKEN_AES_KEY` | token mode | 64 hex characters (AES-256 key shared with the account server). |
| `PN_DSIx_NEX_PASSWORD_SECRET` | token mode | At least 32 bytes in hex, shared with the account server. |
| `PN_DSIx_POSTGRES_URI` | no | PostgreSQL connection string. Enables matchmaking and utility. |
| `HTTP_PORT` | no | Port of the optional companion HTTP server (default `8081`). |
| `HTTP_POSTGRES_URI` | no | Database used by the HTTP `/health` check. |

Minimal local test `.env`:

```env
PN_DSI1_SECURE_HOST=192.168.1.20
PN_DSI1_LOCAL_MODE=1
PN_DSI3_SECURE_HOST=192.168.1.20
PN_DSI3_LOCAL_MODE=1
```

## Protocols

Implemented: Ticket Granting, Secure Connection, NAT Traversal, Matchmaking,
Matchmaking Ext, Matchmake Extension, Utility (the last four need a database).

Not implemented yet: Ranking and Storage Manager (declared for Disney Infinity
1.0, a warning is printed), and Datastore for Disney Infinity 3.0.

Every incoming RMC request is logged with its protocol and method name, which
helps discovering what a game calls.

## Development

```powershell
go vet ./...
go test ./...
docker build -t disney-infinity-nex-global .
```

## Troubleshooting

- **Console loops on CONNECT / error 106-0502**: the NEX version or access key
  of the game is wrong. The NEX version and access key are set per game in
  `version/dsi1.go` / `version/dsi3.go` (Disney Infinity 1.0 is NEX 3.5.1);
  NEX 3.5.0 and newer use structure headers, and below 3.2.0 the legacy
  connection signature is used. Verify these first.
- **`UDP port N is used by both ...`**: two games have the same port, fix
  `PN_DSIx_AUTH_PORT` / `PN_DSIx_SECURE_PORT`.
- **`... environment variable not set`**: a required variable is missing for a
  game you are running; use `GAMES=` to run only one game.
- **`ping postgres: connection refused`**: the URI is set but PostgreSQL is not
  running. Start it or comment the URI out.
- **Login fails in local mode**: the PID is not in `settings.json`
  (`username` must equal the PID).

## License & Copyright

This project is licensed under AGPLv3.0. Copyright © 2026 Protarium Network.

Deployed and maintained as part of the Protarium Network Wii U online service revival project. Derived from Pretendo Network’s original codebase (PretendoNetwork). Copyright (C) Pretendo Network contributors.
