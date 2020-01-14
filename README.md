# vault-mv

Simple tool to copy/move secrets in Vault.

## Install

```
$ git clone https://github.com/sgmac/vault-mv
$ go build 
```

## Usage

Create a token and export  *VAULT_TOKEN*

`$ vault token create -policy=YOUR_POLICY`

Copying secrets from one environment to a new one.
```
$ vault-mv  secret/stg secret/ent
copying secret/stg/app/provider/aws to secret/ent/app/provider/aws
copying secret/stg/app/provider/gcp to secret/ent/app/provider/gcp
copying secret/stg/password to secret/ent/password
copying secret/stg/redis_url to secret/ent/redis_url
copying secret/stg/username to secret/ent/username
```

