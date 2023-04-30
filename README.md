# jetti

제티(jetti)는 고 언어 프로젝트에 어느정도 정규화된 프로젝트 구성을 적용하기 위한 도구입니다.

## 지향점

### 프로젝트 구성

제티가 지향하는 프로젝트 구성은 다음과 같습니다.

1. cmd: 폴더의 하위에 각각의 프로젝트의 엔트리 포인트가 될 패키지들이 존재합니다.
2. lib: 프로젝트 내의 공통 라이브러리들이 존재합니다. 보통 외부에 공개해도 되는 코드들이 포함됩니다.
   1. client: 프로젝트 내에서 사용되는 클라이언트 라이브러리들이 존재합니다. 레디스, 나츠, sql, gRPC 클라이언트 등이 존재합니다.
   2. server: 프로젝트 내에서 사용되는 다른 프로세스나 클라이언트의 요청을 받아줄 서비스 라이브러리들이 존재합니다. gRPC 서버, HTTP 서버 등이 존재합니다.
   3. worker: 프로젝트 내에서 사용되는 백그라운드 워커 라이브러리들이 존재합니다. 주로 주기적으로 동작하는 코드나, 프로세스 내에서 끊임 없이 작업하는 코드가 포함됩니다.
   4. service: 프로젝트 내에서 사용되는 서비스 라이브러리들이 존재합니다. 주로 비즈니스 로직이 포함됩니다.
3. internal: 내부에서만 쓰이는 라이브러리들이 존재합니다. 외부에 공개되지 않는 코드들이 포함됩니다. 개인적으로는 이 폴더 사용을 최소한으로 하길 원합니다.
4. docs: 문서화에 필요한 자료들이 들어갑니다. UML이나 D2같은 문서들이나, 코드 내에 주석으로 담지 못 하는 정보들이 포함됩니다.
5. template: 제티가 프로젝트를 생성할 때 사용하는 파일들이 들어갑니다.
   1. proto: 프로토버퍼 파일들이 들어갑니다.
   2. model: json, yaml, xml 등의 샘플 파일이 포함됩니다.
   3. config: 설정 파일을 생성하기 위한 템플릿 파일들이 들어갑니다.
6. gen: 제티가 생성한 파일들이 들어갑니다. 보통은 template 폴더 내의 파일들을 기반으로 생성됩니다.
   1. proto: 프로토버퍼 파일들이 컴파일되어 생성됩니다.
   2. model: json, yaml, xml 등의 샘플 파일들이 고랭 구조체로 변환되어 생성됩니다.
   3. bean: 프로젝트 내에서 사용하는 빈(bean) 컨테이너를 생성합니다.

### 주의점

1. 단일 `go.mod` 파일을 가지지만, 여러 실행 파일이 생성되기를 원합니다. `cmd` 폴더 내에 여러 엔트리 포인트를 생성합니다.
2. 전역 변수 사용을 지양합니다. 전역 변수의 사용을 대체하기 위해, `bean`을 생성하여 의존성 주입을 합니다.
3. 프로토버퍼를 VO와 DTO로 사용할 것을 권장합니다. VO는 내부에서만 사용하고, DTO는 외부 통신에 사용합니다.
4. 프론트에 따라 json을 사용할 수도 있습니다.
5. 서버 간 통신에는 gRPC를 권장하고, 웹 클라이언트와의 통신에는 json을 권장합니다.
6. 데이터의 흐름은 서버 -> 서비스 -> 워커 or 클라이언트 -> 서비스 -> 서버를 권장합니다.
7. 문서는 대부분 주석으로 처리합니다. `docs`는 더욱 상세한 부분이나 히스토리를 포함합니다.
8. `internal` 폴더는 최대한 사용하지 않습니다. `lib` 폴더를 사용합니다.

## 기능

### 프로젝트 관련

#### 프로젝트 생성

```bash
jetti new --init <project-name>
```

프로젝트를 생성합니다. `--init` 옵션의 값으로 프로젝트 이름을 입력합니다.  
그러면 자동으로 필요한 폴더들을 만들고, `go mod init <project-name>`을 실행합니다.

