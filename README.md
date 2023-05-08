# jetti

제티(jetti)는 고 언어 프로젝트에 어느정도 정규화된 프로젝트 구성을 적용하기 위한 도구입니다.

## 지향점

[WIP]

## 기능

### bean

`bean`은 `context.Context`의 `WithValue`를 편리하게 이용할 수 있게 해주는 기능입니다.

빈의 핵심은 동등, 혹은 하위 문맥에서 동일한 객체 인스턴스를 공유하는 것입니다.

#### 예시

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

### json/yaml to go

`go-jsonstruct` 라이브러리를 이용해서 json/yaml 파일을 go 구조체로 변환할 수 있습니다.

파싱에는 각각 `goccy/go-json`과 `goccy/go-yaml` 라이브러리를 사용합니다.

#### 예시

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