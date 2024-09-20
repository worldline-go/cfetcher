# Config Fetcher

This repo help us to fetch all configs from vault and consul and put them in a nested directory.  
With using this one we can mock all configs or change and copy them to another environment.

Our stack use vault for getting secrets and some generic configs. Before to vault we fetch the consul to load configs.

> In future we will have one config server to manage that in one place.

## Usage

Run the following command to fetch all configs from vault and consul:

```sh
# load credentials
source env/local.sh

# fetch all configs from vault
make run-vault
# fetch all configs from consul
make run-consul
```

### With Binary

Create CLI binary:

```sh
make build-linux build-windows
```

## Mocking

Use [turna](https://rakunlabs.github.io/turna/) tool to mock the loaded configuration for vault and consul.

This mock file designed for vault and consul configs run on same turna server.

```sh
docker run --rm -it \
-e LOG_LEVEL=debug \
-p 8080:8080 \
-v $(pwd)/out/test/finops-consul:/finops-consul -v $(pwd)/out/test/finops-vault:/finops-vault \
-v $(pwd)/mock/turna.yaml:/turna.yaml \
ghcr.io/rakunlabs/turna:latest
```

Test mocking configs

```sh
make run-load
```

<details><summary>Development</summary>

## Testing Vault

First create a vault server

```sh
make vault
```

After that login

```sh
export VAULT_ADDR=http://127.0.0.1:8200
TOKEN=$(docker logs vault |& grep "Root Token:" | cut -d":" -f2 | xargs) make vault-login
```

Create approle with a secret values, get role-id and secret in the output

```sh
make vault-role-enable vault-role vault-secret
```

For testing set this env values
```sh
export VAULT_ROLE_ID=xxx
export VAULT_ROLE_SECRET=xxx
```

</details>
