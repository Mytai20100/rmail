# rmail

![version](https://img.shields.io/badge/version-0.1-black) ![go](https://img.shields.io/badge/go-1.25+-00ADD8) ![cloudflare](https://img.shields.io/badge/cloudflare-email_routing-F38020)

A minimal local web UI for managing Cloudflare Email Routing rules. Create, edit, and delete custom email addresses across your zones without touching the Cloudflare dashboard.

## How it works

rmail runs a small HTTP server on `localhost:7432`. It talks directly to the Cloudflare API using a token you provide — no data leaves your machine except to Cloudflare. The token is saved locally to `config.yml` and reused on next launch.

Each email address is a routing rule: it matches an incoming address (e.g. `hello@yourdomain.com`) and either forwards it to a destination address, sends it to a Worker, or drops it.

## Setup

The token needs the following permissions: **Zone / Zone / Read** and **Zone / Email Routing / Edit**.

```
go mod tidy
go build -o rmail .
./rmail
```

A browser window opens automatically at `http://localhost:7432`.
