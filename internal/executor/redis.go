package executor

import (
	"bytes"
	"fmt"
	"github.com/snowmerak/jetti/internal/finder"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func RedisNew(path string) {
	name := filepath.Base(path)
	dir := filepath.Join("template", "redis", path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	if !strings.HasSuffix(path, ".go") {
		path += ".go"
	}
	f, err := os.Create(filepath.Join(dir, name+".go"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()

	if _, err := f.Write([]byte(fmt.Sprintf("package %s\n", name))); err != nil {
		panic(err)
	}
}

func RedisGenerate() {
	if err := filepath.Walk("./template/redis", func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}(f)

		dataList := finder.FindDirections(f, "redis")

		if len(dataList) == 0 {
			return nil
		}

		generateFile := "generated/redis/" + strings.TrimPrefix(path, "template/redis/")
		if err := os.MkdirAll(filepath.Dir(generateFile), os.ModePerm); err != nil {
			return err
		}

		packageName := filepath.Base(filepath.Dir(generateFile))

		f, err = os.Create(generateFile)
		if err != nil {
			return err
		}
		defer func(f *os.File) {
			if err := f.Close(); err != nil {
				panic(err)
			}
		}(f)

		buffer := bytes.NewBuffer(nil)
		buffer.WriteString(fmt.Sprintf("package %s\n\n", packageName))
		buffer.WriteString("import (\n")
		buffer.WriteString("\t\"context\"\n")
		buffer.WriteString("\t\"github.com/rueian/rueidis\"\n")
		buffer.WriteString(")\n\n")

		for _, data := range dataList {
			split := strings.Split(data, " ")
			if len(split) != 3 {
				continue
			}

			switch split[0] {
			case "string":
				if err := writeRedisStringType(buffer, split[1], split[2]); err != nil {
					return err
				}
			}
		}

		if _, err := f.Write(buffer.Bytes()); err != nil {
			return err
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func writeRedisStringType(w *bytes.Buffer, name, key string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Get(ctx context.Context) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Get().Key(t.key).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) GetDel(ctx context.Context) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Getdel().Key(t.key).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Set(ctx context.Context, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Set().Key(t.key).Value(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SetIfNotExist(ctx context.Context, value string) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Setnx().Key(t.key).Value(value).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SetWithTTL(ctx context.Context, value string, ttlSeconds int64) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Set().Key(t.key).Value(value).ExSeconds(ttlSeconds).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Del(ctx context.Context) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Del().Key(t.key).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Expire(ctx context.Context, ttlSeconds int64) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Expire().Key(t.key).Seconds(ttlSeconds).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) TTL(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Ttl().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Append(ctx context.Context, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Append().Key(t.key).Value(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Incr(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Incr().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) IncrBy(ctx context.Context, value int64) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Incrby().Key(t.key).Increment(value).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) IncrByFloat(ctx context.Context, value float64) (float64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Incrbyfloat().Key(t.key).Increment(value).Build()).ToFloat64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Decr(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Decr().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) DecrBy(ctx context.Context, value int64) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Decrby().Key(t.key).Decrement(value).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) GetRange(ctx context.Context, start, end int64) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Getrange().Key(t.key).Start(start).End(end).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SetRange(ctx context.Context, offset int64, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Setrange().Key(t.key).Offset(offset).Value(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Len(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Strlen().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	return nil
}
