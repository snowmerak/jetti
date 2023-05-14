# jetti

제티(jetti)는 고 언어 프로젝트에 어느정도 정규화된 프로젝트 구성을 적용하기 위한 도구입니다.

## 지향점

[WIP]

## new

`new`는 새로운 프로젝트나 커맨드 패키지를 생성합니다.

### new module

`jetti new <module-name>`을 실행하면 현재 폴더에서 `go mod <module-name>`을 실행하면서 다음과 같은 기본 폴더들을 만들어줍니다.

```
➜  tree .
.
├── README.md
├── cmd
│   └── doc.go
├── go.mod
├── internal
│   └── doc.go
└── lib
    └── doc.go

4 directories, 5 files
```

### new command

`jetti new --cmd <cmd-name>`을 실행하면 현재 폴더 내의 `cmd` 폴더에 `<cmd-name>` 폴더를 만들고, `main.go` 파일을 만들어줍니다.

다음 예시는 `jetti new --cmd prac`를 실행한 결과입니다.

```
➜  jetti new --cmd prac
➜  tree .
.
├── README.md
├── cmd
│   ├── doc.go
│   └── prac
│       └── main.go
├── go.mod
├── internal
│   └── doc.go
└── lib
    └── doc.go

5 directories, 6 files
```

### new proto

`jetti new --proto <path/name>`을 실행하면 현재 폴더 내의 `<path>` 폴더를 만들고, `<name>.proto` 파일을 만듭니다.

다음 예시는 `jetti new --proto model/proto/person`를 실행한 결과입니다.

```protobuf
syntax = "proto3";

package person;

option go_package = "model/proto/person";

```

## run

`run`은 `cmd` 내의 커맨드 패키지를 실행하는 역할을 합니다.

`jetti run <cmd-name>`을 실행하면 `cmd/<cmd-name>` 폴더 내의 고 파일들을 실행합니다.

추가로 `jetti run <cmd-name> <args>...`을 실행하여 커맨드 패키지에 인자를 전달할 수 있습니다.  
사실상 `go run`과 동일합니다.

## bean

`bean`은 `context.Context`의 `WithValue`를 편리하게 이용할 수 있게 해주는 기능입니다.

빈의 핵심은 동등, 혹은 하위 문맥에서 동일한 객체 인스턴스를 공유하는 것입니다.

### 예시

`./lib/config` 폴더를 만들고 `config.go` 파일을 만들어 다음과 같이 작성합니다.

```go
package config

// jetti:bean redis postgres
type Config struct {
}
```

`jetti:bean` 주석을 통해 `redis`와 `postgres` 빈을 등록했습니다.

이제 터미널에 `jetti generate`를 입력하면 `./lib/redis.context.go`와 `./lib/postgres.context.go` 파일이 생성됩니다.

```go
// postgres.context.go
package config

import "context"

type PostgresContextKey string

func PushPostgres(ctx context.Context, v *Config) context.Context {
	return context.WithValue(ctx, PostgresContextKey("Postgres"), v)
}

func GetPostgres(ctx context.Context) (*Config, bool) {
	v, ok := ctx.Value(PostgresContextKey("Postgres")).(*Config)
	return v, ok
}
```

```go
package config

import "context"

type RedisContextKey string

func PushRedis(ctx context.Context, v *Config) context.Context {
	return context.WithValue(ctx, RedisContextKey("Redis"), v)
}

func GetRedis(ctx context.Context) (*Config, bool) {
	v, ok := ctx.Value(RedisContextKey("Redis")).(*Config)
	return v, ok
}
```

이제 단일 컨텍스트를 생성한 후, `Push` 메서드를 통해 빈을 등록하고, `Get` 메서드를 통해 빈을 가져올 수 있습니다.

## optional parameter

`optional parameter`는 `jetti:optional` 주석을 통해 생성할 수 있습니다.

옵셔널 패러미터는 기존의 프리미티브 타입, 혹은 구조체에 기본값과 값 변경을 위한 함수를 받아 기본값을 변형하여 새로운 패러미터를 반환합니다.

