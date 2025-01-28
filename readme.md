# GOTS

**GOTS** is a tool that generates TypeScript interfaces from Go struct definitions (transpiler). It supports embedded types, custom types, and third-party package types. This project is built to work with Go Modules and leverages `go/parser` and `go/types` for analysis and type-checking.

## Features

- Transpiles Go structs into TypeScript interfaces.
- Supports embedded and custom types.
- Handles external types from Go Modules.
- Optional fields for pointer types in Go.

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
