# JellySet: Redis Compatible Set Data Structure in Go

## Overview

`jellyset` is a highly-efficient, straightforward and easy-to-integrate Go package that provides a Redis-like set data structure. It is designed to be simple to use and understand, offering essential set operations with an emphasis on clarity and efficiency. The package includes a high-level interface, Set, which manages multiple sets, each associated with a unique key.

It is onspired by the simplicity of Redis [SET](https://redis.io/docs/data-types/sets/) and influenced by the [arriqaaq/set](https://github.com/arriqaaq/set) project.

## Instalation

To integrate jellyset into your Go project, execute the following command:
```bash
go get -u github.com/davidandw190/jellyset
```

## Usage

### Creating a new Set

```go
import "github.com/your-username/jellyset"

// Create a new set instance
mySet := jellyset.New()
```

### Operations

```go
// Add members to the set
count := mySet.SAdd("mySet", "member1", "member2", "member3")

// Remove and return random members from the set
popped := mySet.SPop("mySet", 3)

// Return random members from the set without removal
randomMembers := mySet.SRandMember("mySet", 3)

// Check if a member exists in the set
exists := mySet.SIsMember("mySet", "member2")

// Remove a member from the set
removed := mySet.SRem("mySet", "member2")

// Move a member from one set to another
moved := mySet.SMove("sourceSet", "destSet", "member2")

// Get the number of elements in the set
size := mySet.SCard("mySet")

// Get all members of the set
members := mySet.SMembers("mySet")

// Get the union of multiple sets
unionResult := mySet.SUnion("set1", "set2")

// Store the union of multiple sets in a new set
unionCount := mySet.SUnionStore("unionSet", "set1", "set2")

// Check if a key exists in the set
keyExists := mySet.SKeyExists("mySet")

// Clear a set
mySet.SClear("mySet")

// Get the difference between two sets
differenceResult := mySet.SDiff("set1", "set2")

// Store the difference between two sets in a new set
differenceCount := mySet.SDiffStore("differenceSet", "set1", "set2")

// Get the intersection of multiple sets
intersectionResult := mySet.SInter("set1", "set2")

// Store the intersection of multiple sets in a new set
intersectionCount := mySet.SInterStore("intersectionSet", "set1", "set2")
```

### Implementation Details

`NOTE:` It's important to note that, at the moment, `jellyset` is not designed to be concurrent.