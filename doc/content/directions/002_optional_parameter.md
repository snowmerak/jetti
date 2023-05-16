---
title: 'Optional Parameter'
date: 2023-05-16T21:12:23+09:00
---

Optional `parameter` is a directive that make boilerplate code for optional parameter.

## Example

We can use `parameter` directive like this.

```go
package person

// jetti:parameter
type Person struct {
	Name string
	Age  int
}
```

## Generated Code

Run `jetti generate` command, `jetti` generates the following code.

```go
package person

type PersonOptional func(*Person) *Person

func ApplyPerson(defaultValue Person, fn ...PersonOptional) *Person {
	param := &defaultValue
	for _, f := range fn {
		param = f(param)
	}
	return param
}
```

`ApplyPerson` function is a function that applies optional parameters to the default value.
