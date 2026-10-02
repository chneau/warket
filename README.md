# warket

A small CLI for the [warframe.market](https://warframe.market) API (v2).

## Requirements

- Go 1.26+

## Install

```bash
go install -v github.com/chneau/warket@latest
```

If on Windows, be sure to have `%userprofile%\go\bin` in your PATH.

Or install the toolchain with Chocolatey (sort of a package manager for Windows):

1. Install Chocolatey: `https://chocolatey.org/install`
2. Install Go: `choco install golang`
3. Install this repo: `go install -v github.com/chneau/warket@latest`

## Usage

```bash
warket snipe --gain 5        # watch newly posted orders and flag good deals
warket watch <user-slug>     # live table of a market user's buy/sell orders
```

Both commands need a terminal. On Windows run `winpty warket ...` if the
console output looks wrong.

### snipe

Subscribes to the realtime feed of newly posted orders and reports ones that
undercut the current best online offer, when the gap to the next best offer is
at least the gain threshold.

| Flag | Default | Description |
| --- | --- | --- |
| `--gain`, `-g` | `5` | minimum gain in platinum |
| `--notification`, `-n` | `true` | desktop notification |
| `--copy`, `-c` | `true` | copy the whisper message to the clipboard |

### watch

Polls a user's orders and renders them in a table with rank, price diff and the
best competing offers.

> The username is the user's warframe.market **slug** (the one in the profile
> URL, e.g. `gh0stman`), not their in-game name. v2 only resolves users by slug.

| Flag | Default | Description |
| --- | --- | --- |
| `--buy`, `-b` | `true` | show buy orders |
| `--sell`, `-s` | `true` | show sell orders |
| `--log`, `-l` | `10` | number of log lines to keep |
| `--ticker`, `-t` | `-1` | refresh interval in seconds (`-1` renders once) |
| `--sort` | `name` | sort column: `name`, `qtt`, `place`, `price`/`plat`, `diff` |

## API

This client talks to the v2 API at `https://api.warframe.market/v2` and uses
the `wfm` WebSocket subprotocol at `wss://ws.warframe.market/socket`.

- Docs: https://docs.warframe.market/
- API docs: https://warframe.market/api_docs
- Other clients: https://github.com/LastExceed/WarframeMarKT

Useful order endpoints (authentication is required only for writes):

```bash
GET   /v2/orders/item/{slug}   # visible orders for an item
GET   /v2/orders/user/{slug}   # a user's orders
POST  /v2/order                # create an order
# {"itemId":"54aae292e7798909064f1575","type":"sell","platinum":210,"quantity":1,"visible":false,"rank":1}
PATCH /v2/order/{id}           # update an order
# {"platinum":211,"quantity":1,"visible":false}
```

Responses use the v2 envelope `{"apiVersion","data","error"}`. Orders carry an
`itemId` rather than an embedded item, so the client loads `/v2/items` once and
re-attaches item details locally.

## dev

```bash
# use nodemon for hotreload
nodemon -e go -x "go run . s -g 0 || false"

# live integration test (hits the real API)
go test ./...
```
