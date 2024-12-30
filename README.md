# go-crawl-on-them

**go-crawl-on-them** is a basic web-crawler written in Go, designed to find dead links inside a website

## Features

- **Finding dead links**: Detects refernces to dead links on a given website and checks every subsite for dead links too.
- **Easy setup**: Easy to setup via **git clone** or **go install**. ([Installation](#installation))
- **User-Experience**: Simple UI, basic Syntax & lightweight.

## Platforms

| Platform       | Developed | Tested |     Version     |
|----------------|:---------:|:------:|:---------------:|
| Windows        |     ✅     |   ✅    | Windows 11 24H2 |
| Linux          |     ✅     |   ✅    | Ubuntu 22.04 on WSL |
| macOS          |     ✅     |   ❌    |        ❌        |

## Installation

### Via `go install`:

1. **Install the package**:
**Linux**:

```bash
go install github.com/hustender/go-crawl-on-them@latest
mv $(go env GOPATH)/bin/go-crawl-on-them $(go env GOPATH)/bin/crawl
```
**Windows**:

PowerShell:
```powershell
go install github.com/hustender/go-crawl-on-them@latest
Rename-Item "$(go env GOPATH)\bin\go-crawl-on-them.exe" "crawl.exe"
```
Command Prompt:
```cmd
go install github.com/hustender/go-crawl-on-them@latest
rename %GOPATH%\bin\go-crawl-on-them.exe crawl.exe
```

2. **Run the program**:
```bash 
crawl <url>
```

### Via `git clone`:

1. **Clone the repository**
```bash
git clone https://github.com/hustender/go-crawl-on-them.git
cd go-crawl-on-them/
```

2. **Build the program**:
**Linux**:

```bash
go build -o $(go env GOPATH)/bin/crawl
```

**Windows**:

PowerShell:
```powershell
go build -o "$(go env GOPATH)\bin\crawl.exe"
```
Command Prompt:
```cmd
go build -o %GOPATH%\bin\crawl.exe
```

3. **Run the program**:
```bash
crawl <url>
```

## Example

![go-crawl-on-them](https://github.com/user-attachments/assets/a44f44e0-160f-4abb-83f7-181088de62cf)

## Contributing

To contribute, please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
