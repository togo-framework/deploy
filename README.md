<!-- togo-header -->
<div align="center">
  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />
  <h1>togo-framework/deploy</h1>
  <p>
    <a href="https://to-go.dev/marketplace"><img src="https://img.shields.io/badge/marketplace-to--go.dev-1FC7DC" alt="marketplace" /></a>
    <a href="https://pkg.go.dev/github.com/togo-framework/deploy"><img src="https://pkg.go.dev/badge/github.com/togo-framework/deploy.svg" alt="pkg.go.dev" /></a>
    <img src="https://img.shields.io/badge/license-MIT-blue" alt="MIT" />
  </p>
  <p><strong>Part of the <a href="https://to-go.dev">togo</a> framework.</strong></p>
</div>

## Install

```bash
togo install togo-framework/deploy
```
<!-- /togo-header -->

togo's **deployment subsystem** — a provider-agnostic `Deployer` contract
(`Provision`/`Deploy`/`Destroy`/`Status`). Real targets ship as **driver plugins**
that call `deploy.RegisterDriver`; pick one with `DEPLOY_PROVIDER` or
`deploy.provider` in `togo.yaml`. The CLI's `togo deploy` resolves the provider
and runs it.

```bash
togo install togo-framework/deploy            # the base
togo install togo-framework/deploy-docker     # a target driver
```

Drivers: `deploy-terraform`, `deploy-docker`, `deploy-kubernetes`,
`deploy-aws`, `deploy-gcp`, `deploy-azure`, `deploy-digitalocean`, `deploy-vultr`,
`deploy-hetzner`, `deploy-ovh`, `deploy-ubuntu`, `deploy-centos`, `deploy-debian`.

## Configure (`togo.yaml`)

```yaml
deploy:
  provider: docker
  host: 1.2.3.4
  user: root
  domain: app.example.com
```

`DEPLOY_PROVIDER` overrides `provider`. The default is `log` (no-op, safe for dev).

<!-- togo-sponsors -->
---
<div align="center">
  <h3>Premium sponsors</h3>
  <p><a href="https://id8media.com"><strong>ID8 Media</strong></a> &nbsp;·&nbsp; <a href="https://one-studio.co"><strong>One Studio</strong></a></p>
  <p><sub><a href="https://github.com/sponsors/fadymondy">Become a sponsor</a>.</sub></p>
</div>
<!-- /togo-sponsors -->