### 예시

`./lib/person` 폴더를 만들고 `person.go` 파일을 만들어 다음과 같이 작성합니다.

```go
package person

// jetti:optional
type Person struct {
	Name string
	Age  int
}
```

`jetti generate`를 실행하면 `./lib/person.optional.go` 파일이 생성됩니다.

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

`ApplyPerson` 함수에 기본값과 변형 함수를 전달하여 새로운 `Person` 구조체를 생성합니다.

## json/yaml to go

`go-jsonstruct` 라이브러리를 이용해서 json/yaml 파일을 go 구조체로 변환할 수 있습니다.

파싱에는 각각 `goccy/go-json`과 `goccy/go-yaml` 라이브러리를 사용합니다.

### 예시

`./config/json_prac.json`과 `./config/yaml_prac.yaml` 파일을 만들고 다음과 같이 작성합니다.

```json
{
  "name": "snowmerak",
  "version": "1.3.2",
  "author": "snowmerak",
  "dependencies": {
    "go": "github.com/golang/go",
    "rust": "github.com/rust-lang/rust"
  }
}
```

```yaml
name: snowmerak
version: 1.3.2
author: snowmerak
dependencies:
  go: github.com/golang/go
  rust: github.com/rust-lang/rust
```

그리고 `jetti generate`를 실행하면 다음 파일 들이 생성됩니다.

```go
// json_prac.json.go
package config

import "github.com/goccy/go-json"
import "io"
import "os"

func JsonPracFromJSON(data []byte) (*JsonPrac, error) {
	v := new(JsonPrac)
	if err := json.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v, nil
}

func JsonPracFromFile(path string) (*JsonPrac, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return JsonPracFromJSON(f)
}

func (jsonprac *JsonPrac) Marshal2JSON() ([]byte, error) {
	return json.Marshal(jsonprac)
}

func (jsonprac *JsonPrac) Encode2JSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(jsonprac)
}

type JsonPrac struct {
	Author       string `json:"author"`
	Dependencies struct {
		Go   string `json:"go"`
		Rust string `json:"rust"`
	} `json:"dependencies"`
	Name    string `json:"name"`
	Version string `json:"version"`
}
```

```go
// yaml_prac.yaml.go
package config

import "github.com/goccy/go-yaml"
import "io"
import "os"

func YamlPracFromYAML(data []byte) (*YamlPrac, error) {
	v := new(YamlPrac)
	if err := yaml.Unmarshal(data, v); err != nil {
		return nil, err
	}
	return v, nil
}

func YamlPracFromFile(path string) (*YamlPrac, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return YamlPracFromYAML(f)
}

func (yamlprac *YamlPrac) Marshal2YAML() ([]byte, error) {
	return yaml.Marshal(yamlprac)
}

func (yamlprac *YamlPrac) Encode2YAML(w io.Writer) error {
	return yaml.NewEncoder(w).Encode(yamlprac)
}

type YamlPrac struct {
	Author       string `yaml:"author"`
	Dependencies struct {
		Go   string `yaml:"go"`
		Rust string `yaml:"rust"`
	} `yaml:"dependencies"`
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}
```

## protobuf/flatbuffers generating

제티는 프로젝트 내부의 프로토버퍼와 플랫버퍼 파일을 찾으면 자동으로 고 코드로 컴파일 하는 커맨드를 실행합니다.

이를 위해 `protoc`와 `flatc`가 필요합니다.

### protoc/grpc 설치

[여기]("https://grpc.io/docs/languages/go/quickstart/")를 참고해 protoc 및 고 코드 생성을 위한 플러그인을 설치합니다.

### flatc 설치

[여기]("https://google.github.io/flatbuffers/flatbuffers_guide_building.html")를 참고해 flatc를 설치합니다.

굳이 빌드 하지 않더라도 사용하는 환경의 패키지 매니저를 통해 설치할 수 있습니다.

