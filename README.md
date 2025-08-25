<div align='center'>
<h1>kirin</h1>
<span>Scaffold a full-stack gRPC application with end-to-end type safety.</span>
</div>

---

<img src="./kirin.jpg" width="155" height="225">

> **Kirin (Qilin)** is a legendary creature from East Asian mythology, often described as a chimeric beast with the body of a deer, scales of a dragon, and sometimes the tail of an ox. Like how the mythical Kirin is composed of parts from multiple magnificent creatures, kirin-generated applications are composed of multiple powerful technologies - gRPC-Gateway for unified gRPC and REST APIs, gRPC-Web for frontend connectivity, and embedded assets for self-contained deployment.
>
> *Image and creature information credited to [Yokai.com](https://yokai.com/kirin/)*

## 🚀 What is kirin?

kirin is a scaffolding tool that creates full-stack gRPC applications with **frontend and backend coexisting in the same folder structure**. Since **Go 1.18**, embedding files directly into binary executables is natively supported, and kirin leverages this feature to create **self-contained applications** that include both frontend assets and backend logic in a single executable.

This approach enables seamless type generation, easier development workflow, and eliminates the complexity of managing separate repositories for frontend and backend.

## 📦 Installation

### Go Install (Recommended)
```bash
go install github.com/thetnaingtn/kirin@latest
```

### npm
```bash
npm install -g kirin
# or
npx kirin
```

### Homebrew (macOS/Linux)
```bash
brew tap thetnaingtn/kirin
brew install kirin
```

### Download Binary
Download the latest binary from [GitHub Releases](https://github.com/thetnaingtn/kirin/releases)

## 🎮 Usage

### Interactive Mode (Recommended)
Launch the interactive prompt to configure your project step by step:

```bash
kirin
```

**Just follow the prompts!** The interactive mode will guide you through:
1. **App Name**: Choose your application name
2. **Module Name**: Set your Go module name (e.g., `github.com/user/myapp`)
3. **Frontend Framework**: Select from React, Vue, or Svelte

No need to remember complex commands or flags - simply answer each question and kirin will handle the rest.

### Command Line Mode
Create a project directly with command line arguments:

```bash
# Basic usage with default React frontend
kirin create myapp

# With custom module name
kirin create myapp github.com/myuser/myapp

# With specific frontend framework
kirin create myapp --frontend vue
kirin create myapp --frontend svelte
kirin create myapp --frontend react
```

### System Health Check
Verify your system has all required dependencies:

```bash
kirin doctor
```

## ⚙️ Available Commands & Configuration

### Commands Overview
| Command | Description | Aliases |
|---------|-------------|---------|
| `kirin` | Launch interactive prompt (default) | - |
| `kirin create` | Create project via command line | `c` |
| `kirin dev` | Start development with live reload (uses Air) | - |
| `kirin generate` | Generate code from protobuf definitions | `gen`, `g` |
| `kirin doctor` | Check system requirements | `d`, `doc` |
| `kirin --help` | Show help information | `-h` |

### Create Command Options
```bash
kirin create <app-name> [module-name] [flags]

Arguments:
  app-name      Name of the application to create
  module-name   Go module name (optional, defaults to app-name)

Flags:
  -f, --frontend string   Frontend framework (default "react")
                         Supported: react, vue, svelte
  -h, --help             Help for create command
```

### Frontend Frameworks
| Framework | Description | Status |
|-----------|-------------|--------|
| **React** | Popular JavaScript library for building UIs | ✅ Supported |
| **Vue** | Progressive JavaScript framework | ✅ Supported |
| **Svelte** | Cybernetically enhanced web apps | ✅ Supported |

### Doctor Command
The `doctor` command checks your system for:
- **Git**: Required for cloning template repositories
- **Go**: Required for building the application
- **Node.js**: Required for frontend development
- **npm, yarn or pnpm**: Required for frontend package management

## 🛠️ Development Workflow

### Live Reload Development
Kirin provides seamless live reload for both frontend and backend development:

```bash
# Navigate to your generated project
cd myapp

# Start development with live reload
kirin dev
```

**What `kirin dev` does:**
- **🔄 Hot Reload**: Automatically restarts the server when Go files change
- **⚡ Fast Builds**: Uses [Air](https://github.com/cosmtrek/air) under the hood for efficient rebuilding
- **🎯 Frontend Integration**: Serves frontend assets with automatic refresh
- **📝 Live Logging**: Real-time logs and error reporting

### Code Generation
Generate TypeScript types and gRPC code from protobuf definitions:

```bash
# Generate code from protobuf files
kirin generate
# Short form
kirin gen
```

**What `kirin generate` does:**
- **🔍 Validation**: Checks for buf.yaml in proto directory first
- **📁 Directory Change**: Switches to proto directory for buf operations
- **🔄 Build Check**: Runs `buf build` to validate protobuf files
- **🚫 Error Prevention**: Stops generation if validation fails
- **⚡ Type Generation**: Creates TypeScript types for frontend
- **🔄 gRPC Code**: Generates Go gRPC server and client code
- **📁 Directory Restore**: Returns to original directory after completion

### Other Development Commands

1. **Run Tests**:
   ```bash
   go test ./...
   ```

2. **Check System Health**:
   ```bash
   kirin doctor
   ```

### Air Configuration
Kirin automatically configures Air for optimal development experience. The live reload includes:
- Go source files (`.go`)
- Template files
- Configuration files
- Static assets

No manual Air setup required - just run `kirin dev` and start coding!

## 🧩 Composite Architecture

Like the mythical Kirin composed of multiple creature parts, generated applications include:
- **gRPC-Gateway**: Unified gRPC and REST API endpoints
- **gRPC-Web**: Frontend connectivity to gRPC services
- **Protocol Buffers**: Type-safe communication contracts
- **Embedded Assets**: Self-contained frontend resources using Go 1.18+ embed

## 📋 Examples

### Creating a React Application
```bash
# Interactive mode
kirin

# Command line mode
kirin create my-dashboard
kirin create my-dashboard github.com/company/my-dashboard
kirin create my-dashboard --frontend react
```

### Creating a Vue Application
```bash
kirin create admin-panel --frontend vue
```

### Creating a Svelte Application
```bash
kirin create analytics-app github.com/company/analytics --frontend svelte
```

### System Check
```bash
kirin doctor
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.