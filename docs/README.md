# deploy — docs

**Deploy.** Provider-agnostic deploy subsystem — provision, deploy, destroy and status over a pluggable driver.

## Install

```bash
togo install togo-framework/deploy
```

Select a driver via **deploy.provider in togo.yaml (or DEPLOY_PROVIDER)**. Drive it from the CLI with **`togo deploy`**.

## Interface

`Deployer` — `Provision`/`Deploy`/`Destroy`/`Status` over a `Spec{App,Dir,BuildCmd,Host,User,Image,Region,Domain}` built from your `togo.yaml`.

## Configuration

| Env var | Description |
|---|---|
| `DEPLOY_PROVIDER` | Selects the deploy driver (alternative to togo.yaml `deploy.provider`). |

## Usage & notes

The base wires `togo deploy` to a driver. Each provider self-registers via `init()`; the CLI resolves it with `deploy.Build(name, k)` (no app boot needed). Set the provider + per-provider settings in `togo.yaml`:
```yaml
deploy:
  provider: docker   # docker|kubernetes|terraform|ubuntu|centos|debian|aws|gcp|azure|digitalocean|vultr|hetzner|ovh
  image: ghcr.io/me/app:latest
  host: 1.2.3.4        # VPS drivers
  user: root
  region: eu-central
  domain: app.example.com
```

## Example

```bash
togo deploy --provider docker --dry-run   # preview the plan
```

## Links

- [Driver plugins](https://to-go.dev/marketplace)
- [Marketplace](https://to-go.dev/marketplace)
- [Source](https://github.com/togo-framework/deploy)
