# deploy — docs

The togo **deploy subsystem**: a provider-agnostic `Deployer` contract. Targets ship as driver plugins (`deploy-docker`, `deploy-kubernetes`, `deploy-terraform`, `deploy-aws`, …), selected with `DEPLOY_PROVIDER` or `deploy.provider` in `togo.yaml`.

```yaml
deploy:
  provider: docker        # docker | kubernetes | terraform | aws | gcp | … | ubuntu | centos | debian
  host: 1.2.3.4
  domain: app.example.com
```
