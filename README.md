# go‑iscsiadm

> A lightweight Go wrapper around the system `iscsiadm` command, designed for use in Container Storage Interface (CSI) drivers.

---
## Table of Contents
1. [Overview](#overview)
2. [Prerequisites & Installation](#prerequisites--installation)
3. [Quick Start](#quick-start)
4. [API Reference](#api-reference)
5. [Error Handling](#error-handling)
6. [FAQ / Troubleshooting](#faq---troubleshooting)
7. [Contributing](#contributing)
8. [License](#license)

## Overview
`go-iscsiadm` exposes a small, well‑typed API that internally runs the `iscsiadm` binary and parses its output.

* **Login / Logout** – Attach or detach an iSCSI target.
* **Discover** – Query a portal for available targets.
* **Remove** – Delete a discovered entry from the local database.
* **Sessions** – List active sessions on the host.

The package is intentionally minimal; it does not implement any of the CSI spec itself but can be used as a building block in a driver implementation.

## Prerequisites & Installation
```bash
# Go 1.20+ required (module uses generics and context)
go get github.com/OchiengEd/go-iscsiadm@v0.1.0
```
* The host must have the `iscsiadm` binary in `$PATH`. It is part of most Linux distributions’ iSCSI utilities.
* If you need to run commands inside a container, set the `useNsenter` option when creating the controller (see API reference).

## Quick Start
```go
package main

import (
    "context"
    "fmt"

    iscsiadm "github.com/OchiengEd/go-iscsiadm"
)

func main() {
    ctrl := iscsiadm.New(iscsiadm.WithNsenter(true)) // optional, only needed in containers

    ctx := context.Background()

    // Discover targets on a portal.
    tgtList, err := ctrl.Discover(ctx, &iscsiadm.DiscoverRequest{Portal: "192.168.100.5:3260"})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Found %d target(s)\n", len(tgtList))

    // Login to the first one.
    if len(tgtList) > 0 {
        dev, err := ctrl.Login(ctx, &iscsiadm.LoginRequest{Portal: tgtList[0].portal, TargetIQN: tgtList[0].name})
        if err != nil { panic(err) }
        fmt.Printf("Device mounted at %s\n", dev)
    }
}
```

## API Reference
| Exported Type / Function | Description |
|--------------------------|-------------|
| `SystemController` | Main entry point. Holds a `command.Runner`. Use `New()` to create it.
| `WithCommandRunner(r command.Runner)` | Option for injecting a custom runner (e.g., mock in tests).
| `WithNsenter(b bool)` | Enable nsenter prefix when running commands inside containers.
| `Login(ctx, req *LoginRequest) (Device, error)` | Log into an iSCSI target. Returns the device path (`/dev/sdX`).
| `Logout(ctx, req *LogoutRequest) (bool, error)` | Terminate a session; returns success flag.
| `Discover(ctx, req *DiscoverRequest) ([]Target, error)` | Discover targets from a portal.
| `Remove(ctx, req *RemoveRequest) (bool, error)` | Delete an entry from the discovery database.
| `Sessions(ctx) ([]Session, error)` | List active sessions on the host.

### Types
```go
type Target struct { name string; portal string }
type Session struct { target, portal, session string }
```
All fields are exported via getters in a full implementation (currently omitted for brevity).

## Error Handling
`command/runner.go` maps `iscsiadm` exit codes to custom errors. The public API returns these as Go errors.
| Exit Code | Custom Error Type | Human‑Readable Message |
|-----------|------------------|------------------------|
| 1 (SessionExists) | `ErrSessionExists` | *"iSCSI session already exists."* |
| 2 (ObjectsNotFound) | `ErrResourceNotFound` | *"Target or portal not found in database."* |
| 3 (LoginFailure) | `ErrLoginFailed` | *"Unable to login; check credentials and network connectivity."* |
| 4 (LogoutFailure) | `ErrLogoutFailed` | *"Unable to logout the session."* |
| 5 (AccessDenied) | `ErrPermissionDenied` | *"Insufficient OS permissions – run as root or with appropriate capabilities."* |

These errors are exported from the package and can be type‑asserted by callers.

## FAQ / Troubleshooting
- **Why does login fail even though I see the target in discovery?**
  - Verify that you’re using the correct portal/target IQN pair. `iscsiadm` may reject duplicate sessions; check for an existing session with `iscsiadm -m node --show`. 
- **I’m running inside a container and commands fail.**
  - Enable nsenter by passing `WithNsenter(true)` when creating the controller.
- **How do I see raw command output?**
  - Use a custom runner that logs stdout/stderr or set up a mock in tests.

## Contributing
Pull requests are welcome! Please run `go test ./...` before submitting. If you add new public functions, update the README accordingly.

## License
MIT © Edmund Ochieng (2026)
