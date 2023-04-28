package executor

import (
	"bytes"
	"fmt"
	"github.com/snowmerak/jetti/internal/finder"
	"io/fs"
	"log"
	"os"
	"os/exec"
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
	moduleName, err := finder.FindModuleName()
	if err != nil {
		panic(err)
	}

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

		generateFile := generated + "/redis/" + strings.TrimPrefix(path, "template/redis/")
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

		dependencies := getDependencies(dataList)

		buffer := bytes.NewBuffer(nil)
		buffer.WriteString(fmt.Sprintf("package %s\n\n", packageName))
		buffer.WriteString("import (\n")
		buffer.WriteString("\t\"context\"\n")
		buffer.WriteString("\t\"github.com/rueian/rueidis\"\n")
		if len(dependencies) > 0 {
			buffer.WriteString("\t\"google.golang.org/protobuf/proto\"\n")
			buffer.WriteString("\t\"encoding/base64\"\n")
		}
		for _, dep := range dependencies {
			buffer.WriteString(fmt.Sprintf("\t\"%s/%s\"\n", moduleName, dep.Import))
		}
		buffer.WriteString(")\n\n")

		for _, data := range dataList {
			split := strings.Split(data, " ")
			if len(split) != 3 {
				continue
			}

			switch split[0] {
			case "string":
				if err := generateRedisString(buffer, split[1], split[2]); err != nil {
					return err
				}
			case "list":
				if err := generateRedisList(buffer, split[1], split[2]); err != nil {
					return err
				}
			case "set":
				if err := generateRedisSet(buffer, split[1], split[2]); err != nil {
					return err
				}
			case "bitmap":
				if err := generateRedisBitMap(buffer, split[1], split[2]); err != nil {
					return err
				}
			default:
				switch {
				case strings.HasPrefix(split[0], "string"):
					dep := dependencies[data]
					if err := generateRedisStringProto(buffer, split[1], split[2], dep.Type); err != nil {
						return err
					}
				case strings.HasPrefix(split[0], "list"):
					dep := dependencies[data]
					if err := generateRedisListProto(buffer, split[1], split[2], dep.Type); err != nil {
						return err
					}
				case strings.HasPrefix(split[0], "set"):
					dep := dependencies[data]
					if err := generateRedisSetProto(buffer, split[1], split[2], dep.Type); err != nil {
						return err
					}
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

	switch output, err := exec.Command("go", "get", "-u", "github.com/rueian/rueidis").Output(); err.(type) {
	case nil:
		log.Println(string(output))
	default:
		panic(err)
	}
}

func generateRedisString(w *bytes.Buffer, name, key string) error {
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

	w.WriteString(fmt.Sprintf("func (t *%s) ExpireAt(ctx context.Context, timestamp int64) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Expireat().Key(t.key).Timestamp(timestamp).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Persist(ctx context.Context) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Persist().Key(t.key).Build()).Error()\n")
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

func generateRedisStringProto(w *bytes.Buffer, name string, key string, typ string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Get(ctx context.Context) (*%s, error) {\n", name, typ))
	w.WriteString(fmt.Sprintf("\tvalue, err := t.client.Do(ctx, t.client.B().Get().Key(t.key).Build()).ToString()\n"))
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tbuf, err := base64.StdEncoding.DecodeString(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := new(%s)\n", typ))
	w.WriteString("\tif err := proto.Unmarshal(buf, result); err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Set(ctx context.Context, value *%s) error {\n", name, typ))
	w.WriteString("\tbuf, err := proto.Marshal(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tencoded := base64.StdEncoding.EncodeToString(buf)\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Set().Key(t.key).Value(encoded).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Del(ctx context.Context) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Del().Key(t.key).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Exists(ctx context.Context) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Exists().Key(t.key).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Expire(ctx context.Context, seconds int64) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Expire().Key(t.key).Seconds(seconds).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) ExpireAt(ctx context.Context, timestamp int64) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Expireat().Key(t.key).Timestamp(timestamp).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Persist(ctx context.Context) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Persist().Key(t.key).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) TTL(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Ttl().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	return nil
}

func generateRedisList(w *bytes.Buffer, name, key string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPush(ctx context.Context, values ...string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lpush().Key(t.key).Element(values...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) RPush(ctx context.Context, values ...string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Rpush().Key(t.key).Element(values...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPop(ctx context.Context) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lpop().Key(t.key).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) RPop(ctx context.Context) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Rpop().Key(t.key).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LRange(ctx context.Context, start, end int64) ([]string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lrange().Key(t.key).Start(start).Stop(end).Build()).AsStrSlice()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LTrim(ctx context.Context, start, end int64) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Ltrim().Key(t.key).Start(start).Stop(end).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LLen(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Llen().Key(t.key).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LIndex(ctx context.Context, index int64) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lindex().Key(t.key).Index(index).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LSet(ctx context.Context, index int64, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lset().Key(t.key).Index(index).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LRem(ctx context.Context, count int64, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lrem().Key(t.key).Count(count).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LInsertBefore(ctx context.Context, pivot, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Linsert().Key(t.key).Before().Pivot(pivot).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LInsertAfter(ctx context.Context, pivot, value string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Linsert().Key(t.key).After().Pivot(pivot).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPos(ctx context.Context, value string, rank int64) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lpos().Key(t.key).Element(value).Rank(rank).Build()).ToInt64()\n")
	w.WriteString("}\n\n")

	return nil
}

func generateRedisListProto(w *bytes.Buffer, name, key, typ string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString(fmt.Sprintf("\tkey string\n"))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPush(ctx context.Context, values ...*%s) error {\n", name, typ))
	w.WriteString("\telements := make([]string, 0, len(values))\n")
	w.WriteString("\tfor _, v := range values {\n")
	w.WriteString("\t\telement, err := proto.Marshal(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\telements = append(elements, base64.StdEncoding.EncodeToString(element))\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lpush().Key(t.key).Element(elements...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) RPush(ctx context.Context, values ...*%s) error {\n", name, typ))
	w.WriteString("\telements := make([]string, 0, len(values))\n")
	w.WriteString("\tfor _, v := range values {\n")
	w.WriteString("\t\telement, err := proto.Marshal(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\telements = append(elements, base64.StdEncoding.EncodeToString(element))\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Rpush().Key(t.key).Element(elements...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPop(ctx context.Context) (*%s, error) {\n", name, typ))
	w.WriteString("\tvalue, err := t.client.Do(ctx, t.client.B().Lpop().Key(t.key).Build()).ToString()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tdata, err := base64.StdEncoding.DecodeString(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := new(%s)\n", typ))
	w.WriteString("\tif err := proto.Unmarshal(data, result); err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) RPop(ctx context.Context) (*%s, error) {\n", name, typ))
	w.WriteString("\tvalue, err := t.client.Do(ctx, t.client.B().Rpop().Key(t.key).Build()).ToString()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tdata, err := base64.StdEncoding.DecodeString(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := new(%s)\n", typ))
	w.WriteString("\tif err := proto.Unmarshal(data, result); err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LRange(ctx context.Context, start, stop int64) ([]*%s, error) {\n", name, typ))
	w.WriteString("\tvalues, err := t.client.Do(ctx, t.client.B().Lrange().Key(t.key).Start(start).Stop(stop).Build()).AsStrSlice()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := make([]*%s, 0, len(values))\n", typ))
	w.WriteString("\tfor _, v := range values {\n")
	w.WriteString("\t\tdata, err := base64.StdEncoding.DecodeString(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString(fmt.Sprintf("\t\titem := new(%s)\n", typ))
	w.WriteString("\t\tif err := proto.Unmarshal(data, item); err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tresult = append(result, item)\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LLen(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Llen().Key(t.key).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LTrim(ctx context.Context, start, stop int64) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Ltrim().Key(t.key).Start(start).Stop(stop).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LIndex(ctx context.Context, index int64) (*%s, error) {\n", name, typ))
	w.WriteString("\tvalue, err := t.client.Do(ctx, t.client.B().Lindex().Key(t.key).Index(index).Build()).ToString()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tdata, err := base64.StdEncoding.DecodeString(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := new(%s)\n", typ))
	w.WriteString("\tif err := proto.Unmarshal(data, result); err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LSet(ctx context.Context, index int64, element *%s) error {\n", name, typ))
	w.WriteString("\tdata, err := proto.Marshal(element)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tvalue := base64.StdEncoding.EncodeToString(data)\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lset().Key(t.key).Index(index).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LRem(ctx context.Context, count int64, element *%s) error {\n", name, typ))
	w.WriteString("\tdata, err := proto.Marshal(element)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tvalue := base64.StdEncoding.EncodeToString(data)\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lrem().Key(t.key).Count(count).Element(value).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) LPos(ctx context.Context, element *%s, rank int64) (int64, error) {\n", name, typ))
	w.WriteString("\tdata, err := proto.Marshal(element)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn 0, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tvalue := base64.StdEncoding.EncodeToString(data)\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Lpos().Key(t.key).Element(value).Rank(rank).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	return nil
}

func generateRedisSet(w *bytes.Buffer, name, key string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey    string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SAdd(ctx context.Context, members ...string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Sadd().Key(t.key).Member(members...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SCard(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Scard().Key(t.key).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SIsMember(ctx context.Context, member string) (bool, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Sismember().Key(t.key).Member(member).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SMembers(ctx context.Context) ([]string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Smembers().Key(t.key).Build()).AsStrSlice()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SPop(ctx context.Context) (string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Spop().Key(t.key).Build()).ToString()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SRandMember(ctx context.Context, count int64) ([]string, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Srandmember().Key(t.key).Count(count).Build()).AsStrSlice()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SRem(ctx context.Context, members ...string) error {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Srem().Key(t.key).Member(members...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SScan(ctx context.Context, cursor uint64, match string, count int64) (uint64, []string, error) {\n", name))
	w.WriteString("\tse, err := t.client.Do(ctx, t.client.B().Sscan().Key(t.key).Cursor(cursor).Match(match).Count(count).Build()).AsScanEntry()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn 0, nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn se.Cursor, se.Elements, nil\n")
	w.WriteString("}\n\n")

	return nil
}

func generateRedisSetProto(w *bytes.Buffer, name, key, typ string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey    string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SAdd(ctx context.Context, members ...*%s) error {\n", name, typ))
	w.WriteString("\tvalues := make([]string, 0, len(members))\n")
	w.WriteString("\tfor _, m := range members {\n")
	w.WriteString("\t\tdata, err := proto.Marshal(m)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tvalues = append(values, base64.StdEncoding.EncodeToString(data))\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Sadd().Key(t.key).Member(values...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SCard(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Scard().Key(t.key).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SIsMember(ctx context.Context, member *%s) (bool, error) {\n", name, typ))
	w.WriteString("\tdata, err := proto.Marshal(member)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn false, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Sismember().Key(t.key).Member(base64.StdEncoding.EncodeToString(data)).Build()).ToBool()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SMembers(ctx context.Context) ([]*%s, error) {\n", name, typ))
	w.WriteString("\tvalues, err := t.client.Do(ctx, t.client.B().Smembers().Key(t.key).Build()).AsStrSlice()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := make([]*%s, 0, len(values))\n", typ))
	w.WriteString("\tfor _, v := range values {\n")
	w.WriteString("\t\tdata, err := base64.StdEncoding.DecodeString(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString(fmt.Sprintf("\t\tm := new(%s)\n", typ))
	w.WriteString("\t\tif err := proto.Unmarshal(data, m); err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tresult = append(result, m)\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SPop(ctx context.Context) (*%s, error) {\n", name, typ))
	w.WriteString("\tvalue, err := t.client.Do(ctx, t.client.B().Spop().Key(t.key).Build()).ToString()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\tdata, err := base64.StdEncoding.DecodeString(value)\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tm := new(%s)\n", typ))
	w.WriteString("\tif err := proto.Unmarshal(data, m); err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn m, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SRandMember(ctx context.Context, count int64) ([]*%s, error) {\n", name, typ))
	w.WriteString("\tvalues, err := t.client.Do(ctx, t.client.B().Srandmember().Key(t.key).Count(count).Build()).AsStrSlice()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := make([]*%s, 0, len(values))\n", typ))
	w.WriteString("\tfor _, v := range values {\n")
	w.WriteString("\t\tdata, err := base64.StdEncoding.DecodeString(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString(fmt.Sprintf("\t\tm := new(%s)\n", typ))
	w.WriteString("\t\tif err := proto.Unmarshal(data, m); err != nil {\n")
	w.WriteString("\t\t\treturn nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tresult = append(result, m)\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn result, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SRem(ctx context.Context, members ...*%s) error {\n", name, typ))
	w.WriteString("\tvalues := make([]string, 0, len(members))\n")
	w.WriteString("\tfor _, member := range members {\n")
	w.WriteString("\t\tdata, err := proto.Marshal(member)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tvalues = append(values, base64.StdEncoding.EncodeToString(data))\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Srem().Key(t.key).Member(values...).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) SScan(ctx context.Context, cursor uint64, match string, count int64) (uint64, []*%s, error) {\n", name, typ))
	w.WriteString("\tse, err := t.client.Do(ctx, t.client.B().Sscan().Key(t.key).Cursor(cursor).Match(match).Count(count).Build()).AsScanEntry()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn 0, nil, err\n")
	w.WriteString("\t}\n")
	w.WriteString(fmt.Sprintf("\tresult := make([]*%s, 0, len(se.Elements))\n", typ))
	w.WriteString("\tfor _, v := range se.Elements {\n")
	w.WriteString("\t\tdata, err := base64.StdEncoding.DecodeString(v)\n")
	w.WriteString("\t\tif err != nil {\n")
	w.WriteString("\t\t\treturn 0, nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString(fmt.Sprintf("\t\tm := new(%s)\n", typ))
	w.WriteString("\t\tif err := proto.Unmarshal(data, m); err != nil {\n")
	w.WriteString("\t\t\treturn 0, nil, err\n")
	w.WriteString("\t\t}\n")
	w.WriteString("\t\tresult = append(result, m)\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn se.Cursor, result, nil\n")
	w.WriteString("}\n\n")

	return nil
}

func generateRedisBitMap(w *bytes.Buffer, name, key string) error {
	w.WriteString(fmt.Sprintf("type %s struct {\n", name))
	w.WriteString("\tclient rueidis.Client\n")
	w.WriteString("\tkey    string\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func New%s(client rueidis.Client) *%s {\n", name, name))
	w.WriteString(fmt.Sprintf("\treturn &%s{client: client, key: \"%s\"}\n", name, key))
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) BitCount(ctx context.Context, start, end int64) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Bitcount().Key(t.key).Start(start).End(end).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Set(ctx context.Context, offset int64, value bool) error {\n", name))
	w.WriteString("v := int64(0)")
	w.WriteString("\tif value {\n")
	w.WriteString("\t\tv = 1\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Setbit().Key(t.key).Offset(offset).Value(v).Build()).Error()\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Get(ctx context.Context, offset int64) (bool, error) {\n", name))
	w.WriteString("\tv, err := t.client.Do(ctx, t.client.B().Getbit().Key(t.key).Offset(offset).Build()).AsInt64()\n")
	w.WriteString("\tif err != nil {\n")
	w.WriteString("\t\treturn false, err\n")
	w.WriteString("\t}\n")
	w.WriteString("\treturn v == 1, nil\n")
	w.WriteString("}\n\n")

	w.WriteString(fmt.Sprintf("func (t *%s) Count(ctx context.Context) (int64, error) {\n", name))
	w.WriteString("\treturn t.client.Do(ctx, t.client.B().Bitcount().Key(t.key).Build()).AsInt64()\n")
	w.WriteString("}\n\n")

	return nil
}
