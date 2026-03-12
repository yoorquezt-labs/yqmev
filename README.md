<p align="center">
  <strong>quezt</strong><br>
  <em>MEV terminal dashboard by YoorQuezt</em>
</p>

<p align="center">
  <a href="https://github.com/yoorquezt-labs/quezt/releases"><img src="https://img.shields.io/github/v/release/yoorquezt-labs/quezt?style=flat-square" alt="Release"></a>
  <a href="https://github.com/yoorquezt-labs/quezt/actions"><img src="https://img.shields.io/github/actions/workflow/status/yoorquezt-labs/quezt/release.yaml?style=flat-square" alt="CI"></a>
  <a href="https://www.npmjs.com/package/quezt"><img src="https://img.shields.io/npm/v/quezt?style=flat-square" alt="npm"></a>
  <a href="https://github.com/yoorquezt-labs/quezt/blob/master/LICENSE"><img src="https://img.shields.io/badge/license-MIT-blue?style=flat-square" alt="License"></a>
</p>

---

A terminal UI for monitoring, analyzing, and interacting with [YoorQuezt](https://yoorquezt.com) MEV infrastructure. Think [k9s](https://k9s.io) for MEV.

```
quezt --gateway ws://your-gateway:9099/ws
```

## Install

```bash
# npm
npm install -g quezt

# homebrew
brew install yoorquezt-labs/tap/quezt

# curl
curl -fsSL https://quezt.dev/install | sh

# go
go install github.com/yoorquezt-labs/quezt/cmd/quezt@latest

# from source
git clone https://github.com/yoorquezt-labs/quezt.git
cd quezt && make install
```

## Tabs

| # | Tab | Description |
|---|-----|-------------|
| 1 | **Dashboard** | Live overview — health, mempool, blocks, bundles, relays |
| 2 | **Bundles** | Browse and inspect MEV bundles in the auction pool |
| 3 | **Auction** | Real-time sealed-bid auction feed via WebSocket |
| 4 | **Blocks** | Built blocks with profit breakdown |
| 5 | **Relay** | Relay marketplace — registration, reputation, stats |
| 6 | **OFA Protect** | Order Flow Auction — protected txs, sandwich detection, rebates |
| 7 | **Intents** | Intent submission and solver resolution |
| 8 | **Orderflow** | Orderflow payment tracking and analytics |
| 9 | **System** | Simulation cache, auth stats, gateway diagnostics |
| 10 | **Q AI** | AI chat assistant with live MEV engine context |
| 11 | **Logs** | Real-time WebSocket event stream |
| 12 | **Help** | Keyboard shortcuts and navigation guide |
| 13 | **Settings** | Gateway URL, API key, AI provider configuration |

## Usage

```bash
# Connect to gateway
quezt --gateway ws://localhost:9099/ws

# With API key
quezt --gateway ws://mev.example.com:9099/ws --api-key YOUR_KEY

# AI chat with Claude
quezt --ai-provider claude --ai-key sk-ant-...

# AI chat with OpenAI
quezt --ai-provider openai --ai-key sk-...

# AI chat with local Ollama
quezt --ai-provider ollama --ai-model llama3.1
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Next / previous tab |
| `1`-`9`, `0` | Jump to tab |
| `r` | Refresh current view |
| `q` / `Ctrl+C` | Quit |
| `/` | Filter (in list views) |
| `Enter` | Select / submit |
| `Esc` | Back / cancel |

## Q AI Assistant

The Q AI tab provides an AI-powered MEV assistant with live engine context. It knows your current mempool size, bundle count, relay status, OFA stats, and recent events — no copy-pasting needed.

Supported providers:
- **Claude** (Anthropic) — default, best MEV analysis
- **OpenAI** (GPT-4o)
- **Ollama** — fully local, no API key needed

Example questions:
```
> What's the current state of the auction pool?
> Why did my last bundle revert?
> How can I improve my bid strategy?
> Show me the most profitable arb opportunities right now
```

## SDK

The `pkg/` packages can be imported directly in your Go projects:

```go
import (
    "github.com/yoorquezt-labs/quezt/pkg/client"
    "github.com/yoorquezt-labs/quezt/pkg/jsonrpc"
    "github.com/yoorquezt-labs/quezt/pkg/types"
)

// Connect to the MEV gateway
c, err := client.Dial(client.Config{
    GatewayURL: "ws://localhost:9099/ws",
    APIKey:     "your-key",
})
defer c.Close()

// Submit a bundle
result, err := c.SendBundle(ctx, types.BundleMessage{
    BundleID:     "bundle-001",
    Transactions: txs,
    BidWei:       "1000000000000000",
})

// Subscribe to auction events
subID, events, err := c.Subscribe(ctx, jsonrpc.TopicAuction)
for event := range events {
    fmt.Println("Auction event:", string(event))
}
```

## Architecture

```
quezt (public binary)
  |
  |  WebSocket JSON-RPC 2.0
  v
yqmev gateway (private)
  |
  v
MEV engine, auction, relay, OFA, intents
```

The `quezt` binary is a pure client — it connects to the YoorQuezt MEV gateway over WebSocket and renders data in the terminal. No private MEV engine code is included.

## Development

```bash
make build      # Build binary to bin/quezt
make install    # Install to $GOPATH/bin
make test       # Run tests
make snapshot   # GoReleaser snapshot (test release locally)
make release    # GoReleaser full release (requires tag)
```

## License

MIT
