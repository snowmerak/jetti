---
title: 'Bean'
date: 2023-05-16T21:12:23+09:00
---

`bean` direction is an instruction for `jetti` to create a bean container based on the structure or interface, type alias.

`bean` is a core of dependency injection.

## Usage

```go
package config

// jetti:bean redis postgres
type Config struct {
}
```

Above code mark `Config` as a bean container named `redis`, `postgres`.

## Generated code

Run `jetti generate` command, `jetti` generates the following code.

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

## Example

We can use the generated code as follows.

```go
package main

import (
    "context"
    "fmt"

    "../config"
)

func main() {
	ctx := context.Background()
    ctx = config.PushPostgres(ctx, &config.Config{})
    ctx = config.PushRedis(ctx, &config.Config{})

    postgres, exists := config.GetPostgres(ctx)
	if !exists {
		panic("postgres not exists")
	}
    redis, exists := config.GetRedis(ctx)
	if !exists {
		panic("redis not exists")
	}

    fmt.Println(postgres, redis)
}
```
