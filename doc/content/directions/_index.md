---
title: 'Direction'
date: 2023-05-16T21:12:23+09:00
weight: 2
---

`Direction` specifies the code that `jetti` should generate.

## Usage

```go
package person

// jetti:bean
type Person struct {
	Name string
	Age int
}
```

`jetti:bean` is an instruction for `jetti` to create a bean container based on the Person structure.

Typically, you prefix `jetti:` and follow it with directives like `bean`.

I will continue to provide other directives and their usage.
