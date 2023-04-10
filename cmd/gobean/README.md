# go bean

go bean is a simple single tone bean container for golang.

## Installation

Install `gobean` with go install.

```bash
go install github.com/snowmerak/go-bean/cmd/gobean@latest
```

## Usage

Write `//go:bean` above the struct you want to register to the bean container.

```go
package person

//go:bean
type Person struct {
    Name string
    Age  int
}
````

Run `gobean` in the root directory of your project.

```bash
gobean
```

Then the code below will be generated.

```go
package bean

type Bean struct {
	person *person.Person
}

type Builder struct {
	bean *Bean
}

func New() *Builder {
	return &Builder{}
}

func (b *Builder) Build() *Bean {
	return b.bean
}

func (b *Builder) Person(person *person.Person) *Builder {
	b.bean.person = person
	return b
}

func (b *Bean) Person() *person.Person {
	return b.person
}
```

We do not import dependency package automatically, so you need to import it manually.

Just save it.
