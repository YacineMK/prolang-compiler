# ProLang Compiler: Architecture & Implementation

**A complete compiler from a pedagogical language to Intel 8086 assembly, written in Go**

---

## Table of Contents

1. [Introduction & The Case for Go](#chapter-1-introduction--the-case-for-go)
2. [Lexical Analysis: The Tokenizer](#chapter-2-lexical-analysis--the-tokenizer)
3. [Semantic Analysis: Parsing, Symbol Tables & FNV Hashing](#chapter-3-semantic-analysis--parsing-symbol-tables--fnv-hashing)
4. [Intermediate Representation: Three-Address Code via Quads](#chapter-4-intermediate-representation--three-address-code-via-quads)
5. [The Optimizer: Five Passes to Better Code](#chapter-5-the-optimizer--five-passes-to-better-code)
6. [Code Generation: From Quads to 8086 Assembly](#chapter-6-code-generation--from-quads-to-8086-assembly)

---

## Chapter 1: Introduction & The Case for Go

### 1.1 What is ProLang?

ProLang is a **pedagogical programming language** designed for teaching compiler construction. Its syntax draws inspiration from Algol and Pascal, two languages that shaped the structured programming paradigm now taken for granted in every modern language. The language is intentionally simple enough that a complete compiler can be understood by a single developer, yet rich enough to demonstrate every major compiler concept: lexical analysis, recursive-descent parsing, semantic validation, intermediate representation, optimization, and target code generation.

A ProLang program has a distinctive structure that immediately signals its educational heritage:

```
BeginProject MyProgram;

Setup:
    define x, y : integer;
    define buffer : [integer; 100];
    const MaxSize : integer = 1024;

Run:
{
    x <-- 10;
    y <-- x * 2;
    out("Result: ", y);
}
EndProject;
```

The language supports two primitive types (`integer` and `float`), one-dimensional arrays, constants, and a modest set of control-flow constructs: `if/then/else/endIf`, `loop while/endloop`, and `for/in/to/endfor`. Input and output are handled through built-in `input()` and `out()` functions — the latter accepting a mix of string literals and expressions, a pragmatic convenience that mirrors how real-world logging works.

Variables are declared using the `define` keyword with an unusual **pipe syntax** for grouping: `define x | y | z : integer;` declares three integer variables at once. Assignment uses a visual arrow operator `<--` rather than the conventional `=`, a deliberate choice that reinforces the distinction between definition (`=`) and mutation (`<--`). Constants use `const Name : type = value;` and are immutable after declaration — enforced at the semantic analysis stage.

Comments come in two flavors: line comments delimited by `%%` and block comments wrapped in `//* ... *//`. The logical operators `AND`, `OR`, and `NON` use French-inspired uppercase naming, a subtle nod to the language's likely origins in a francophone academic context.

### 1.2 The Compiler Pipeline

The ProLang compiler implements a classic **multi-pass architecture** with five distinct stages, each implemented as a separate Go package:

```
Source File (.pl)
       |
       v
[ Lexer ] ------> Token Stream
       |
       v
[ Semantic Analyzer ] ------> Symbol Table + Validated Program
       |
       v
[ Quad Generator ] ------> Three-Address Code (Quads)
       |
       v
[ Optimizer ] ------> Optimized Quads (5 passes)
       |
       v
[ Assembly Generator ] ------> 8086 .asm File
```

Each stage is a pure transformation: the lexer converts raw text to tokens, the semantic analyzer validates and enriches those tokens with type information, the quad generator lowers the validated program to an intermediate representation, the optimizer improves that representation through a sequence of local transformations, and the assembly generator emits target code. The stages communicate through well-defined data structures — there is no shared mutable state between them beyond what is explicitly passed.

This architecture was chosen over a single-pass design for several reasons. First, it makes each stage testable in isolation: the lexer can be tested without a parser, the optimizer can be tested on hand-crafted quad sequences, and the assembly generator can be verified against known IR inputs. Second, it allows the intermediate representation to be inspected and debugged — the compiler prints every quad before and after optimization, giving the developer full visibility into what the optimizer did. Third, it mirrors the structure taught in compiler textbooks, making the codebase itself a pedagogical artifact.

### 1.3 Why Go?

The choice of Go as the implementation language was not arbitrary. It reflects a careful evaluation of the demands of compiler writing against the trade-offs offered by different language ecosystems.

**Performance without complexity.** A compiler is a CPU-bound program that processes potentially large inputs through multiple passes. Go compiles to native code and offers performance within a factor of 2–3 of C for most workloads, without the memory-safety nightmares that make C compilers fragile. The Go garbage collector handles the lifetime of tokens, AST nodes, and quad sequences automatically, eliminating entire categories of bugs (use-after-free, double-free, memory leaks) that plague compilers written in C or C++.

Consider what happens during a single compilation: the lexer produces hundreds of `Token` structs, the quad generator produces hundreds more `Quad` structs, and the optimizer allocates temporary maps and slices for each pass. In C, every one of those allocations would require explicit `free()` calls — and missing even one would leak memory. In Rust, the borrow checker would fight against the natural pattern of passing mutable references through recursive descent parsers. In Go, the garbage collector handles it all, and the developer focuses on the compilation logic.

**Fast compilation for fast iteration.** Compiler development involves hundreds of small edit-compile-test cycles. Go's compilation speed — often measured in milliseconds for projects of this size — means the feedback loop stays tight. A full rebuild of the ProLang compiler from scratch takes under a second, even on modest hardware. This is not a minor convenience; it directly affects how often the developer is willing to experiment with a new optimization pass or restructure the IR.

**A standard library built for systems programming.** Go's standard library includes `hash/fnv` (used for the symbol table), `strconv` (used for parsing numeric literals and formatting assembly operands), and `os` (used for file I/O). All of these would require third-party libraries or manual implementation in a language like C. The `testing` package makes unit testing trivial. Cross-compilation — producing a Windows binary from a Linux machine, or vice versa — requires nothing more than setting `GOOS` and `GOARCH` environment variables.

**Readability for a pedagogical project.** Go's syntax is famously simple: no inheritance, no generics (until recently), no operator overloading, no exceptions. For a codebase meant to be studied by students learning compiler construction, this simplicity is a feature. A reader can trace from `cmd/main.go` through each package and understand the full pipeline without needing to decode template metaprogramming or macro expansions. The explicit error handling (`if err != nil`) makes failure modes visible at every call site rather than hiding them behind exception handlers.

**A single static binary.** The compiled ProLang compiler is a single executable with no runtime dependencies — no DLLs, no JVM, no shared libraries. A student can copy the binary to any machine running the same OS and architecture and compile ProLang programs immediately. This is particularly valuable in a teaching environment where students may have diverse and unpredictable machine configurations.

The one external dependency is `charmbracelet/lipgloss`, a terminal styling library that provides the colored, box-drawn output used in the compiler's TUI. It was chosen over raw ANSI escape codes because it handles terminal capability detection and provides a clean, composable API for building styled output. The dependency is purely cosmetic — the compiler would function identically with plain `fmt.Println` calls — but the visual polish makes the tool feel professional and helps students distinguish error messages from informational output at a glance.

### 1.4 The Target: 8086 Real-Mode Assembly

The compiler targets **Intel 8086 real-mode assembly** in MASM/TASM syntax. This is an unconventional choice in 2026 — why not target x86-64, ARM, or WebAssembly?

The answer is context. The 8086 is the processor at the heart of the original IBM PC. Its architecture — 16-bit registers, segmented memory, a small but complete instruction set — is simple enough to be taught in a few lectures, yet real enough that the generated code runs on actual hardware (or, more practically, in DOSBox, the de facto standard DOS emulator). Students can single-step through their compiled programs in Turbo Debugger, watching registers change and memory locations update in real time. They can inspect the stack, set breakpoints, and modify values — all within a visual debugger that was designed for exactly this kind of exploration.

This is fundamentally different from targeting a modern architecture. An x86-64 assembly listing is overwhelming: 16 general-purpose registers, SSE/AVX vector extensions, multiple addressing modes, position-independent code requirements, and platform-specific calling conventions. A student trying to understand why their `out("Hello")` call segfaults would need to debug through glibc, the dynamic linker, and the kernel's syscall interface. On the 8086 under DOS, that same call is three instructions: load the string address, set the function code, trigger the interrupt.

The trade-off is that floating-point operations are not truly supported — all values are truncated to integers at the assembly level — and the input routine only reads single decimal digits. These limitations are documented and understood; they are acceptable for a pedagogical compiler where the goal is understanding the pipeline, not producing production binaries.

---

## Chapter 2: Lexical Analysis — The Tokenizer

### 2.1 Design Overview

The lexer is the compiler's first point of contact with source code. Its job is conceptually simple: consume a stream of characters and produce a stream of tokens, discarding whitespace and comments along the way. In practice, the lexer must handle dozens of edge cases — multi-character operators, escaped strings, numeric formats, keyword disambiguation — while maintaining precise line and column tracking for error reporting.

The ProLang lexer is implemented in `internal/lexer/lexer.go` and `internal/lexer/token.go`, totaling approximately 380 lines of Go. It is a **hand-written scanner** rather than a generated one (no `lex` or `ragel`), giving full control over error recovery and token representation.

The core data structure is the `Lexer` struct:

```go
type Lexer struct {
    source []rune   // source code as rune slice
    pos    int      // current position in source
    line   int      // current line (1-indexed)
    col    int      // current column (1-indexed)
    Tokens []Token  // accumulated token list
    Errors []string // accumulated error messages
}
```

The use of `[]rune` rather than `[]byte` is deliberate. Go source code is UTF-8, but indexing into a UTF-8 byte slice is error-prone — a multi-byte character would be split across indices, producing garbage. By converting to runes upfront, the lexer can advance one character at a time regardless of encoding width, and `line`/`col` tracking remains accurate for any Unicode identifier.

### 2.2 The Token Set

Every token in ProLang is represented by a `Token` struct carrying four fields:

```go
type Token struct {
    Kind  TokenKind // the category (e.g., IDENT, INT_CONST, IF)
    Value string    // the lexeme as it appeared in source
    Line  int       // source line (1-indexed)
    Col   int       // source column (1-indexed)
}
```

`TokenKind` is a string type alias rather than an integer enum. This is a pragmatic Go idiom: it makes debug output self-documenting (a token appears as `{IF "if" line:12 col:5}` rather than `{17 "if" line:12 col:5}`), and it costs nothing at runtime since string constants are interned by the compiler.

The token set is organized into logical groups:

| Category | Tokens | Purpose |
|----------|--------|---------|
| **Literals** | `INT_CONST`, `FLOAT_CONST`, `STRING`, `IDENT` | Value-bearing tokens |
| **Structure** | `BEGINPROJECT`, `ENDPROJECT`, `SETUP`, `RUN` | Program scaffolding |
| **Declarations** | `DEFINE`, `CONST`, `INTEGER`, `FLOAT` | Variable/constant definition |
| **Control Flow** | `IF`, `THEN`, `ELSE`, `ENDIF`, `LOOP`, `WHILE`, `ENDLOOP`, `FOR`, `IN`, `TO`, `ENDFOR` | Branching and looping |
| **Operators** | `PLUS`, `MINUS`, `STAR`, `SLASH`, `LT`, `GT`, `LTE`, `GTE`, `EQ`, `NEQ` | Arithmetic and comparison |
| **Logic** | `AND`, `OR`, `NON` | Boolean connectives |
| **Punctuation** | `PIPE`, `COLON`, `SEMI`, `COMMA`, `LPAREN`, `RPAREN`, `LBRACE`, `RBRACE`, `LBRACK`, `RBRACK` | Delimiters |
| **Assignment** | `ASSIGN` | The `<-` / `<--` operator |
| **Other** | `EQUAL`, `IN_FUNC`, `OUT_FUNC`, `EOF`, `ILLEGAL` | Special-purpose tokens |

Keywords are stored in a `map[string]TokenKind` literal:

```go
var Keywords = map[string]TokenKind{
    "BeginProject": BEGINPROJECT,
    "EndProject":   ENDPROJECT,
    "Setup":        SETUP,
    "Run":          RUN,
    "define":       DEFINE,
    "const":        CONST,
    // ... etc
}
```

When the lexer encounters an identifier-like token, it performs a single map lookup to determine if it's a keyword or a user-defined name. This is O(1) and trivial to extend — adding a keyword is a one-line change to the map.

### 2.3 The Scanning Loop

The main scanning loop in `Lex()` is deceptively simple:

```go
func (l *Lexer) Lex() {
    for l.pos < len(l.source) {
        l.scan()
    }
    l.Tokens = append(l.Tokens, Token{Kind: EOF, Line: l.line, Col: l.col})
}
```

All the complexity lives in `scan()`, which is a single 120-line method that dispatches on the current character. The dispatch order matters:

1. **Whitespace** is skipped first — spaces, tabs, newlines, and carriage returns advance the position without emitting tokens.

2. **Comments** are detected next. A `%%` triggers line-comment mode: the scanner advances until it hits a newline or EOF, consuming everything in between. A `//*` triggers block-comment mode: the scanner advances until it finds the closing `*//` sequence, tracking line count for nested newlines. Block comments that reach EOF without closing produce an error.

3. **Strings** begin with `"`. The `scanString()` method accumulates characters until the closing quote, handling the case of unterminated strings (newline before closing quote) with a specific error message. The string value is stored without the surrounding quotes.

4. **Numbers** begin with a digit. The `scanNumber()` method accumulates digit characters, then checks for a decimal point. If found, it switches to float mode and accumulates fractional digits. The token is emitted as `INT_CONST` or `FLOAT_CONST` accordingly.

5. **Identifiers and keywords** begin with a letter. The `scanIdent()` method accumulates alphanumeric characters and underscores, then checks the keyword map. A special case handles the `in` keyword: if followed by whitespace and `(`, it's classified as `IN_FUNC` (the input built-in); otherwise it's the `IN` keyword used in `for` loops. This two-token lookahead is the only place the lexer needs context beyond the current character.

6. **Multi-character operators** are matched greedily. The scanner checks for `<=`, `>=`, `==`, `!=`, `<--` (the long arrow assignment), `<-` (the short arrow), `<`, and `>` in that order. The greedy approach ensures that `<=` is recognized as a single `LTE` token rather than `LT` followed by `EQUAL`.

7. **Single-character tokens** are matched via a map lookup. Characters like `+`, `-`, `*`, `/`, `|`, `:`, `;`, `,`, `(`, `)`, `{`, `}`, `[`, `]` each map to their corresponding token kind.

8. **Anything else** is an illegal character. The scanner emits an `ILLEGAL` token and records an error with the exact line and column.

This dispatch structure is intentionally flat. A table-driven DFA would be more efficient for large grammars, but for ProLang's ~40 token types, the if/switch chain is readable, debuggable, and fast enough — a full lexical pass over a 400-token program completes in microseconds.

### 2.4 Position Tracking

Accurate error reporting requires precise line and column information. The lexer tracks position through two mechanisms:

The `advance()` method is the only function that consumes a character. It increments `pos` and `col`, and when it encounters a `\n`, it increments `line` and resets `col` to 1:

```go
func (l *Lexer) advance() rune {
    r := l.source[l.pos]
    l.pos++
    if r == '\n' {
        l.line++
        l.col = 1
    } else {
        l.col++
    }
    return r
}
```

Lookahead is provided by `peek()` and `peekAt(offset int)`, which return the rune at a future position without consuming it. These are used for multi-character operator detection and keyword disambiguation without disturbing the position state.

### 2.5 Error Recovery Strategy

The ProLang lexer uses a **permissive error recovery** strategy. When it encounters an illegal character, it does not abort — it emits an `ILLEGAL` token, records the error, and continues scanning. This means a single source file can produce both valid tokens and error annotations, letting the developer fix multiple issues in one compilation cycle.

Lexical errors are also **non-fatal**: the main pipeline in `cmd/main.go` prints lexical errors but continues to semantic analysis. This is in contrast to semantic errors, which abort the compilation. The reasoning is that a misspelled keyword or stray character rarely prevents the parser from understanding the program's structure, whereas a type mismatch or undeclared variable makes further analysis meaningless.

---

## Chapter 3: Semantic Analysis — Parsing, Symbol Tables & FNV Hashing

### 3.1 The Dual Role of the Analyzer

The semantic analyzer in `internal/semantic/analyzer.go` serves two roles that are traditionally separated: **parsing** (determining whether the token sequence conforms to the language grammar) and **semantic validation** (checking that the program follows type rules, variables are declared before use, constants are not reassigned, etc.).

This unification is a pragmatic choice for a language of ProLang's size. The grammar is LL(1) for the most part, meaning a recursive-descent parser can make decisions by looking at a single token of lookahead. The entire analyzer fits in approximately 500 lines of Go — small enough that separating parsing from validation would add architectural overhead without corresponding clarity.

The analyzer is structured as a **recursive-descent parser** — a pattern where each grammar non-terminal is implemented as a method, and the structure of the code mirrors the grammar. For example, expression parsing follows the classic precedence-climbing pattern:

```
parseExpr()      ->  parseAddSub()      (lowest precedence: +, -)
parseAddSub()    ->  parseMulDiv()      (middle precedence: *, /)
parseMulDiv()    ->  parsePrimary()     (highest precedence: literals, identifiers, parens)
```

Each method returns type information (so the caller can check type compatibility) but does not build an AST node. Instead, the analyzer performs validation inline — checking that an integer variable is not assigned a float expression, that constants are not used as assignment targets, that loop variables exist and have the correct type.

### 3.2 The Symbol Table

The symbol table is the central data structure of semantic analysis. It maps names (identifiers) to their properties: what kind of symbol they are (variable, constant, array), what type they hold (integer, float), what their initial value is, and whether they have been assigned to.

The table is implemented in `internal/ast/symbol.go` as a **fixed-size hash table with separate chaining**:

```go
const TABLE_SIZE = 1024

type SymbolTable struct {
    buckets [TABLE_SIZE]*entry
    count   int
}

type entry struct {
    key  string
    sym  Symbol
    next *entry
}
```

The fixed size of 1024 buckets is appropriate for the expected use case. A typical ProLang program declares between 10 and 50 symbols. With 1024 buckets and a reasonable hash function, the average chain length is near zero — most lookups hit on the first comparison. Even a pathological program with hundreds of declarations would see chain lengths in the single digits.

Insertion uses **head insertion** on the bucket chain: the new entry is linked to the current bucket head, then replaces it. This is O(1) after the O(1) duplicate check (a lookup on the same key). Lookup computes the hash, walks the chain comparing keys, and returns the found symbol or nil. Updates find the existing symbol by lookup and modify its `Value` or `Initialized` fields in place.

The `Print()` method renders the table as a Unicode box-drawn table, a nice touch that makes debugging symbol resolution straightforward:

```
┌──────┬──────────────┬───────────┬─────────┬───────┬────────────┐
│ Idx  │ Nom          │ Catégorie │ Type    │ Ligne │ Valeur     │
├──────┼──────────────┼───────────┼─────────┼───────┼────────────┤
│ 135  │ x            │ variable  │ integer │ 4     │            │
│ 452  │ i            │ variable  │ integer │ 5     │            │
└──────┴──────────────┴───────────┴─────────┴───────┴────────────┘
```

### 3.3 FNV Hashing: Why and How

The symbol table uses the **FNV-1a (Fowler-Noll-Vo) hash function** via Go's standard library `hash/fnv` package:

```go
func hashStr(key string) int {
    h := fnv.New32a()
    h.Write([]byte(key))
    return int(h.Sum32()) % TABLE_SIZE
}
```

The choice of FNV-1a over other hash functions is worth examining in detail.

**What is FNV?** Fowler-Noll-Vo is a family of non-cryptographic hash functions designed for speed and good distribution. The "1a" variant processes each byte by XOR-ing it with the current hash value, then multiplying by a prime constant. For the 32-bit variant, the prime is 16777619. The algorithm is:

```
hash = 2166136261 (offset basis)
for each byte b:
    hash = hash XOR b
    hash = hash * 16777619
```

In Go's implementation, this is approximately:

```go
func fnv32a(data []byte) uint32 {
    h := uint32(2166136261)
    for _, b := range data {
        h ^= uint32(b)
        h *= 16777619
    }
    return h
}
```

**Why not DJB2?** The original implementation of this compiler used DJB2 (Bernstein hash), another classic non-cryptographic hash. DJB2 is:

```
hash = 5381
for each byte b:
    hash = hash * 33 + b
```

The switch from DJB2 to FNV-1a (visible in commit `f41f9b1`) was motivated by FNV's superior **avalanche property** — the tendency for a single-bit change in the input to flip roughly half the bits in the output. DJB2 with multiplier 33 produces weaker mixing for short inputs (which dominate in symbol tables, where most identifiers are under 12 characters). FNV-1a's XOR-then-multiply structure ensures that even single-character differences between identifiers ("x" vs "y", "temp1" vs "temp2") produce widely separated bucket indices.

**Why not a cryptographic hash?** SHA-256 or BLAKE3 would produce perfect distribution but at significant cost. A symbol table lookup happens on every identifier reference — potentially hundreds of times per compilation. Cryptographic hashes are designed to resist preimage attacks, not to be fast. FNV-1a processes each byte with a single XOR and multiply; SHA-256 processes each 64-byte block through 64 rounds of bit rotations, shifts, and table lookups. The performance difference for short strings is roughly 50–100x, and the improved distribution offers no practical benefit when the table has 1024 buckets and fewer than 100 entries.

**Why not Go's built-in map?** This is the most interesting question. Go's `map[string]Symbol` would handle hashing, collision resolution, and resizing automatically. Why implement a custom hash table?

The answer is pedagogical. A compiler's symbol table is one of its most important data structures, and understanding how it works — hashing, chaining, collision resolution — is a core learning objective. Using Go's built-in map would hide all of that behind a black-box abstraction. The custom implementation exposes every detail: the bucket array, the entry chain, the hash computation, the head-insertion strategy. A student can set a breakpoint in `hashStr()` and watch the hash value be computed, or inspect `buckets` after insertion to see how collisions are resolved.

Additionally, the fixed-size design avoids the complexity of rehashing. Go's built-in map grows automatically when the load factor exceeds a threshold, which would make the bucket indices unstable across insertions. The fixed 1024-bucket table guarantees that a symbol inserted at index 135 will always be found at index 135, simplifying debugging and reasoning about the table's behavior.

### 3.4 Parsing Declarations

Declaration parsing handles the `Setup:` section of a ProLang program. The grammar is:

```
Setup:
    (define | const)*
Run:
```

**Variable declarations** use the `define` keyword followed by a pipe-separated list of names, a colon, and a type:

```
define x | y | z : integer;
define arr : [integer; 50];
define total : integer = 0;
```

The parser collects names in a slice, handling the pipe separators, then dispatches based on whether `[` follows (array declaration) or not (scalar declaration). For arrays, it validates that the size is a positive integer. For scalars, it checks for an optional `= value` initializer and validates type compatibility — a float literal like `3.14` cannot initialize an integer variable.

**Constant declarations** are simpler:

```
const Pi : float = 3.14159;
const Max : integer = 100;
```

The parser validates that the declared type matches the literal's type and inserts the symbol as `KIND_CONST` with `Initialized: true`. Constants are permanently immutable — any attempt to assign to a constant is caught by the assignment parser.

After all declarations are processed, the analyzer calls `table.Print()`, rendering the symbol table so the developer can verify that every declared symbol was correctly categorized.

### 3.5 Expression Parsing with Precedence

Expression parsing uses the standard **precedence climbing** technique. The grammar is:

```
expr       -> add_sub
add_sub    -> mul_div ( (+ | -) mul_div )*
mul_div    -> primary ( (* | /) primary )*
primary    -> INT_CONST | FLOAT_CONST | IDENT [ '[' expr ']' ] | '(' expr ')'
```

Each level calls the next higher precedence level for operands, then loops on its own operators. This naturally enforces that `*` and `/` bind tighter than `+` and `-`:

```
a + b * c    parses as    a + (b * c)    // correct
a * b + c    parses as    (a * b) + c    // correct
```

The parser also handles **unary minus** through a special case for parenthesized signed literals like `(-5)` and `(+3.14)`. These are detected in `parsePrimary()` and converted to negative/positive constant expressions.

Type checking during expression parsing uses `MergeTypes()`, which implements the standard arithmetic promotion rules: if either operand is `float`, the result is `float`; otherwise it's `integer`. This allows mixed-type expressions like `a + Pi` (where `a` is integer and `Pi` is float) to be correctly typed as float expressions, while flagging an error if the result is assigned to an integer variable.

### 3.6 Statement Parsing

**Assignment** is the most common statement. The parser distinguishes between scalar assignment (`x <-- expr`) and array assignment (`Tab[i] <-- expr`). For scalars, it validates that the target exists in the symbol table, is not a constant, and is type-compatible with the right-hand side. For arrays, it additionally validates that the target is actually declared as an array and that the index expression is valid.

**If statements** support an optional `else` branch:

```
if (condition) then:
    { ... }
else:
    { ... }
endIf;
```

The parser validates that the condition is a boolean expression (comparison or logical combination), then recursively parses both branches, checking that variables assigned in one branch are compatible with usage after the if statement.

**Loop statements** come in two flavors. The `loop while` form:

```
loop while (condition)
    { ... }
endloop;
```

And the `for` form:

```
for variable in start_expr to end_expr
    { ... }
endfor;
```

The `for` parser validates that the loop variable is declared and has integer type (floating-point loop counters are prohibited), and that both the start and end expressions produce integer-compatible values.

### 3.7 Error Handling in the Analyzer

Semantic errors are **fatal**: the analyzer collects all errors in a slice, and if the slice is non-empty after `Analyze()` returns, `main.go` prints all errors and exits with code 1. This is appropriate because semantic errors indicate fundamental problems — undeclared variables, type mismatches, invalid control flow — that would produce meaningless quads or assembly if ignored.

Error messages include the source location (line and column) and the problematic token value:

```
SEMANTIC_ERROR [ligne:15, col:10] 'x': variable 'x' non déclarée
SEMANTIC_ERROR [ligne:22, col:5] 'Pi': impossible de modifier une constante
```

The French-language error messages are a deliberate choice reflecting the language's likely use in francophone educational settings.

---

## Chapter 4: Intermediate Representation — Three-Address Code via Quads

### 4.1 Why an Intermediate Representation?

Between the semantic analyzer and the assembly generator sits an **intermediate representation** (IR). An IR decouples the frontend (which understands the source language) from the backend (which understands the target machine). This decoupling has three major benefits:

1. **Optimization becomes language- and target-agnostic.** The optimizer works on quads without knowing whether they came from ProLang, C, or any other language, and without knowing whether the target is 8086, ARM, or a virtual machine. An optimization pass written once applies to all combinations.

2. **The assembly generator becomes simpler.** Instead of translating complex nested control flow directly to jumps and labels, it translates flat sequences of simple operations. A `while` loop becomes a sequence of quads with `BF` (branch-if-false) and `BR` (branch) operations — the assembly generator doesn't need to know about loop semantics.

3. **Debugging becomes tractable.** The compiler prints every quad before and after optimization. A developer can inspect the IR to determine whether a bug is in the frontend (wrong quads generated), the optimizer (quads incorrectly transformed), or the backend (correct quads translated to wrong assembly).

### 4.2 The Quad Data Structure

A **quad** (short for quadruple) is a four-tuple representing a single operation in three-address code form:

```go
type Quad struct {
    Op     string // operation code
    Arg1   string // first operand
    Arg2   string // second operand
    Result string // destination
}
```

The name "three-address code" comes from the fact that most operations involve at most three addresses (two source operands and one destination). The quad representation makes this explicit: `(+, x, y, T0)` means "add x and y, store the result in T0," and `(:=, 42, , x)` means "assign 42 to x."

The ProLang IR uses these operation codes:

| Op | Meaning | Example |
|----|---------|---------|
| `:=` | Assignment | `(:=, 10, , x)` — x = 10 |
| `+`, `-`, `*`, `/` | Arithmetic | `(+, x, y, T0)` — T0 = x + y |
| `NEG` | Unary negation | `(NEG, x, , T0)` — T0 = -x |
| `>`, `<`, `>=`, `<=`, `==`, `!=` | Comparison | `(>, x, 5, T0)` — T0 = (x > 5) |
| `AND`, `OR` | Logical connectives | `(AND, T0, T1, T2)` — T2 = T0 AND T1 |
| `NON` | Logical NOT | `(NON, T0, , T1)` — T1 = NOT T0 |
| `BF` | Branch if false | `(BF, T0, , L5)` — if !T0 goto L5 |
| `BR` | Unconditional branch | `(BR, , , L2)` — goto L2 |
| `[]` | Array read | `([], Tab, i, T0)` — T0 = Tab[i] |
| `[]=` | Array write | `([]=, val, i, Tab)` — Tab[i] = val |
| `in` | Input | `(in, , , x)` — read x |
| `out` | Output expression | `(out, x, , )` — print x |
| `out_str` | Output string | `(out_str, "Hello", , )` — print "Hello" |
| `LABEL` | Informational label | Not emitted as code |

### 4.3 The Quad Manager

The `quad.Manager` is the factory for IR construction:

```go
type Manager struct {
    Quads      []Quad
    tempCount  int
    labelCount int
}
```

It provides three key services:

**Emission:** `Emit(op, arg1, arg2, result string) int` appends a quad and returns its index. The index is used for backpatching (see below).

**Temporary generation:** `NewTemp() string` returns names like `T0`, `T1`, `T2`, incrementing an internal counter. Each temporary is guaranteed unique, and by construction, each temporary is assigned exactly once — a property that the optimizer exploits for safe constant propagation.

**Label generation:** `NewLabel() string` returns names like `L0`, `L1`, `L2`. Labels are used as branch targets in `BF` and `BR` quads. Like temporaries, they are unique and monotonically increasing.

### 4.4 Control Flow and Backpatching

The most interesting aspect of quad generation is **backpatching** — the technique for resolving forward branch targets when the target address is not yet known.

Consider an `if` statement:

```
if (x > 5) then:
    x <-- 0;
else:
    x <-- 1;
endIf;
out(x);
```

When the quad generator encounters the `if`, it emits quads for the condition (`x > 5` → `T0`), then emits a `BF` (branch-if-false) quad with an **empty target**. At this point, the generator doesn't know how many quads the then-block will produce, so it can't compute the jump offset. It records the index of the `BF` quad and continues.

After generating the then-block, if there's an `else`, the generator emits a `BR` (unconditional branch) with an empty target — this will jump over the else-block after the then-block executes. Now the generator knows where the else-block starts (the current quad index), so it **backpatches** the `BF` quad: `qm.Backpatch(bfIndex, elseLabel)`. After generating the else-block, it backpatches the `BR` quad to point past the `endIf`.

The backpatching function is simple:

```go
func (qm *Manager) Backpatch(index int, label string) {
    qm.Quads[index].Result = label
}
```

The same technique handles loops. A `loop while`:

```
loop while (i < 10):
    i <-- i + 1;
endloop;
```

The generator records the index of the condition check. After generating the body, it emits a `BR` to the condition check's index. The `BF` at the end of the condition is backpatched to point past the `endloop` (the first quad after the loop body).

For `for` loops, there's an additional complexity: the loop variable must be initialized before the condition check, and incremented after the body. The quad sequence for `for i in 0 to 10 { body }` is:

```
L_init:   := 0 i
L_cond:   <= i 10 Tcmp
          BF Tcmp L_exit
          ... body ...
          + i 1 Tinc
          := Tinc i
          BR L_cond
L_exit:   ...
```

### 4.5 Array Access in Three-Address Code

Array operations are a revealing example of how high-level language features decompose into IR. The source statement:

```
Tab[i + 1] <-- x * 2;
```

generates this quad sequence:

```
+ i 1 T_idx        // T_idx = i + 1        (compute index)
* x 2 T_val        // T_val = x * 2        (compute value)
[]= T_val T_idx Tab // Tab[T_idx] = T_val  (store to array)
```

Each sub-expression gets its own temporary, and the final array store is a single quad with the array name in the `Result` field (for `[]=`) or the value in `Result` (for `[]` reads). This decomposition is what makes optimization possible — the index computation can be constant-folded if `i` is known, and the value computation can be simplified independently of the array operation.

---

## Chapter 5: The Optimizer — Five Passes to Better Code

### 5.1 Optimization Philosophy

The ProLang optimizer performs **local, intra-procedural** optimization on the quad sequence. It does not perform global analysis (no control-flow graph, no data-flow analysis, no alias analysis). Each pass scans the quad list linearly and applies transformations based on local information — the current quad and its immediate operands.

This design reflects a deliberate trade-off. Global optimizations (like loop-invariant code motion or global value numbering) would require building and analyzing a control-flow graph, which would roughly double the complexity of the optimizer. The local passes, while less powerful, catch the most common cases: expressions involving only constants, repeated computations, and assignments to dead variables. For a pedagogical compiler, this is the right balance — the passes are simple enough to be understood in a single reading, yet effective enough to demonstrate tangible improvements.

The pipeline runs five passes in order:

```
propagate -> fold -> simplify -> cse -> deadCode
```

Each pass transforms the quad list and passes it to the next. The order matters: propagation resolves variable references to constants, which enables folding to evaluate expressions at compile time; simplification cleans up after folding (e.g., `x + 0` → `x`); CSE finds repeated expressions; and dead code elimination removes the now-unnecessary temporaries.

### 5.2 Pass 1: Copy and Constant Propagation

The `propagate` pass is the most complex — and the one where the most significant bugs were found and fixed. Its purpose is to replace variable references with their known values. If `x` is assigned `10` and never reassigned, every later use of `x` can be replaced with `10`, enabling the folding pass to evaluate expressions like `x + 5` at compile time.

The pass maintains an **environment map** — a Go `map[string]string` from variable names to their known values. For each quad, it resolves both operands through the environment (following chains transitively: if `env[a] = "b"` and `env[b] = "10"`, then `resolve("a")` returns `"10"`). Then it updates the environment:

- If the quad is `:= literal target` and the target is safe to propagate, add `env[target] = literal`.
- If the quad is `:= temp target` and the target is safe to propagate, add `env[target] = temp`.
- Otherwise, if the quad has a result, delete it from the environment (its value is no longer the known constant).

**The multi-assignment bug.** The original implementation had a critical flaw: it propagated values for variables that were assigned **multiple times**. Consider a loop:

```
x <-- 0;
loop while (x < 5):
    x <-- x + 1;
endloop;
out(x);
```

The loop condition `x < 5` appears **before** the loop body in the quad list. When the propagator reaches the condition, `env[x] = "0"` (from the initial assignment). It substitutes `0` for `x`, producing `0 < 5`, which the folding pass evaluates to `1` (always true). The loop condition is now a constant, and the loop runs forever — the program never reaches `out(x)`.

A second manifestation occurs with if/else:

```
if (cond) then:
    result <-- 10;
else:
    result <-- 20;
endIf;
out(result);
```

The propagator processes the then-branch's `result := 10` (adding `env[result] = "10"`), then the else-branch's `result := 20` (overwriting `env[result] = "20"`). The final `out(result)` resolves to `20`, regardless of which branch executes at runtime. The result is always wrong when the condition is true.

**The fix** is a pre-scan that counts how many times each variable appears as the target of a `:=` quad:

```go
assignCount := map[string]int{}
for _, q := range quads {
    if q.Op == ":=" {
        assignCount[q.Result]++
    }
}
```

Then, only variables with `assignCount <= 1` are added to the environment. This implements an **SSA-like** constraint: variables that are assigned exactly once (like temporaries, which the quad generator creates fresh each time) are safe to propagate; variables assigned multiple times (loop counters, accumulators, variables set in both if/else branches) are read from memory at every use, preserving runtime-correct behavior.

This fix is conservative — it misses some optimization opportunities (e.g., a variable assigned twice where the first definition never reaches the second use because of control flow) — but it is **sound**: it never produces incorrect code. The alternative would require building a control-flow graph and performing reaching-definition analysis, which is substantially more complex.

### 5.3 Pass 2: Constant Folding

The `fold` pass evaluates operations where both operands are compile-time constants. For example:

```
(+, 10, 5, T0)  →  (:=, 15, , T0)
(*, 3, 7, T1)   →  (:=, 21, , T1)
(>, 10, 5, T2)  →  (:=, 1, , T2)
```

The `eval()` function handles the actual computation:

```go
func eval(op, a, b string) (string, bool) {
    fa, _ := strconv.ParseFloat(a, 64)
    fb, _ := strconv.ParseFloat(b, 64)
    switch op {
    case "+": return fmtFloat(fa + fb), true
    case "-": return fmtFloat(fa - fb), true
    // ... etc
    }
}
```

Folding also handles unary negation of literals: `(NEG, 5, , T0)` → `(:=, -5, , T0)`.

The `fmtFloat()` helper formats the result intelligently: if the float value has no fractional part (e.g., `15.0`), it's formatted as the integer `15` rather than `15.000000`. This keeps the IR readable and avoids unnecessary type distinctions in the assembly output.

Folding is safe because it only triggers when both operands are literal numeric strings — there is no risk of incorrectly "folding" an expression that depends on runtime values. The heavy lifting is done by propagation (Pass 1), which substitutes constants for variable references, creating the literal-literal patterns that folding can then evaluate.

### 5.4 Pass 3: Algebraic Simplification

The `simplify` pass applies algebraic identities that are always valid regardless of operand values:

| Pattern | Replacement | Rationale |
|---------|-------------|-----------|
| `x + 0` / `0 + x` | `x` | Additive identity |
| `x - 0` | `x` | Subtractive identity |
| `x * 1` / `1 * x` | `x` | Multiplicative identity |
| `x * 0` / `0 * x` | `0` | Zero product |
| `0 / x` | `0` | Zero numerator |
| `x / 1` | `x` | Divisive identity |

These simplifications often trigger after folding reduces part of an expression to a constant. For instance, `x + (5 - 5)` folds to `x + 0`, which simplify then reduces to just `x`. Without simplification, the assembly generator would emit an unnecessary `ADD AX, 0` instruction.

Note that division by zero is not caught here — it would be a runtime error. The compiler could add a compile-time check for literal zero divisors, but since the target is DOS (where a division-by-zero interrupt is well-defined), this is left to the runtime.

### 5.5 Pass 4: Common Subexpression Elimination (CSE)

The `cse` pass identifies repeated computations and reuses the first result:

```
(+, a, b, T0)   // first occurrence
...              // (a and b are not modified here)
(+, a, b, T1)   // second occurrence — same computation
```

The second quad is replaced with `(:=, T0, , T1)`, saving a redundant addition at runtime.

The pass maintains a map from expression keys (`op|arg1|arg2`) to the temporary holding the result. When it encounters an arithmetic quad (`+`, `-`, `*`, `/`), it checks whether the expression has been seen before. If so, it emits a copy from the cached temporary instead of a fresh computation.

**The user-variable safety fix.** Like the propagator, the original CSE implementation had a subtle correctness bug. It cached ALL expressions, including those with user-variable operands:

```
(+, x, 5, T0)   // cached as "+|x|5" → T0
...
x <-- 10;        // x is reassigned!
...
(+, x, 5, T1)   // same key "+|x|5" — reused T0, but x changed!
```

The second expression would reuse `T0` (which holds `old_x + 5`), but `x` has been reassigned. The result is incorrect.

The fix complements the propagation fix: CSE now skips expressions where either operand is a user variable (not a temp, not a literal). Temps are safe because each is assigned exactly once; literals are safe because they never change. User variables, which may be reassigned, are excluded:

```go
if isUserVar(q.Arg1) || isUserVar(q.Arg2) {
    out = append(out, q)  // don't cache
    continue
}
```

This is conservative — it misses CSE opportunities for user variables that happen to be single-assignment — but it is sound. The propagator already converts single-assignment user variables to their literal or temp values, so they appear as literals/temps in CSE anyway.

### 5.6 Pass 5: Dead Code Elimination

The `deadCode` pass removes assignments to temporaries that are never used. A temporary `T_k` is dead if the `uses` count (computed by scanning all `Arg1` and `Arg2` fields) is zero:

```go
for _, q := range quads {
    if q.Op == ":=" && isTemp(q.Result) && uses[q.Result] == 0 {
        keep[i] = false
    }
}
```

Dead temporaries arise naturally from earlier passes. When folding replaces `(+, 10, 5, T0)` with `(:=, 15, , T0)`, and `T0` is used exactly once (by the quad that originally consumed the addition result), the `:=` becomes a dead store if the consumer already got the constant 15 propagated into it.

After removing dead quads, branch targets must be remapped. If quad 7 (a `BF` target) was removed, the branch must point to the next surviving quad. The `remapTarget()` function handles this:

```go
func remapTarget(target string, oldToNew map[int]int, oldLen int) string {
    old, _ := strconv.Atoi(target)
    if nt, ok := oldToNew[old]; ok {
        return fmt.Sprintf("%d", nt)  // target survived; use new index
    }
    for j := old + 1; j < oldLen; j++ {
        if nt, ok := oldToNew[j]; ok {
            return fmt.Sprintf("%d", nt)  // target removed; find next
        }
    }
    return target
}
```

This remapping preserves the semantics: if a branch targeted a now-deleted no-op assignment, jumping to the next instruction is equivalent.

### 5.7 Optimization Example

Consider this ProLang fragment:

```
x <-- 10;
y <-- 5;
z <-- x + y;
out(z);
```

The unoptimized quads:
```
000: (:=, 10, , x)
001: (:=, 5, , y)
002: (+, x, y, T0)
003: (:=, T0, , z)
004: (out, z, , )
```

After propagation (x→10, y→5, z→T0):
```
000: (:=, 10, , x)
001: (:=, 5, , y)
002: (+, 10, 5, T0)
003: (:=, T0, , z)
004: (out, T0, , )
```

After folding (10+5=15):
```
000: (:=, 10, , x)
001: (:=, 5, , y)
002: (:=, 15, , T0)
003: (:=, T0, , z)
004: (out, T0, , )
```

After dead code elimination (x, y, z never used — only T0 is used by `out`):
```
002: (:=, 15, , T0)
004: (out, T0, , )
```

The final assembly is simply:
```asm
MOV AX, 15
CALL PRINT_INT
```

The original 5 quads (and all runtime operations on x, y, z) have been reduced to a single constant load and a print call.

---

## Chapter 6: Code Generation — From Quads to 8086 Assembly

### 6.1 The 8086 Memory Model

The 8086 processor operates in a **segmented memory model**. Memory is divided into logical segments, each up to 64KB, addressed by a segment register (CS for code, DS for data, SS for stack) plus a 16-bit offset. The generated assembly follows this model with three segments:

```asm
STACK SEGMENT STACK
    DW 100 DUP(?)
STACK ENDS

DATA SEGMENT
    ; variables, temporaries, string constants
DATA ENDS

CODE SEGMENT
    ASSUME CS:CODE, DS:DATA, SS:STACK
START:
    MOV AX, DATA
    MOV DS, AX
    ; ... program body ...
    MOV AH, 4Ch
    INT 21h
CODE ENDS
END START
```

The `STACK SEGMENT STACK` directive tells the linker this is the stack segment, and `DW 100 DUP(?)` reserves 100 words (200 bytes) of uninitialized stack space. The `START:` label is the program entry point. The first two instructions load the DS segment register with the data segment address — without this, all variable accesses would target the wrong memory.

### 6.2 Data Segment Layout

The data segment is populated from the symbol table and quad analysis. User variables are emitted first, ordered as the symbol table iterator returns them:

```asm
DATA SEGMENT
    ; user variables
Tabfloat DW 30 DUP(?)
x DW ?
k DW ?
somme DW 0
Tabint DW 50 DUP(?)
Pi DW 3
; ...
```

Initialized variables get their value: `somme DW 0`. Uninitialized variables get `DW ?` (the assembler allocates space but doesn't initialize). Arrays get `DW N DUP(?)` for N elements. Constants get their literal value. All values are 16-bit words (`DW` = Define Word) — even floats are truncated to integers by the `integerize()` function, which parses the string as float64 and formats the integer part.

Compiler-generated temporaries follow user variables:

```asm
    ; compiler-generated temporaries
T0 DW ?
T1 DW ?
; ...
```

They are sorted by numeric suffix (T0, T1, ..., T10, T11 rather than T0, T1, T10, T2, ...) using a bubble sort on the extracted number.

String constants for `out()` calls are stored as DOS $-terminated strings:

```asm
    ; string constants
str0 DB 'Valeur finale de x: ', '$'
str1 DB 'Somme: ', '$'
```

The `$` terminator is the DOS convention for interrupt 21h, function 09h (print string). The `collectStrings()` method in the generator deduplicates strings by content — if `out("Hello")` appears twice, only one `str` label is emitted.

### 6.3 Quad-to-Assembly Mapping

Each quad is emitted at a numbered label (`L0:`, `L1:`, ...) corresponding to its index in the optimized quad list. This label scheme serves two purposes: it maps branch targets (quads reference each other by index) and it makes the generated assembly auditable (you can trace from a quad back to its source).

The mapping from quad ops to assembly is handled by `emitQuad()`, a switch statement dispatching to specialized emitters. Here is the complete mapping:

**Assignment (`:=`):**
```asm
; (:=, 42, , x)        ->  MOV x, 42
; (:=, T0, , y)        ->  MOV AX, T0 / MOV y, AX
```
When the source is an immediate value, a direct `MOV mem, imm` is emitted. When the source is a variable or temp, the value goes through AX.

**Arithmetic (`+`, `-`):**
```asm
; (+, x, 5, T0)        ->  MOV AX, x / ADD AX, 5 / MOV T0, AX
; (-, T0, y, T1)       ->  MOV AX, T0 / MOV CX, y / SUB AX, CX / MOV T1, AX
```
When the second operand is an immediate, it can be used directly in the ADD/SUB instruction. When it's a memory reference, it must go through the CX register first (8086 ADD/SUB with memory destination doesn't support memory source for both operands).

**Multiplication (`*`):**
```asm
; (*, x, y, T0)        ->  MOV AX, x / MOV CX, y / IMUL CX / MOV T0, AX
```
The 8086 `IMUL` instruction always multiplies AX by the operand, placing the 32-bit result in DX:AX. Since ProLang only uses 16-bit values, only AX is stored.

**Division (`/`):**
```asm
; (/, x, 2, T0)        ->  MOV AX, x / MOV CX, 2 / CWD / IDIV CX / MOV T0, AX
```
`CWD` (Convert Word to Doubleword) sign-extends AX into DX:AX, which is required before signed division. `IDIV CX` divides DX:AX by CX, placing the quotient in AX (and the remainder in DX, which is discarded).

**Comparison (`<`, `>`, `<=`, `>=`, `==`, `!=`):**
```asm
; (>, x, 5, T0)
    MOV AX, x
    CMP AX, 5
    MOV AX, 0         ; assume false
    JLE S3_end         ; jump if NOT greater
    MOV AX, 1         ; true
S3_end:
    MOV T0, AX
```
Comparison is the most verbose translation. The pattern sets AX to 0 (false), tests the inverse condition (if we want "greater than", we check "less than or equal"), jumps past the "set to 1" instruction if the inverse is true, and stores the result. The `cmpJump()` function returns the inverse mnemonic:

| Desired | Inverse Jump |
|---------|-------------|
| `>` | `JLE` |
| `<` | `JGE` |
| `>=` | `JL` |
| `<=` | `JG` |
| `==` | `JNE` |
| `!=` | `JE` |

**Control flow (`BF`, `BR`):**
```asm
; (BF, T0, , L5)      ->  MOV AX, T0 / CMP AX, 0 / JZ L5
; (BR, , , L2)        ->  JMP L2
```
Branch-if-false tests whether the condition value is zero. Unconditional branch is a direct `JMP`.

**Array access (`[]`, `[]=`):**
```asm
; ([], Tabint, i, T0)
    MOV SI, i
    ADD SI, SI          ; word scaling: index * 2
    MOV AX, Tabint[SI]
    MOV T0, AX

; ([]=, val, i, Tabint)
    MOV SI, i
    ADD SI, SI
    MOV AX, val
    MOV Tabint[SI], AX
```
Array indexing uses the SI (Source Index) register. The `ADD SI, SI` instruction multiplies the index by 2 because each array element is a word (2 bytes). This is a common 8086 idiom for word-sized array access. The base+index addressing mode (`Tabint[SI]`) adds the array's offset to SI to compute the final address.

**Input (`in`):**
```asm
; (in, , , x)
    MOV AH, 01h
    INT 21h
    SUB AL, '0'
    MOV AH, 0
    MOV x, AX
```
The DOS interrupt 21h function 01h reads a single character from standard input into AL. The character is an ASCII digit, so subtracting '0' (48) converts it to its numeric value. Only single-digit input is supported — a limitation documented in the compiler's design.

**String output (`out_str`):**
```asm
; (out_str, "Resultat: ", , )
    LEA DX, str0
    MOV AH, 09h
    INT 21h
```
DOS function 09h prints a $-terminated string starting at DS:DX. `LEA` (Load Effective Address) computes the address of `str0` and loads it into DX.

### 6.4 The PRINT_INT Runtime Routine

Integer-to-string conversion is the most complex piece of the assembly generator, implemented as a runtime procedure rather than emitted inline:

```asm
PRINT_INT PROC
    PUSH AX / PUSH BX / PUSH CX / PUSH DX    ; save registers

    CMP AX, 0
    JGE print_pos                            ; handle negative
    PUSH AX
    MOV DL, '-' / MOV AH, 02h / INT 21h      ; print '-'
    POP AX / NEG AX                          ; make positive

print_pos:
    MOV CX, 0 / MOV BX, 10                   ; digit counter, divisor
print_div:
    MOV DX, 0 / DIV BX                       ; AX / 10
    PUSH DX / INC CX                         ; push remainder, count it
    CMP AX, 0
    JNE print_div                            ; repeat until quotient is 0

print_digits:
    POP DX                                   ; get digit (reverse order)
    ADD DL, '0'                              ; convert to ASCII
    MOV AH, 02h / INT 21h                    ; print character
    LOOP print_digits                        ; repeat for all digits

    POP DX / POP CX / POP BX / POP AX        ; restore registers
    RET
PRINT_INT ENDP
```

The algorithm is the classic "divide by 10, push remainders, pop and print" approach:

1. Save all modified registers (AX, BX, CX, DX) so the caller doesn't see side effects.
2. Check for negative values. If negative, print a minus sign and negate the value.
3. Repeatedly divide by 10, pushing the remainder (the current least-significant digit) onto the stack and counting digits in CX.
4. Pop digits off the stack (which reverses the order — most significant digit comes off first), convert to ASCII by adding '0', and print via DOS function 02h.
5. `LOOP` decrements CX and jumps if not zero, printing exactly the right number of digits.
6. Restore saved registers and return.

This routine handles values from -32768 to 32767 (the range of a signed 16-bit integer). Zero is handled correctly: the division loop pushes one zero digit, and the print loop prints it.

### 6.5 The DOS Interrupt Interface

All I/O in the generated program goes through DOS interrupts (INT 21h). The functions used are:

| AH | Function | Input | Output |
|----|----------|-------|--------|
| 01h | Read character with echo | — | AL = ASCII character |
| 02h | Write character | DL = ASCII character | — |
| 09h | Write string | DS:DX = address of $-terminated string | — |
| 4Ch | Terminate program | AL = return code | (program exits) |

This is the standard DOS programming interface circa 1981. It works on any DOS-compatible system (MS-DOS, FreeDOS, DOSBox). The `$-terminated` string convention (as opposed to C's null-terminated strings) is a DOS-ism: the `$` character was chosen because it's rarely needed in output text, and it avoids the need for a length parameter.

### 6.6 Limitations and Design Trade-offs

The code generator makes several deliberate trade-offs in favor of simplicity:

**Integer-only arithmetic.** All floating-point values are truncated to integers via `integerize()`. The language has a `float` type and the semantic analyzer tracks float vs. integer types, but the backend has no x87 FPU code. Adding floating-point support would require: FPU initialization, `FLD`/`FSTP` instructions for load/store, `FADD`/`FSUB`/`FMUL`/`FDIV` for arithmetic, `FCOM` for comparison, and a float-to-string routine for printing. This would roughly double the size of the assembly generator. For a pedagogical compiler, the integer restriction is acceptable — it demonstrates the pipeline without the complexity of the x87 programming model.

**Single-digit input.** The `in` operation reads one character and subtracts '0', supporting only values 0–9. Multi-digit input would require a string-to-integer conversion routine similar to PRINT_INT in reverse. This is left as an exercise for students extending the compiler.

**No register allocation.** Every quad result is stored to a memory temporary (`T0 DW ?`), and every operand is loaded from memory into a register. A real compiler would perform register allocation to keep frequently-used values in registers, using graph coloring or linear scan. The ProLang compiler uses a trivial "AX for computation, CX for second operand, SI for indexing" scheme that produces correct but bloated code. Teaching register allocation is a natural next step after understanding this codebase.

**No peephole optimization.** The assembly generator emits each quad independently. Adjacent instructions like `MOV T0, AX` followed by `MOV AX, T0` (a redundant round-trip through memory) are not eliminated. A peephole pass over the emitted assembly could remove these, but the optimizer's dead code elimination catches many of them at the quad level.

**MASM/TASM compatibility.** The generated assembly uses syntax accepted by both Microsoft Macro Assembler (MASM) and Borland Turbo Assembler (TASM). The segment directives, `ASSUME`, `DUP(?)`, and `$`-terminated strings are standard across both assemblers. Some minor differences (like `IDEAL` vs `MASM` mode in TASM) do not affect this code.

### 6.7 Running the Generated Code

To execute the generated assembly:

1. **Assemble:** `tasm final.pl.asm` or `masm final.pl.asm;`
2. **Link:** `tlink final.pl.obj` or `link final.pl.obj;`
3. **Run:** Execute `final.pl.exe` in DOSBox or on actual DOS hardware.

In DOSBox, the debugger can be activated with the `debug=true` configuration option or by pressing Ctrl+F11 at runtime. This opens a CPU view showing registers, memory, and the current instruction — the window the user sees in YouTube tutorials. Breakpoints can be set at labels like `L47` or `L90` to inspect loop conditions and output values.

---

## Appendix A: Project Structure

```
prolang-compiler/
├── cmd/
│   └── main.go                  # Entry point; orchestrates the pipeline
├── internal/
│   ├── lexer/
│   │   ├── token.go             # Token kind constants and Token struct
│   │   └── lexer.go             # Lexer: []rune-based scanner
│   ├── ast/
│   │   ├── ast.go               # AST node type definitions
│   │   └── symbol.go            # Symbol table with FNV-1a hashing
│   ├── semantic/
│   │   ├── analyzer.go          # Recursive-descent parser + semantic analysis
│   │   └── quad_generator.go    # Token stream → quad IR pass
│   ├── quad/
│   │   └── quad.go              # Quad struct, Manager, backpatching
│   ├── optimizer/
│   │   └── optimizer.go         # 5-pass optimizer
│   ├── asm/
│   │   └── generator.go         # Quad IR → 8086 assembly emitter
│   └── utils/
│       └── tui.go               # Terminal styling (lipgloss)
├── example/
│   ├── final.pl                 # Comprehensive test program
│   ├── test_loop.pl             # Simple while loop test
│   ├── test_ifelse.pl           # If/else branching test
│   ├── test_nested_loop.pl      # Nested while loops test
│   ├── test_array.pl            # Array read/write test
│   └── test_mixed_expr.pl       # Mixed expression test
├── out/                         # Generated .asm files
├── docs/
│   └── prolang-compiler-documentation.md  # This document
├── go.mod
├── go.sum
└── main                         # Compiled binary (go build)
```

## Appendix B: Key Design Decisions Summary

| Decision | Rationale |
|----------|-----------|
| Go implementation language | Fast compilation, native binaries, GC for memory management, readable syntax for pedagogy |
| Hand-written lexer (not generated) | Full control over error recovery, simpler build process, ~300 lines is manageable |
| Parser + semantic analysis combined | Grammar is small enough (LL(1)) that separation adds complexity without clarity |
| FNV-1a hash for symbol table | Better avalanche than DJB2, faster than cryptographic hashes, standard library available |
| Fixed 1024-bucket hash table | Simplifies debugging (no rehashing), ample for expected program size |
| Three-address code (quads) as IR | Classic representation, easy to print/debug, enables optimizer to be language-agnostic |
| 5-pass local optimizer | Catches common cases without requiring control-flow graph construction |
| SSA-like constraint on propagation | Prevents incorrect constant propagation across reassignments without full data-flow analysis |
| CSE limited to temps/literals | Sound without alias/definition analysis; single-assign temps are safe to cache |
| 8086 real-mode target | Simpler than x86-64, debuggable in DOSBox/Turbo Debugger, pedagogically appropriate |
| MASM/TASM-compatible syntax | Works with widely available DOS tooling |
| Truncation of floats to integers | Avoids x87 FPU complexity; acceptable for a pedagogical compiler |
| Runtime PRINT_INT procedure | Avoids duplicating digit-conversion code at every output site |
