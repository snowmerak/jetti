---
title: 'JSON/YAML converter'
date: 2023-05-16T21:12:23+09:00
---

`Jetti` provides a feature to convert JSON to YAML and YAML to JSON.

`Jetti` uses `go-jsonstruct` to convert JSON to YAML and YAML to JSON.

And `Jetti` uses `go-yaml` and `go-json` to convert YAML/JSON to go struct.

## Example

Make `./config/json_prac.json` and `./config/yaml_prac.yaml` files.

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

## Generated

Run `jetti generate` command, then `./config/json_prac.yaml.go` and `./config/yaml_prac.json.go` files will be generated.

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
