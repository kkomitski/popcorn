# 🍿 Popcorn Language

A lightweight, interpreted programming language built with Go. Popcorn features a clean syntax, dynamic typing, and an interactive REPL with syntax highlighting.

## ✨ Features

- **Interactive REPL** with syntax highlighting and formatted output
- **File execution** for running `.pop` scripts
- **Dynamic typing** with runtime type checking
- **First-class functions** with closures and early returns
- **Array literals** with index-based access
- **Object literals** for structured data
- **Lexical scoping** with nested environments
- **Constants and variables** with immutability enforcement

## 🚀 Quick Start

### Installation

Run the install script to build and install Popcorn:

```bash
./install.sh
```

This will:
- Build the `pop` binary
- Install it to `~/bin`
- Add `~/bin` to your PATH (if needed)

After installation, restart your terminal or run:

```bash
source ~/.zshrc
```

### Running Popcorn

**Start the REPL:**
```bash
pop
```

**Execute a file:**
```bash
pop script.pop
```

**Uninstall:**
```bash
./uninstall.sh
```

## 📝 Syntax Guide

### Variables

```javascript
// Declare mutable variables
let x = 10;
let name = "Popcorn";

// Declare immutable constants
const PI = 3.14159;
const MAX_SIZE = 100;
```

### Numbers

```javascript
let a = 42;
let b = 3.14;
let sum = a + b;
let product = a * b;
```

### Arithmetic Operations

```javascript
let addition = 5 + 3;        // 8
let subtraction = 10 - 4;    // 6
let multiplication = 6 * 7;  // 42
let division = 20 / 4;       // 5
let modulo = 17 % 5;         // 2
```

### Functions

Functions are first-class citizens and support closures:

```javascript
// Function declaration
fn add(a, b) {
  a + b
};

// Function call
let result = add(5, 10);  // 15

// Early return with 'pop' keyword
fn max(a, b) {
  let result = a;
  
  pop result;  // Returns immediately with result
  
  // Code after 'pop' is not executed
};

// Functions can pop without a value
fn checkValue(x) {
  pop;  // Returns null immediately
};

// Functions with multiple statements
fn greet(name) {
  let message = "Hello, ";
  message
};

// Nested functions and closures
fn makeCounter() {
  let count = 0;
  fn increment() {
    count = count + 1;
    count
  };
  increment
};

let counter = makeCounter();
counter();  // 1
counter();  // 2
```

### Arrays

Arrays are ordered collections of values with zero-based indexing:

```javascript
// Array literal
let numbers = [1, 2, 3, 4, 5];

// Access elements by index
numbers[0];  // 1
numbers[2];  // 3

// Arrays can contain mixed types
let mixed = [42, "hello", { x: 10 }];

// Nested arrays
let matrix = [[1, 2], [3, 4], [5, 6]];
matrix[0];     // [1, 2]
matrix[1][0];  // 3

// Empty array
let empty = [];
```

### Objects

```javascript
// Object literal
let person = {
  name: "Alice",
  age: 30,
  city: "NYC"
};

// Access properties
person.name;  // "Alice"

// Shorthand property syntax
let x = 10;
let y = 20;
let point = { x, y };  // Same as { x: x, y: y }
```

### Assignment

```javascript
let x = 5;
x = x + 1;  // x is now 6

const y = 10;
y = 20;  // Error: Cannot reassign constant variable
```

### Comments

```javascript
// This is a single-line comment
let value = 42;  // Inline comment
```

## 🎯 REPL Commands

The interactive REPL provides an enhanced development experience:

| Command | Description |
|---------|-------------|
| `exit` | Exit the REPL |
| `clear` | Clear the screen and redisplay header |

### REPL Features

- **Syntax Highlighting**: Keywords, numbers, operators, and identifiers are color-coded
- **Pretty Output**: Results are formatted with colors and arrows
- **Persistent State**: Variables and functions persist across inputs
- **Error Messages**: Clear error reporting with helpful context

### Example Session

```bash
🍿 >> let x = 10;
   → let x = 10;
   ← {Value:10}

🍿 >> fn double(n) { n * 2 };
   → fn double(n) { n * 2 };
   ← {Name:double Params:[n] ...}

🍿 >> double(x);
   → double(x);
   ← {Value:20}
```

## 🏗️ Architecture

```
popcorn/
├── main.go              # Entry point
├── frontend/            # Lexer and Parser
│   ├── lexer.go        # Tokenization
│   ├── parser.go       # AST generation
│   └── types/
│       ├── tokens/     # Token definitions
│       └── ast/        # AST node types
├── backend/             # Interpreter
│   ├── interpreter.go  # AST evaluation
│   ├── environment.go  # Variable scoping
│   ├── types.go        # Runtime value types
│   └── run.go          # REPL and file execution
└── lib/                 # Utility functions
    └── utils.go
```

### Execution Pipeline

1. **Lexical Analysis**: Source code → Tokens (`lexer.go`)
2. **Parsing**: Tokens → Abstract Syntax Tree (`parser.go`)
3. **Evaluation**: AST → Runtime Values (`interpreter.go`)

## 🎨 Runtime Types

| Type | Description | Example |
|------|-------------|---------|
| `Number` | Floating-point numbers | `42`, `3.14` |
| `Boolean` | True/false values | `true`, `false` |
| `Null` | Null/undefined value | `null` |
| `Array` | Ordered collections | `[1, 2, 3]` |
| `Object` | Key-value collections | `{ x: 10, y: 20 }` |
| `Function` | User-defined functions | `fn add(a,b) { a+b }` |
| `NativeFunction` | Built-in functions | *(future feature)* |

## 🔧 Development

### Prerequisites

- Go 1.23.2 or higher

### Building from Source

```bash
# Clone or navigate to the project directory
cd popcorn

# Build the binary
go build -o pop main.go

# Run directly
./pop

# Or run a file
./pop examples/test.pop
```

### Running Tests

```bash
go test ./...
```

## 📚 Examples

### Fibonacci Sequence

```javascript
fn fibonacci(n) {
  let a = 0;
  let b = 1;
  let i = 0;
  
  fn loop() {
    let temp = a;
    a = b;
    b = temp + b;
    i = i + 1;
    
    i
  };
  
  loop()
};

fibonacci(10);
```

### Array Operations

```javascript
// Sum of array elements
fn sum(arr) {
  let total = 0;
  let i = 0;
  
  // Note: Length function coming soon
  total
};

// Find element in array
fn findAt(arr, index) {
  pop arr[index];
};

let numbers = [10, 20, 30, 40];
findAt(numbers, 2);  // 30
```

### Calculator

```javascript
fn calculate(op, a, b) {
  op(a, b)
};

fn add(x, y) { x + y };
fn multiply(x, y) { x * y };

calculate(add, 5, 3);       // 8
calculate(multiply, 4, 7);  // 28
```

## 🗺️ Roadmap

- [x] Array/list data structures with index access
- [x] Early return with `pop` keyword
- [x] Comparison operators (`==`, `!=`, `<`, `>`, `<=`, `>=`)
- [x] String type and string operations
- [ ] Boolean logic operators (`&&`, `||`, `!`)
- [ ] Implement boolean keywords (`true`, `false`)
- [ ] Control flow (`if`, `while`, `for`)
- [ ] Array methods (push, pop, length, map, filter)
- [ ] Built-in standard library functions
- [ ] Module system and imports
- [ ] Error handling (try/catch)
- [ ] Type annotations (optional)

## 📄 License

This project is open source and available for educational purposes.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the issues page or submit a pull request.

---

**Made with ❤️ and 🍿**
