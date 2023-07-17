# Static Service

> Reference https://docs.google.com/document/d/1QV3K0fMbFkGOrqRnfwm5QbXkwtXoQVFfnMDYvhXDzhY/edit?usp=sharing

Independent service that stores the information
about your project on the storage.

There are multiple types of the storage that 
this service could use. The storages are
plugged using `database` extensions.

## W3Storage Extension
**Static Service** along with 
**W3storage Extension** stores the data
on IPFS using web3.storage.

Below are the list of the routes for public access.

***
### /abi/abi_id.json

The route to store the ABI of the smartcontracts.
The *abi_id* is the hash sum of the whole abi file.

Static Service Command to set abi is

```javascript
  "set_abi"
```

***
### /smartcontract/topic_id.json

The route to store the smartcontract parameters.
The *topic_id* in the name is the unique identifier instead
of the smartcontract address.

The *topic_id* is derived from with the `org`, `net` and `name` parameters:

```plain
import "github.com/ahmetson/common-lib/topic"
topic.Topic{}
```

***
### /configuration/organization/project.json

The route to store the list of smartcontracts used by the dapp.
The `organization` is the name of the organization that is accountable.
The `project` is the name of the configurations.

The example of the json body:

```json
[
  "org-seascape.net-id.name-seascape"
]
```

***
### /configuration/organization.json

The route to get the list of all project configurations within the organization.

```json
[
  "dev"
]
```

***
### /smartcontract/organization.json

The route to get the list of all smartcontracts within the organization.

```json
[
  "org-name.net-1.name-crowns"
]
```

***

## Commands

### Set Abi

Set Abi command stores the information about the smartcontract
in the static server.

**Request Message**

```typescript
// instance of message.Request
let request = {
    command: "set_abi",
    parameters: {
        bytes: ""  
    }
}
```

The `request.parameters.bytes` is the serialized ABI of the smartcontract.

**Reply**

```typescript
// instance of message.Reply
let reply = {
    status: "OK",
    parameters: {
        abi_id: "unique_id"
    }
}
```

Upon the successful response, the Abi will be available on CDN:

For example with **AliOss Seascape Extension** the abi will be on
`/abi/unique_id.json`

***
#### Set Smartcontract

Set Smartcontract command stores the information about the smartcontract in the
static server.

**Request Message**

```typescript
let request = {
    command: "set_smartcontract",
    parameters: {
        topic: {
            org: "seascape", // organization
            proj: "main", // part of core business
            net: "1", // 1 for ethereum mainnet
            group: "token", // classification
            name: "Crowns"
        },
        transaction_id: "0xbeefdead",
        owner: "0xdead",
        verifier: "0xdead",
        specific: {
            address: "0xdead",
            abi_id: "unique_id"
        }
    }
}
```

**Reply Message**

```typescript
let reply = {
    status: "OK",
    parameters: {},
    message: ""
}
```

Upon successful update, the following two parameters will be changed:

1. `/smartcontract/org-seascape.net-1.name-crowns.json` &ndash; with the
smartcontract information.
2. `/smartcontract/seascape.json` &ndash; adds the smartcontract.


***
#### Set Configuration

Set Configuration command stores the dapp 
configuration with the list of smartcontracts in this dapp.

**Request Message**

```typescript
let request = {
    command: "set_configuration",
    parameters: {
        // topic.Topic with 'org' and 'proj'
        id: {
            org: "seascape", // organization
            proj: "dev",
        },
        // topic.Topic[] with 'org', 'net' and 'name'
        smartcontracts: [
            {
                org: "seascape",
                net: "1",
                name: "Crowns"
            }  
        ]
    }
}
```

**Reply Message**

```typescript
let reply = {
    status: "OK",
    parameters: {},
    message: ""
}
```

Upon successful update, the following two parameters will be changed:

1. `/configuration/seascape/dev.json` &ndash; with the project
2. `/configuration/seascape.json` &ndash; with the new project in the list.