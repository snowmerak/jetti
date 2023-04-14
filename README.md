# jetti

## Installation

Install `jetti` with go install.

```bash
go install github.com/snowmerak/jetti/cmd/jetti@latest
```

## commands

### bean

bean is a simple single tone bean container for golang.

#### Usage

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
jetti -bean
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

### init

`init` is a simple command to initialize a project.

#### Usage

```bash
jetti -init <project-name>
```

Then, a project with the following structure will be created.

```bash
.
├── README.md
├── cmd
├── src
├── internal
├── proto
├── configs
├── uml
├── go.mod
├── go.sum
```

And, add the dependencies to `go.mod`.

```bash
go get github.com/goccy/go-json
go get github.com/goccy/go-yaml
go get google.golang.org/protobuf
go get google.golang.org/grpc
```

### proto

`proto` is a simple command to generate protobuf and gRPC code from proto files.

#### Usage

```bash
jetti -proto
```

Then, jetti is going to generate protobuf and gRPC code from proto files in `proto` directory.

### proto-make

`proto-make` is a simple command to generate protobuf and gRPC code from proto files.

#### Usage

```bash
jetti -proto-make <path-to-proto-file>
```

Then, jetti is going to generate protobuf and gRPC code from proto files in `proto` directory.

For example, `jetti -proto-make person/person.proto` will generate protobuf and gRPC code to `proto/person/person.proto`.
