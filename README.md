## vault-backend-migrator

`vault-backend-migrator` is a tool to export and import (migrate) data across vault clusters.

Right now this tool really only supports the `generic` backend. Other mount points often create dynamic secrets that likely should be recreated on the new cluster.

### Usage

##### Setup

You'll want to `go get -f -u github.com/adamdecaf/vault-backend-migrator` this project. It will pull down the vault source code to your GOPATH too.

##### Exporting

After pulling the code it's helpful to set a few environment variables. (These may not all be required for your setup.) These variables will be for the vault you're exporting from.

```
export VAULT_ADDR=http://127.0.0.1:8200/
export VAULT_CACERT=<full filepath to .crt bundle>
export VAULT_TOKEN=<vault token>
```

Note: You'll need to make sure the VAULT_TOKEN has permissions to list and read all vault paths.

Then you should be able to run an export command:

```
$ ./vault-backend-migrator -export secret/ -file secrets.json
```

This will create a file called `secrets.json` that has all the keys and paths. (Note: This is literally all the secrets from the generic backend. Don't share this file with anyone! The secret data is **encoded** in base64, but there's no protection over this file.)

##### Importing

Once you've created an export you are able to reconfigure the vault environment variables (`VAULT_ADDR` and `VAULT_TOKEN` usually) to run an import command.

```
$ ./vault-backend-migrator -import secret/ -file secrets.json
```

This will output each key the tool is writing to. After that a simple `vault list` command off the vault cli will show the secrets there.

Note: It's recommended that you now delete `secrets.json` if you don't need it. If you can install a tool like `srm` to really delete this file.
- OSX: `brew install srm`
- Ubuntu: `apt-get install secure-delete`

### Configuration

This tool reads all the `VAULT_*` environment variables as the vault cli does. You likely need to specify those for the address, CA certs, etc.
