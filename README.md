# HARA DID Root SDK – Monorepo Overview

> **Status:**  
> This project is **currently untested** and will likely undergo **significant changes** to align with future business requirements and updates in the smart contracts.

This Go monorepo contains two primary SDKs:

- **`core-general-sdk/`** – HARA Core Blockchain Library (general blockchain utilities)
- **`did-root-sdk/`** – DID Root SDK (RootFactory + RootRegistry/RootStorage clients)

This README focuses only on **folder structure**, **module imports**, and the **repository references** required for development.

---

# Repository List (Clone Guide)

For full development across the DID Root stack, you need these repositories:

| Component | Repository URL | Description |
|----------|----------------|-------------|
| **Smart Contracts** | https://dev.azure.com/dattabot/Blockchain%202024/_git/sc-hara | RootFactory, RootRegistry, RootStorage, interfaces, structs. |
| **Core General SDK** | https://dev.azure.com/dattabot/Blockchain%202024/_git/core-general-sdk | Network helpers, Blockchain client, ABI loader, utils. |
| **DID Root SDK (this repository)** | https://github.com/HARA-DID/did-root-sdk | RootFactory (write), RootStorage (read), HNS loading. |

---

# Folder Structure

## Root

### **`core-general-sdk/`**
- Go module: `github.com/HARA-DID/hara-core-blockchain-lib`
- Main packages:
  - `pkg/` – public blockchain utilities:
    - `network.go` – RPC network wrapper (multi-endpoint, logging, health checks)
    - `blockchain.go` – wallet, ABI resolution, contract helpers
  - `utils/` – shared helper functions and types
    - `types.go` – `TransactionParams`, `RPCRequest`, `ContractDetail`, etc.
    - `functions.go` – `Namehash`, `EncodeArgs`, `DefaultLogConfig`, etc.
  - `internal/` – ABI cache, RPC implementation, logger internals

---

### **`did-root-sdk/`**
- Go module: `github.com/HARA-DID/did-root-sdk`
- Depends on `core-general-sdk`

#### Packages:

#### **`pkg/rootfactory/`**
- High-level write SDK for **RootFactory**
- Provides typed methods for:
  - `CreateDID`, `UpdateDID`, `DeactivateDID`, `ReactivateDID`
  - `TransferDIDOwner`, `StoreData`, `DeleteData`
  - `AddKey`, `RemoveKey`, `AddClaim`, `RemoveClaim`
- Uses:
  - `Network.ArgBuilder`
  - `utils.EncodeArgs`
- Types match on-chain `IRootStructParams`

#### **`pkg/rootstorage/`**
- Read-only client for RootRegistry, RootList, RootStorage
- Provides:
  - DID resolution (`ResolveDID`, `VerifyDIDOwnership`)
  - Key inspection (`GetKey`, `GetKeysByDID`)
  - Claim inspection (`GetClaim`, `GetClaimsByDID`, `VerifyClaim`)
  - Generic view methods (`GetData`, `GetOriginalKey`, etc.)
- Types match on-chain `DIDDocument`, `Key`, `Claim`

#### **Examples**
- `example_rootfactory.go` – full write flow
- `example_rootstorage.go` – read-only flow

---

# Go Modules and Local Development

Each SDK is a separate module:

- `core-general-sdk/go.mod`  
  `module github.com/HARA-DID/hara-core-blockchain-lib`

- `did-root-sdk/go.mod`  
  `module github.com/HARA-DID/did-root-sdk`

During local development, `did-root-sdk` uses a **replace** directive to point to the local path:

```go
require github.com/HARA-DID/hara-core-blockchain-lib v0.0.0

replace github.com/HARA-DID/hara-core-blockchain-lib => ./core-general-sdk
```

When publishing the SDK, remove `replace` and use semantic versioning.

---

# Importing the Libraries

## Importing `core-general-sdk`

```go
import (
    "github.com/HARA-DID/hara-core-blockchain-lib/pkg"
    "github.com/HARA-DID/hara-core-blockchain-lib/utils"
)
```

Basic usage:

```go
urls := []string{"https://rpc-1", "https://rpc-2"}

net := pkg.NewNetwork(urls, "2.0", 1, utils.DefaultLogConfig())
bc := pkg.NewBlockchain("your-seed-phrase", net, big.NewInt(1212))
```

---

## Importing `did-root-sdk` (RootFactory – write)

```go
import (
    "github.com/HARA-DID/did-root-sdk/pkg/rootfactory"
    "github.com/HARA-DID/hara-core-blockchain-lib/pkg"
    "github.com/HARA-DID/hara-core-blockchain-lib/utils"
)
```

Initialization:

```go
net := pkg.NewNetwork([]string{"http://rpc:5625"}, "2.0", 1, utils.DefaultLogConfig())
bc := pkg.NewBlockchain("mnemonic", net, big.NewInt(1212))

ctx := context.Background()

rf, err := rootfactory.NewRootFactoryWithHNS(
    ctx,
    "rootfactorydevharadid.hara.ethnet",
    bc,
)
```

Example call:

```go
keyID, _ := rootfactory.GenerateKeyIdentifier()
hashes, err := rf.CreateDID(ctx, "", keyID, false)
```

---

## Importing `did-root-sdk` (RootStorage – read-only)

```go
import (
    rootstorage "github.com/HARA-DID/did-root-sdk/pkg/rootstorage"
    "github.com/HARA-DID/hara-core-blockchain-lib/pkg"
    "github.com/HARA-DID/hara-core-blockchain-lib/utils"
)
```

Initialization:

```go
net := pkg.NewNetwork([]string{"http://rpc:5625"}, "2.0", 1, utils.DefaultLogConfig())
bc := pkg.NewBlockchain("mnemonic", net, big.NewInt(1212))

rs, err := rootstorage.NewClientWithHNS(
    ctx,
    "rootregistrydevharadid.hara.ethnet",
    bc,
)
```

Read examples:

```go
doc, _ := rs.ResolveDID(ctx, didHash)
keys, _ := rs.GetKeysByDID(ctx, didHash)
claims, _ := rs.GetClaimsByDID(ctx, didHash)
data, _ := rs.GetData(ctx, keyCode)
```

---

# Summary

- `core-general-sdk/` → RPC, network manager, blockchain utilities.  
- `did-root-sdk/pkg/rootfactory/` → All **write** operations to RootFactory.  
- `did-root-sdk/pkg/rootstorage/` → All **read-only** operations to RootRegistry/RootStorage.  
- Repo URLs included for full-stack development.  
- Project is currently **untested** and will likely change significantly.