#### 엔트리 포인트 생성

```bash
jetti cmd --new <cmd-name>
```

`cmd` 폴더 내에 새로운 엔트리 포인트를 생성합니다. `--new` 옵션의 값으로 엔트리 포인트 이름을 입력합니다.

#### 엔트리 포인트 실행

```bash
jetti cmd --run=<cmd-name>,<args1>,<args2>,...
```

`cmd` 폴더 내에 해당하는 이름의 엔트리 포인트를 실행합니다. `--run` 옵션의 값으로 엔트리 포인트 이름과 인자들을 입력합니다.

#### 엔트리 포인트 빌드

```bash
jetti cmd --build=<cmd-name>,<option1>,<option2>,...
```

`cmd` 폴더 내에 해당하는 이름의 엔트리 포인트를 빌드합니다. `--build` 옵션의 값으로 엔트리 포인트 이름과 빌드 옵션들을 입력합니다.  
빌드된 결과물은 프로젝트 루트 디렉토리에 생성됩니다.

### 빈 컨테이너

#### 빈 구조체 표시

```go
package person

//go:bean
type Person struct {
	name string
	age  int
}
```

`go:bean` 주석을 추가하여 해당 구조체가 빈 컨테이너에 등록될 수 있도록 합니다.

#### 빈 컨테이너 생성

```bash
jetti bean --generate
```

프로젝트의 고 언어 파일을 모두 읽어서 `go:bean` 주석이 달린 구조체들을 찾아서 빈 컨테이너를 생성합니다.

#### 빈 컨테이너 사용

제티를 통해 빈 컨테이너를 생성하면 `gen/bean` 폴더 아래에 다음과 같은 파일이 생성됩니다.

```go
package bean

import (
	".../person"
)

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

func (b *Builder) AddPerson(person *person.Person) *Builder {
	b.bean.person = person
	return b
}

func (b *Bean) Person() *person.Person {
	return b.person
}
```

간단한 빌더 패턴으로 빈 컨테이너를 생성할 수 있습니다.

```go
package main

import (
    ".../bean"
    ".../person"
)

func main() {
    person := &person.Person{
        name: "jetti",
        age:  20,
    }
    bean := bean.New().AddPerson(person).Build()
    println(bean.Person().name)
}
```

`person` 구조체가 싱글톤이길 원하고, 프로젝트 전반에 써야한다면, 앞으로 생성될 하위 로직들에 `bean`을 주입해주면 됩니다.

### 프로토버퍼

프로토버퍼 파일을 컴파일 하기 위해선 `protoc`와 `protoc-gen-go`가 필요합니다.

