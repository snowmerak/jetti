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

#### generate

```bash
jetti bean --generate
```

Then the code below will be generated in `gen/bean/bean.go`.

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

### new

#### init

`init` is a simple command to initialize a project.

##### Usage

```bash
jetti new --init <project-name>
```

Then, a project with the following structure will be created.

```bash
.
├── README.md
├── cmd
├── lib
├── internal
├── proto
├── configs
├── uml
├── go.mod
├── go.sum
```

### proto

`proto` is a simple command to generate protobuf and gRPC code from proto files.

#### new

`new` is a simple command to initialize a proto project.

```bash
jetti proto --new <path>
```

Then, a proto file will be created in the specified path.

If, for example, you run `jetti proto --new person/person.proto`, the file will be created in `proto/person/person.proto`.

#### build

```bash
jetti proto --generate
```

Then, jetti is going to generate protobuf and gRPC code from proto files in `proto` directory.

The generated code will be created in `gen/proto`.

### cmd

`cmd` is a simple command to manage executable package.

#### new

`new` is a simple command to initialize a cmd project.

```bash
jetti cmd --new <name>
```

Then, a cmd file will be created in `cmd/<name>/main.go`.

#### build

```bash
jetti cmd --build=<name>,<option1>,<option2>,...
```

Then, jetti is going to build executable package from cmd files in `cmd` directory.

The executable file will be created in `bin` directory.

#### run

```bash
jetti cmd --run=<name>,<arg1>,<arg2>,...
```

Then, jetti is going to run executable package from cmd files in `cmd` directory.

### pprof

`pprof` is a simple command to make pprof server.

#### http1

```bash
jetti pprof --http-1 <addr>
```

Then, jetti is going to make pprof server with http1 in `gen/pprof/http1`.

#### http2

```bash
jetti pprof --http-2 <addr>
```

Then, jetti is going to make pprof server with http2 in `gen/pprof/http2`.

#### http3

```bash
jetti pprof --http-3 <addr>
```

Then, jetti is going to make pprof server with http3 in `gen/pprof/http3`.

### redis

`redis` is a simple command to make some types to using redis.

#### new

```bash
jetti redis --new <path-name>
```

Then, jetti is going to make some types to using redis in `template/redis/<path-name>`.

Example, if you run `jetti redis --new person`, jetti is going to make a file in `template/redis/person`.

#### directions

##### string

```go
//go:redis string <name> <key>
```

Jetti will generate the redis client for string type.

#### string[generic]

```go
//go:redis string[<path>/<package>.<name>] <key>
```

Jetti will generate the redis client for string type with generic.

> You must use protocol buffers type for generic.

If you want to use the `Person` protobuf in `gen/model/person`, you can use it like this.

```go
//go:redis string[gen/model/person.Person] <key>
```

#### list

```go
//go:redis list <name> <key>

//go:redis list[<path>/<package>.<name>] <key>
```

#### set

```go
//go:redis set <name> <key>

//go:redis set[<path>/<package>.<name>] <key>
```

#### bitmap

```go
//go:redis bitmap <name> <key>
```
