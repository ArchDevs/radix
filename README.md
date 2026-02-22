# Radix

Self-hosted messenger with cryptographic identity.

Fast. Private. Yours.

---

## Overview

Radix is a real-time messaging platform combining the usability of modern chat apps with cryptographic security and self-hosted sovereignty.

- **No phone numbers** — Ed25519 keypairs, 12-word recovery
- **No central control** — Self-host or use hosted (in the future), your choice
- **Real-time** — WebSocket delivery with offline persistence

---

## Architecture

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Desktop   │───▶│   Server    │◀───│   Mobile    │
│  (Wails 3)  │     │    (Go)     │     │  (Future)   │
└─────────────┘     └──────┬──────┘     └─────────────┘
                           │
                    ┌──────┴──────┐
                    │   SQLite    │
                    │  (embedded) │
                    │             │
                    └─────────────┘
```

---

## Components

| Layer      | Technology               | Purpose                    |
|------------|--------------------------|----------------------------|
| Identity   | Ed25519, BIP39           | Key generation, recovery   |
| Auth       | Challenge-response, JWT  | Passwordless login         |
| Transport  | WebSocket, HTTP/1.1      | Real-time + fallback       |
| Encryption | NaCl box, Double Ratchet | E2EE messaging             |
| Storage    | SQLite                   | Local + server persistence |

---

## Quick Start

```bash
# Clone
git clone https://github.com/ArchDevs/radix
cd radix

# Server
make build # Build new API binary
make run # Build and run the API
```

---

## Development

### Requirements

- Go 1.25+

### Structure

```
.
├── cmd/
│   ├── server/         # Headless server
│   └── desktop/        # Wails application
├── internal/
│   ├── auth/           # JWT, challenge-response
│   ├── user/           # User related
│   ├── challenge/      # Challenge/Identity logic
│   ├── wsocket/        # WebSocket server
│   └── crypto/         # E2EE implementation
└── migrations/         # Database migrations
```

### Database

```bash
# Run migrations
make migrate
```

---

## Protocol

Radix uses a simple JSON protocol over WebSocket with HTTP fallback.

| Action    | Endpoint                        | Description                 |
|-----------|---------------------------------|-----------------------------|
| Connect   | `WS /ws?token=<jwt>`            | Authenticated WebSocket     |
| Challenge | `GET /challenge?address=<addr>` | Get nonce to sign           |
| Verify    | `POST /verify`                  | Sign challenge, receive JWT |
| Register  | `POST /register`                | Create new identity         |
| Search    | `GET /search?q=<query>`         | Find users by username      |

### Message Format

```json
{
  "type": "dm",
  "to": "rad:A4f7c28d...",
  "content": "<base64-encrypted-content>"
}
```

---

## Security

| Feature         | Implementation                            |Status |
|-----------------|-------------------------------------------|-------|
| Identity        | Ed25519 keypairs, BIP39 mnemonic recovery |  [x]  |
| Authentication  | Challenge-response, no passwords stored   |  [x]  |
| Transport       | TLS 1.3, WebSocket secure                 |  [ ]  |
| Encryption      | NaCl box (X25519 + XSalsa20 + Poly1305)   |  [ ]  |
| Forward secrecy | Ephemeral keys per session                |  [ ]  |

---

## Roadmap

- [x] Cryptographic identity
- [x] Passwordless authentication
- [x] Real-time messaging
- [x] Offline message persistence
- [ ] Encryption
- [ ] Mobile applications
- [ ] Voice and video calls (WebRTC)
- [ ] Group conversations
- [ ] Federation between servers

---

## License

GPL-3.0