프로토버퍼 컴파일러는 [이곳](https://github.com/protocolbuffers/protobuf#protocol-compiler-installation)을 참고해서 설치해주세요.
> 맥과 홈브루를 사용하신다면, `brew install protobuf`로 설치하실 수 있습니다.

고 언어 용 프로토버퍼 플러그인은 `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`을 실행하여 설치하실 수 있습니다.

#### 새 프로토버퍼 파일 생성

```bash
jetti proto --new <path/proto-name>
```

`proto` 폴더 내에 새로운 프로토버퍼 파일을 생성합니다. `--new` 옵션의 값으로 프로토버퍼 파일의 경로와 이름을 입력합니다.

예를 들어, `jetti proto -new person/person.proto`를 입력하면 `proto/person/person.proto` 파일이 생성됩니다.

#### 프로토버퍼 파일 생성

```bash
jetti proto --generate
```

`proto` 폴더 내에 있는 모든 프로토버퍼 파일을 컴파일하여 `gen/proto` 폴더에 생성합니다.

#### 프로토버퍼 파일 사용

생성된 프로토버퍼 메시지나 gRPC 서비스를 바로 사용하실 수 있습니다.

프로토버퍼 파일을 DTO로 적극 활용하기를 권장합니다.

### 플랫버퍼

플랫버퍼를 사용하기 위해선 `flatc`를 설치해야합니다.

[이곳](https://github.com/google/flatbuffers/releases)에서 최신 버전을 받은 후, 환경 변수에 추가해주세요.

#### 새 플랫버퍼 파일 생성

```bash
jetti fbs --new <path/fbs-name>
```

`fbs` 폴더 내에 새로운 플랫버퍼 파일을 생성합니다. `--new` 옵션의 값으로 플랫버퍼 파일의 경로와 이름을 입력합니다.

예를 들어, `jetti fbs -new person/person.fbs`를 입력하면 `fbs/person/person.fbs` 파일이 생성됩니다.

#### 플랫버퍼 파일 생성

```bash
jetti fbs --generate
```

`fbs` 폴더 내에 있는 모든 플랫버퍼 파일을 컴파일하여 `gen/fbs` 폴더에 생성합니다.

### 프로파일링 서버 생성

고 언어에서 지원하는 `net/pprof` 서버에 대해 `gen/proto` 폴더 밑에 서버 코드를 생성하는 기능을 제공합니다.

#### http1

```bash
jetti pprof --http-1 <addr>
```

`--http-1` 옵션의 값으로 http1 서버의 주소를 입력합니다.

그러면 `gen/pprof/http1` 폴더에 http1 서버 코드가 생성됩니다.

#### http2

```bash
jetti pprof --http-2 <addr>
```

`--http-2` 옵션의 값으로 http2 서버의 주소를 입력합니다.

그러면 `gen/pprof/http2` 폴더에 http2 서버 코드가 생성됩니다.

### http3

```bash
jetti pprof --http-3 <addr>
```

`--http-3` 옵션의 값으로 http3 서버의 주소를 입력합니다.

그러면 `gen/pprof/http3` 폴더에 http3 서버 코드가 생성됩니다.

### json, yaml 모델 생성

제티에는 `github.com/twpayne/go-jsonstruct/v2` 라이브러리의 도움으로 json, yaml 파일을 바로 고 언어 구조체로 변환하는 기능을 제공합니다.

#### 모델 파일 생성

```bash
jetti model --new <path/name>.(json|yaml)
```

`model` 폴더 내에 새로운 모델 파일을 생성합니다. `--new` 옵션의 값으로 모델 파일의 경로와 이름을 입력합니다.

- `jetti model -new person/person.json`를 입력하면 `template/model/person/person.json` 파일이 생성됩니다.  
- `jetti model -new person/person.yaml`를 입력하면 `template/model/person/person.yaml` 파일이 생성됩니다.

#### 모델 파일 변환

```bash
# json
jetti model --json <path/name>.json

# yaml
jetti model --yaml <path/name>.yaml
```

- `--json` 옵션의 값으로 json 파일의 경로와 이름을 입력하면 `gen/model` 폴더에 json 파일을 바로 고 언어 구조체로 변환하여 생성합니다.
- `--yaml` 옵션의 값으로 yaml 파일의 경로와 이름을 입력하면 `gen/model` 폴더에 yaml 파일을 바로 고 언어 구조체로 변환하여 생성합니다.

### 설정 파일 생성

제티는 `github.com/google/go-jsonnet` 라이브러리의 도움으로 `jsonnet` 파일을 `json`으로 컴파일 할 수 있습니다.

#### 설정 파일 생성

```bash
jetti config --new <path/name>.jsonnet
```

`config` 폴더 내에 새로운 설정 파일을 생성합니다. `--new` 옵션의 값으로 설정 파일의 경로와 이름을 입력합니다.

예를 들어, `jetti config -new person/person.jsonnet`를 입력하면 `template/config/person/person.jsonnet` 파일이 생성됩니다.

#### 설정 파일 컴파일

```bash
jetti config --jsonnet <path/name>.jsonnet
```

`--jsonnet` 옵션의 값으로 설정 파일의 경로와 이름을 입력하면 `gen/config` 폴더에 json 파일로 컴파일하여 생성합니다.

예를 들어, `jetti config --jsonnet person/person.jsonnet`를 입력하면 `gen/config/person/person.json` 파일이 생성됩니다.
