# disney-infinity-nex-global

NEX (Quazal Rendez-Vous) game server for Disney Infinity, built on
[nex-go](https://github.com/PretendoNetwork/nex-go). Game parameters live in
`version/`, so adding a game means adding one file there.

## Layout

| Path | Role |
| --- | --- |
| `main.go`, `init.go` | Entry point; loads `.env`, picks the game, connects to Postgres |
| `nex/` | Authentication (ticket granting) and secure servers |
| `globals/` | Shared state, accounts, login validation |
| `version/` | Per-game config (access key, NEX version, protocols) |
| `database/` | Postgres connection |
| `http/` | Optional companion HTTP server with `/health` (`go run ./http`) |

Implemented protocols: Ticket Granting, Secure Connection, NAT Traversal,
Matchmaking, Matchmaking Ext, Matchmake Extension, Utility. Ranking and
Storage Manager are declared for Disney Infinity 1.0 but not implemented yet;
Disney Infinity 3.0 has no datastore yet either.

## Run

Requires Go 1.27+ and PostgreSQL.

```powershell
Copy-Item example.env .env
Copy-Item settings.example.json settings.json   # local mode only
go run .
```

Disney Infinity 1.0 and 3.0 run side by side in the same process, each with
its own access key and ports (see `version/dsi1.go`, `version/dsi3.go`):

| Game | Env prefix | Auth / secure UDP ports |
| --- | --- | --- |
| Disney Infinity 1.0 | `PN_DSI1_` | 29100 / 29101 |
| Disney Infinity 3.0 | `PN_DSI3_` | 29200 / 29201 |

`GAMES=disney-infinity-1.0` runs a single game. Ports can be overridden with
`PN_DSIx_AUTH_PORT` / `PN_DSIx_SECURE_PORT`; startup fails if two games share
one. Use one database per game (`PN_DSIx_POSTGRES_URI`).

Two auth modes per game, chosen with `PN_DSIx_LOCAL_MODE`:

- `1`: isolated preservation mode. Passwords come from `settings.json`
  (`username` is the PID) and tokens are not validated. Never expose it publicly.
- otherwise: tokens are decrypted with `PN_DSIx_NEX_TOKEN_AES_KEY` and
  passwords derived with `PN_DSIx_NEX_PASSWORD_SECRET`.

`PN_DSIx_SECURE_HOST` must be short (about 15 characters) because the retail
client truncates it.

## Test

```powershell
go vet ./...
go test ./...
```

## Notes

Disney Infinity 1.0 is configured as NEX 2.0.0 in `version/dsi1.go`, which
selects the legacy connection signature and no structure headers. If clients
loop on CONNECT, that version is the first thing to verify.
