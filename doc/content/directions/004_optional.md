---
title: 'Optional'
date: 2023-05-16T21:12:23+09:00
---

`Jetti` can generate a `optional` type with `optional` direction.

## Example

```go
package person

// jetti:optional
type Person struct {
	Name string
	Age  int
}

// jetti:optional
type People [100]Person
```

## Generated Code

Run `jetti generate` command, `jetti` generates the following code.

### normal

```go
package person

type OptionalPeople struct {
	value *People
	valid bool
}

func (o *OptionalPeople) Unwrap() *People {
	if !o.valid {
		panic("unwrap a none value")
	}
	return o.value
}

func (o *OptionalPeople) IsSome() bool {
	return o.valid
}

func (o *OptionalPeople) IsNone() bool {
	return !o.valid
}

func (o *OptionalPeople) UnwrapOr(defaultValue *People) *People {
	if !o.valid {
		return defaultValue
	}
	return o.value
}

func SomePeople(value *People) OptionalPeople {
	return OptionalPeople{
		value: value,
		valid: true,
	}
}

func NonePeople() OptionalPeople {
	return OptionalPeople{
		valid: false,
	}
}
```

`people`'s optional type is `OptionalPeople`.

We can use `SomePeople` and `NonePeople` to create `OptionalPeople` type.

### same name with package

```go
package person

type OptionalPerson struct {
	value *Person
	valid bool
}

func (o *OptionalPerson) Unwrap() *Person {
	if !o.valid {
		panic("unwrap a none value")
	}
	return o.value
}

func (o *OptionalPerson) IsSome() bool {
	return o.valid
}

func (o *OptionalPerson) IsNone() bool {
	return !o.valid
}

func (o *OptionalPerson) UnwrapOr(defaultValue *Person) *Person {
	if !o.valid {
		return defaultValue
	}
	return o.value
}

func Some(value *Person) OptionalPerson {
	return OptionalPerson{
		value: value,
		valid: true,
	}
}

func None() OptionalPerson {
	return OptionalPerson{
		valid: false,
	}
}
```

`person`'s optional type is `OptionalPerson`.

We can use `Some` and `None` to create `OptionalPerson` type.
