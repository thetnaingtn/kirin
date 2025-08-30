<div align='center'>
<h1>kirin</h1>
<span>Scaffold a Full-stack Go gRPC Application With End-to-end Type Safety.</span>
</div>

---
## What is kirin?

**[kirin](https://yokai.com/kirin/)** is a scaffolding tool that creates full-stack gRPC applications with **frontend and backend coexisting in the same folder structure**. Since **Go 1.18**, embedding files directly into binary executables is natively supported, and kirin leverages this feature to create **self-contained applications** that include both frontend assets and backend logic in a single executable.

This approach enables seamless type generation, an easier development workflow, and eliminates the complexity of managing separate repositories for frontend and backend.

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

### Step 1: Check System Dependencies

Before creating your first project, verify that all necessary dependencies are installed:

```bash
kirin doctor
```

The doctor command checks for these dependencies:

**Required Dependencies:**

* **git**: Required for cloning template repositories
* **protoc**: Protocol buffer compiler for gRPC services
* **protoc-gen-go**: Go plugin for protoc compiler
* **protoc-gen-go-grpc**: Go gRPC plugin for protoc compiler

**Optional Dependencies (Good to Have):**

* **buf**: Modern protobuf tooling for validation and code generation
* **Node.js & npm/yarn/pnpm**: Required for frontend development and TypeScript client consumption

If any dependencies are missing, the doctor command will provide installation links and instructions.

### Step 2: Create Your Project

#### Option A: Interactive Mode (Recommended)

Launch the interactive prompt to configure your project step by step:

```bash
kirin
```

**Just follow the prompts!** The interactive mode will guide you through:

1. **App Name**: Choose your application name
2. **Module Name**: Set your Go module name (e.g., `github.com/user/myapp`)
3. **Frontend Framework**: Select from React, Vue, or Svelte

No need to remember complex commands or flags — simply answer each question and kirin will handle the rest.

#### Option B: Command Line Mode

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

### Step 3: Generate gRPC Code (Server + TypeScript Client)

Once your project is scaffolded, generate code from your `.proto` definitions:

```bash
# Generate Go server stubs and TypeScript client code
kirin generate
```

What this does:

* **Server:** Produces Go service interfaces and gRPC server stubs for your backend.
* **Client:** Produces a **TypeScript client** for your chosen frontend framework.

By default, kirin looks for your proto files under the `proto` directory. If you placed them elsewhere, pass a custom folder:

```bash
kirin generate --proto-folder api/proto
```

### Step 4: Start Development Server

After generating code, start developing with live reload capabilities.

#### Backend Development

Start the backend development server with automatic reloading:

```bash
kirin dev
```

The `dev` command uses [air](https://github.com/air-verse/air) under the hood for live reloading, which means **all air flags are supported** and can be passed directly:

```bash
# Use Air flags for custom configuration
kirin dev --build.cmd "go build -o ./tmp/main ." --build.bin "./tmp/main"
kirin dev --tmp_dir custom_tmp --build.delay 1000
kirin dev --color.build red --color.runner green
```

#### Frontend Development (Vite)

For frontend development, navigate to your frontend folder and start the development server. kirin projects use **Vite** by default for frontend builds, so this is the same as running a standard Vite application:

```bash
# Navigate to frontend directory (default: 'web')
cd web

# Start frontend dev server (depends on your package manager)
npm run dev      # if using npm
yarn dev         # if using yarn  
pnpm dev         # if using pnpm
```

If you used a custom frontend folder name during project creation, replace `web` with your specified folder name:

```bash
# Example with custom frontend folder
cd ui            # if you specified --frontend-folder=ui
npm run dev
```

### Step 5: Build Full Application

When you are ready to package your application:

```bash
kirin build
```

This command will:

1. Build the frontend using **Vite**, producing production-ready assets.
2. Use those build assets when compiling the Go application.
3. Embed the assets directly into the resulting Go binary, producing a fully self-contained executable.

This way, deployment is simplified — you only need to distribute a single binary.

### Configuration Persistence (.kirin.toml)

When you run `kirin generate` or `kirin build` with custom options (such as `--frontend-folder` or `--main-folder`), kirin saves these preferences into a `.kirin.toml` file in your project root. Subsequent runs of these commands will automatically read from this configuration file.

This means you only need to specify custom flags once — after the first run, your preferences are remembered.

## ⚙️ Available Commands & Configuration

### Commands Overview

| Command          | Description                                                             | Aliases    |
| ---------------- | ----------------------------------------------------------------------- | ---------- |
| `kirin`          | Launch interactive prompt (default)                                     | -          |
| `kirin create`   | Create project via command line                                         | `c`        |
| `kirin build`    | Build full-stack application (frontend + backend)                       | `b`        |
| `kirin dev`      | Start development with live reload (uses Air)                           | -          |
| `kirin generate` | Generate code from protobuf definitions (Go server + TypeScript client) | `gen`, `g` |
| `kirin doctor`   | Check system requirements                                               | `d`, `doc` |
| `kirin --help`   | Show help information                                                   | `-h`       |

### Command Usage

#### Create Command

```bash
# Basic usage with default React frontend
kirin create myapp

# With custom module name
kirin create myapp github.com/myuser/myapp

# With specific frontend framework
kirin create myapp --frontend svelte

# Available flags
--frontend, -f    Frontend framework (default "react")
                  Supported: react, vue, svelte
--help, -h        Help for create command
```

#### Build Command

```bash
# Build full-stack application
kirin build

# With custom directories
kirin build --frontend-folder ui --main-folder app

# With specific package manager
kirin build --pkg-manager pnpm

# With custom output name
kirin build --output myapp

# Available flags
--frontend-folder     Frontend directory name (default "web")
--main-folder         Main directory name (default "cmd")
--pkg-manager         Package manager to use (npm, yarn, pnpm)
                      Auto-detected from lock files if not specified
--output, -o          Output binary name (default: derived from directory)
--help, -h            Help for build command
```

#### Development Command

```bash
# Start development server with live reload
kirin dev

# All Air flags are supported and delegated to Air
kirin dev --build.cmd "go build -o ./tmp/main ."
kirin dev --build.bin "./tmp/main"

# Available flags
# All Air configuration flags are supported
--help, -h        Help for dev command
```

#### Generate Command

```bash
# Generate code from protobuf definitions
kirin generate

# With custom proto directory
kirin generate --proto-folder api/proto
```

**What gets generated**

* **Go backend**: Service interfaces and gRPC server stubs you can implement.
* **TypeScript frontend**: Strongly-typed client code your UI can import and call.

**Flags**

* `--proto-folder, -p`    Proto directory name (default `proto`)
* `--help, -h`            Help for generate command

#### Doctor Command

```bash
# Check system dependencies
kirin doctor

# Available aliases
kirin doc
kirin d

# No additional flags available
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the GNU GPL v3 License — see the [LICENSE](LICENSE) file for details.