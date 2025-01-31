<p align="center">
<img alt="Go Status Checker" src="./gots-logo.svg" width="250" height="100" style="max-width: 100%;">
</p>

<p align="center">
A little transpiler to generate interface types from go struct types
</p>

---

**GOTS** is a tool that generates TypeScript interfaces from Go struct definitions (transpiler). It supports embedded types, custom types, and third-party package types. This project is built to work with Go Modules and leverages `go/parser` and `go/types` for analysis and type-checking.

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/BelkacemYerfa/GOTS.git
   cd GOTS
   ```
2. **Install the deps:\_**
   ```bash
   go mod tidy
   ```
3. **Build the project:**
   ```bash
   go build -o build/ main.go
   ```
   Or use Makefile:
   ```bash
   make bin
   ```
4. **Run the project:**
   ```bash
    ./build/main
   ```

## Usage

### Source file

Specify the Go source file that contains the struct definitions.

```bash
gots -s <pathname-to-file> -o <output-path>
```

### Source folder

Specify the Go source directory that contains the struct definitions.

```bash
gots -f <pathname-to-folder> -o <output-path>
```

### Custom config file

Specify the custom configuration file for the third party types.

```bash
gots -s <pathname-to-folder> -c <pathname-to-config-file> -o <output-path>
```

### Activate recursive mode

Specify the recursive mode to transpile all the files in the folder.

```bash
gots -f <pathname-to-folder> -r -o <output-path>
```

### Example

```bash
gots -s ./example.go -o ./output.ts
```

## Features Checklist

### ✅ Implemented Features

- [x] **Generate TypeScript Interfaces**
  - Transpiles Go structs into TypeScript interfaces.
- [x] **Embedded and Custom Types**
  - Supports Go embedded structs and custom types.
- [x] **External Types from Go Modules**
  - Handles third-party package types.
- [x] **Optional Fields for Pointer Types**
  - Pointer fields in Go are translated into optional fields in TypeScript.
- [x] **Recursive Mode**
  - Transpiles all files in a directory recursively.
- [x] **Custom Configuration File**
  - Supports YAML configuration files for defining third-party type mappings.

---

### 🚀 Upcoming Features

#### **General Enhancements**

- [ ] **Multi-package Support in YAML File**
  - Extend configuration file to support multiple packages.
- [ ] **Comments Transpiler**
  - Add a flag (`-tc` | `--transpile-comments`) to generate TypeScript documentation from Go comments.
- [ ] **Multi-Output Target**
  - Add a flag (`-mo` | `--multi-output`) for generating files at multiple locations or with custom file names.
- [ ] **File Versioning**
  - Embed version information in the generated files.
- [ ] **Support for TypeScript `type` Definitions**
  - Add support for generating TypeScript `type` in addition to `interface`.
- [ ] **Improved Error Handling**
  - Enhance error messages and validation for better developer experience.

#### **Mode Enhancements**

- [ ] **Enhanced Recursive Mode**
  - Add filtering options:
    - Include files by patterns (e.g., `*.go`).
    - Exclude specific files or directories (`-ex` | `--exclude testdata/*`).
  - Add recursive filtering rules in the YAML configuration.

#### **Advanced Type Support**

- [ ] **Enum Support (Go → TS)**
  - Automatically convert Go enums into TypeScript union types.

## Configuration file

The configuration file is a Yaml file that contains the third-party types and their corresponding TypeScript types.

The configuration file should have the following structure:

```yaml
config:
  gorm.Model:
    pkg: "time"
    value: 'type gormModel struct { ID uint `gorm:"primarykey"`; CreatedAt time.Time; UpdatedAt time.Time; DeletedAt time.Time `gorm:"index"`}'
```

**Notes:**

1. The `pkg` field is the packages that are used of the third-party type, and the `value` field is the Go struct definition.
2. The provided example is for the `gorm.Model` type from the `gorm` package, which is already defined in the `gots.yaml` file.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Happy hacking 🚀.
Enjoy!
