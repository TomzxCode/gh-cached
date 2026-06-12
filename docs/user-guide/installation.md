# Installation

## Binary download

Download the latest pre-built binary from the [releases page](https://github.com/TomzxCode/gh-cached/releases/tag/latest).

=== "darwin-amd64"
    ```bash
    curl -LO https://github.com/TomzxCode/gh-cached/releases/latest/download/gh-cached-darwin-amd64
    chmod +x gh-cached-darwin-amd64
    mv gh-cached-darwin-amd64 /usr/local/bin/gh-cached
    ```

=== "darwin-arm64"
    ```bash
    curl -LO https://github.com/TomzxCode/gh-cached/releases/latest/download/gh-cached-darwin-arm64
    chmod +x gh-cached-darwin-arm64
    mv gh-cached-darwin-arm64 /usr/local/bin/gh-cached
    ```

=== "linux-amd64"
    ```bash
    curl -LO https://github.com/TomzxCode/gh-cached/releases/latest/download/gh-cached-linux-amd64
    chmod +x gh-cached-linux-amd64
    mv gh-cached-linux-amd64 /usr/local/bin/gh-cached
    ```

=== "windows-amd64"
    ```powershell
    curl -LO https://github.com/TomzxCode/gh-cached/releases/latest/download/gh-cached-windows-amd64.exe
    rename-item gh-cached-windows-amd64.exe gh-cached.exe
    ```

## Build from source

Requires [Go](https://go.dev/) 1.21 or later.

```bash
go install github.com/tomzxcode/gh-cached@main
```

This places the binary at `$(go env GOPATH)/bin/gh-cached`. Make sure that directory is on your `PATH`.

## Verify

```bash
gh-cached --version
```
