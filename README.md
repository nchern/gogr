gogr - Go Grep
==============

`gogr` makes it easier to perform sophisticated queries to Go sources using `grep`.

## Examples

```bash
# find all methods
find /path/to/module -name '*.go' | gogr | grep ':method:'

# find all structs that have a field X of type float64 
find /path/to/module -name '*.go' | gogr | grep ':struct:' | grep '.X float64'

# find all structs that have a method with a given signature: 'Foo() string'
find /path/to/module -name '*.go' | gogr | grep ':method:' | grep '.Foo() string'

# find all calls to log that write a given substring
find /path/to/module -name '*.go' | gogr | grep ':call:' | grep 'log.Printf' | grep -E '".*processor.*"'
```

## How it works
The utility parses Go sources and transforms them in the following way:
- multiline calls turned into a single line
- multiline function declarations turned into a single line
- multiline conditions in `if` / `for` statements turned into a single line
- members of interfaces and structures are enriched with their type name and printed one member per line