### 사용법

`./proto` 디렉토리를 만들고 `./proto/test/test.proto` 파일을 만듭니다.

```proto
syntax = "proto3";

package test;

option go_package = "./test";

message Test {
  string name = 1;
  int32 age = 2;
}
```

그리고 `jetti generate`를 실행하면 `./gen/grpc/proto/test/test.pb.go`에 파일을 생성합니다.

## object pool

제티는 `sync.Pool`과 `chan T`을 사용한 두가지 풀을 만들 수 있습니다.

### sync pool

`jetti:pool`을 주석에 작성함으로 풀을 생성할 수 있습니다.

두 가지 풀 중, sync.Pool은 `jetti:pool sync:<alias>`로 생성할 수 있습니다.

`<alias>`는 풀의 이름을 지정합니다.

```go
// jetti:pool sync:people
type Person struct {
    Name string
    Age  int
}
```

위와 같이 주석을 작성하면 `sync.Pool`을 이용한 `people` 풀이 생성됩니다.

```go
package person

import "sync"
import "errors"
import "runtime"

var errPeopleCannotGet error = errors.New("cannot get people")

type PeoplePool struct {
	pool *sync.Pool
}

func (p *PeoplePool) Get() (*Person, error) {
	v := p.pool.Get()
	if v == nil {
		return nil, errPeopleCannotGet
	}
	return v.(*Person), nil
}

func (p *PeoplePool) GetWithFinalizer() (*Person, error) {
	v := p.pool.Get()
	if v == nil {
		return nil, errPeopleCannotGet
	}
	runtime.SetFinalizer(v, func(v interface{}) {
		p.pool.Put(v)
	})
	return v.(*Person), nil
}

func (p *PeoplePool) Put(v *Person) {
	p.pool.Put(v)
}

func NewPeoplePool() PeoplePool {
	return PeoplePool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(Person)
			},
		},
	}
}

func IsPeopleCannotGetErr(err error) bool {
	return errors.Is(err, errPeopleCannotGet)
}
```

### chan pool

채널을 사용한 풀은 `jetti:pool chan:<alias>`로 생성할 수 있습니다.

이 풀의 경우엔, 최대 풀링 가능한 오브젝트 수를 제한할 때 유용하게 사용할 수 있습니다.

```go
// jetti:pool sync:people chan:candidate
type Person struct {
	Name string
	Age  int
}
```

방금 예제에서 `sync:people` 뒤에 `chan:candidate`를 추가해서 생성하면, 추가적으로 다음 파일도 생성됩니다.

```go
package person

import (
	"runtime"
	"time"
)

type CandidatePool struct {
	pool    chan *Person
	timeout time.Duration
}

func (c *CandidatePool) Get() *Person {
	after := time.After(c.timeout)
	select {
	case v := <-c.pool:
		return v
	case <-after:
		return new(Person)
	}
}

func (c *CandidatePool) GetWithFinalizer() *Person {
	after := time.After(c.timeout)
	resp := (*Person)(nil)
	select {
	case v := <-c.pool:
		resp = v
	case <-after:
		resp = new(Person)
	}
	runtime.SetFinalizer(resp, func(v interface{}) {
		c.pool <- v.(*Person)
	})
	return resp
}

func (c *CandidatePool) Put(v *Person) {
	select {
	case c.pool <- v:
	default:
	}
}

func NewCandidatePool(size int, timeout time.Duration) CandidatePool {
	pool := make(chan *Person, size)
	return CandidatePool{
		pool:    pool,
		timeout: timeout,
	}
}
```

sync pool과 다른 점으로 전체 채널 길이와 채널에서 값을 가져올 시간의 제한을 지정합니다.

## show

### imports

프로젝트 루트 폴더에서 `jetti show --imports`을 실행함으로 프로젝트 내부에서 각 패키지들이 의존하는 관계를 그래프로 그려줍니다.

그려진 그래프는 루트 폴더 내의 `imports.svg` 파일로 저장됩니다.
