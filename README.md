# go-critbit

A high-performance Go implementation of Crit-bit trees (also known as PATRICIA trees) with support for arbitrary binary keys and generic values.


## Overview

Crit-bit trees are space-efficient binary trees that store keys in a way that eliminates unnecessary internal nodes. Each internal node represents the first critical bit where two subtrees diverge. This implementation provides O(k) operations where k is the key length, making it excellent for applications requiring fast lookups with binary keys.


## Features

- **Generic Values**: Support for any value type using Go generics (`Tree[V any]`)
- **Flexible Key Types**: Built-in support for strings, integers (uint64/32/16/8), and arbitrary byte sequences
- **Bit-Level Operations**: Fine-grained control with bit-level key manipulation
- **Memory Efficient**: Compact tree structure that eliminates redundant nodes
- **Ordered Iteration**: Natural lexicographical ordering with forward/reverse traversal
- **Longest Prefix Matching**: Built-in LPM support for routing and networking applications
- **Modern Go**: Uses Go 1.23+ features including `iter.Seq` for iteration


## Quick Start

```go
package main

import (
    "fmt"

    "github.com/koji-hirono/go-critbit"
)

func main() {
    // Create a new critbit tree for string values
    var tree critbit.Tree[string]

    // Insert key-value pairs using same key types
    tree.Set(critbit.StringKey("apple"), "fruit")
    tree.Set(critbit.StringKey("application"), "software")
    tree.Set(critbit.StringKey("ant"), "insect")

    // Get values
    if value, found := tree.Get(critbit.StringKey("apple")); found {
        fmt.Printf("Found: %s\n", value) // Found: fruit
    }

    // Check length
    fmt.Printf("Tree contains %d items\n", tree.Len())

    // Iterate over all key-value pairs
    for key, value := range tree.All() {
        fmt.Printf("Key: %v, Value: %s\n", key, value)
    }
}
```

## Key Types

The library provides several convenience functions for creating keys:

```go
// String keys
key := critbit.StringKey("hello world")

// Integer keys (stored in big-endian format)
key := critbit.Uint64Key(1234567890)
key := critbit.Uint32Key(12345)
key := critbit.Uint16Key(123)
key := critbit.Uint8Key(42)

// Byte slice keys
key := critbit.BytesKey([]byte("raw bytes"))

// Bit-level keys (custom bit length)
key := critbit.BitsKey([]byte{0b10110000}, 5) // Only first 5 bits
```

## API Reference

### Tree Operations

```go
// Create a new tree (zero value is ready to use)
var tree critbit.Tree[ValueType]

// Set inserts or updates a key-value pair
func (t *Tree[V]) Set(key Key, value V)

// Get retrieves a value by key
func (t *Tree[V]) Get(key Key) (V, bool)

// Delete removes a key from the tree
func (t *Tree[V]) Delete(key Key)

// Len returns the number of items in the tree
func (t *Tree[V]) Len() int
```

### Longest Prefix Matching

```go
// Find the longest prefix match for a given key
func (t *Tree[V]) Longest(key Key) (V, bool)
```

### Iteration

The library provides Go 1.23+ iterator support:

```go
// Iterate over all key-value pairs
func (t *Tree[V]) All() iter.Seq2[Key, V]

// Iterate over keys only
func (t *Tree[V]) Keys() iter.Seq[Key]

// Iterate over values only
func (t *Tree[V]) Values() iter.Seq[V]
```


## Examples

### Longest Prefix Matching (LPM)

Perfect for IP routing tables and network applications:

```go
import "net/netip"

var routingTable critbit.Tree[string]

// Add routes with their network prefixes
addRoute := func(cidr string, gateway string) {
    prefix := netip.MustParsePrefix(cidr)
    key := critbit.Key{
        Data:  prefix.Addr().AsSlice(),
        Nbits: prefix.Bits(),
    }
    routingTable.Set(key, gateway)
}

addRoute("10.0.0.0/8", "gateway-1")
addRoute("10.1.0.0/16", "gateway-2")
addRoute("10.1.2.0/24", "gateway-3")

// Find best matching route
lookupRoute := func(ip string) {
    addr := netip.MustParseAddr(ip)
    key := critbit.Key{
        Data:  addr.AsSlice(),
        Nbits: len(addr.AsSlice()) * 8,
    }

    if gateway, found := routingTable.Longest(key); found {
        fmt.Printf("Route for %s: %s\n", ip, gateway)
    }
}

lookupRoute("10.1.2.100") // Route for 10.1.2.100: gateway-3
lookupRoute("10.1.5.1")   // Route for 10.1.5.1: gateway-2
lookupRoute("10.5.0.1")   // Route for 10.5.0.1: gateway-1
```

### Iteration and Traversal

```go
var tree critbit.Tree[string]
tree.Set(critbit.StringKey("cat"), "animal")
tree.Set(critbit.StringKey("car"), "vehicle")
tree.Set(critbit.StringKey("card"), "payment")

// Forward iteration (lexicographical order)
fmt.Println("Forward iteration:")
for key, value := range tree.All() {
    // Convert key back to string for display
    keyStr := string(key.Data[:key.Nbits/8])
    fmt.Printf("%s: %s\n", keyStr, value)
}
// Output: car: vehicle, card: payment, cat: animal


### Working with Bit-Level Keys

```go
var tree critbit.Tree[string]

// Insert keys with specific bit lengths
tree.Set(critbit.BitsKey([]byte{0b10000000}, 1), "prefix-1")  // Just "1"
tree.Set(critbit.BitsKey([]byte{0b10100000}, 3), "prefix-101") // "101"
tree.Set(critbit.BitsKey([]byte{0b10110000}, 4), "prefix-1011") // "1011"

// The tree will organize these by their bit patterns
for key, value := range tree.All() {
    fmt.Printf("Bits: %d, Value: %s\n", key.Nbits, value)
}
```


## Use Cases

This critbit tree implementation is particularly well-suited for:

- **IP Routing Tables**: Longest prefix matching for network routing
- **Autocomplete Systems**: Prefix-based text search and suggestions
- **Binary Protocol Parsing**: Efficient lookup of bit patterns
- **Compressed Data Structures**: Space-efficient key storage
- **Ordered Key-Value Storage**: When you need both fast lookup and ordering
- **Networking Applications**: Any application requiring LPM functionality


## Thread Safety

This implementation is **not thread-safe**. For concurrent access, use external synchronization:

```go
import "sync"

var (
    tree critbit.Tree[string]
    mu   sync.RWMutex
)

// Safe read
func safeGet(key critbit.Key) (string, bool) {
    mu.RLock()
    defer mu.RUnlock()
    return tree.Get(key)
}

// Safe write
func safeSet(key critbit.Key, value string) {
    mu.Lock()
    defer mu.Unlock()
    tree.Set(key, value)
}
```


## Implementation Details

- Uses Go generics for type-safe value storage
- Implements efficient bit manipulation for key operations
- Supports variable-length bit keys (not just byte-aligned)
- Uses Go 1.23+ iterators for modern iteration patterns
- Optimized for both memory usage and speed


## License

This software is released under the MIT License. See [LICENSE](LICENSE) for details.


## References

- [Crit-bit trees](https://cr.yp.to/critbit.html) by D. J. Bernstein
- [PATRICIA Practical Algorithm to Retrieve Information Coded in Alphanumeric](https://dl.acm.org/doi/10.1145/321479.321481)
- [Go 1.23 Iterator Patterns](https://go.dev/blog/range-functions)
