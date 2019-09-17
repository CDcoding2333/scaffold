package nosql

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

var ErrNil = errors.New("nil return")
var ErrWrongType = errors.New("wrong type")
var ErrWrongArgsNum = errors.New("args num error")

type Reply interface {
	Int(interface{}, error) (int, error)
	Int64(interface{}, error) (int64, error)
	Uint64(interface{}, error) (uint64, error)
	Float64(interface{}, error) (float64, error)
	Bool(interface{}, error) (bool, error)
	Bytes(interface{}, error) ([]byte, error)
	String(interface{}, error) (string, error)
	Strings(interface{}, error) ([]string, error)
	Values(interface{}, error) ([]interface{}, error)
	Ints(interface{}, error) ([]int, error)
}

type kvStore interface {
	Get(string) (interface{}, error)
	Set(string, interface{}) error
	SetEx(string, interface{}, int64) error
	SetNx(string, interface{}) (int64, error)
	// nil表示成功，ErrNil表示数据库内已经存在这个key，其他表示数据库发生错误
	SetNxEx(string, interface{}, int64) error
	Del(...string) (int64, error)
	Incr(string) (int64, error)
	IncrBy(string, int64) (int64, error)
	Expire(string, int64) (int64, error)
	ExpireAt(string, int64) (int64, error)
	TTl(string) (int64, error)
}

type KVStore interface {
	Reply
	kvStore
}

type hashStore interface {
	HGet(string, string) (interface{}, error)
	HSet(string, string, interface{}) error
	HMGet(string, ...string) (interface{}, error)
	HMSet(string, ...interface{}) error
	HExpire(string, int64) error
	HGetAll(string) (map[string]interface{}, error)
	HIncrBy(string, string, int64) (int64, error)
	HDel(string, ...string) (int64, error)
	HClear(string) error
	HLen(string) (int64, error)
}

type HashStore interface {
	Reply
	hashStore
}

type setStore interface {
	SAdd(string, ...interface{}) (int64, error)
	SIsMember(string, interface{}) (bool, error)
	SRem(string, ...interface{}) (int64, error)
	SMembers(string) (interface{}, error)
}

type SetStore interface {
	Reply
	setStore
}

type zSetStore interface {
	ZAdd(string, ...interface{}) (int64, error)
	ZRem(string, ...string) (int64, error)
	ZExpire(string, int64) error
	ZRange(string, int64, int64, bool) (interface{}, error)
	ZRangeByScoreWithScore(string, int64, int64) (map[string]int64, error)
	ZClear(string) error
}

type ZSetStore interface {
	Reply
	zSetStore
}

type listStore interface {
	LRange(string, int64, int64) (interface{}, error)
	LLen(string) (int64, error)
	LPop(string) (interface{}, error)
	RPop(string) (interface{}, error)
	BLPop(string, int) (interface{}, error)
	BRPop(string, int) (interface{}, error)
	LPush(string, ...interface{}) error
	RPush(string, ...interface{}) error
	MGet(...interface{}) ([]interface{}, error)
}

type ListStore interface {
	Reply
	listStore
}

type HashZSetStore interface {
	Reply
	hashStore
	zSetStore
}

type KVHashStore interface {
	Reply
	kvStore
	hashStore
}

type Store interface {
	Reply
	kvStore
	hashStore
	zSetStore
	listStore
	setStore
	SetMaxIdle(int)
	SetMaxActive(int)
}

func Open(connurl string) (Store, error) {
	scheme, host, password, port, db, err := parseConnurl(connurl)
	if err != nil {
		return nil, err
	}
	switch scheme {
	case "redis":
		return NewRedisStore(host, password, port, db)
	default:
		return nil, fmt.Errorf("invalid connection url %s", connurl)

	}

}

func parseConnurl(connurl string) (scheme string, host, password string, port int, db int, err error) {
	url, err := url.Parse(connurl)
	if err != nil {
		return
	}
	scheme = url.Scheme
	if scheme != "ledis" && scheme != "redis" && scheme != "memory" {
		err = fmt.Errorf("invalid connection url %s", connurl)
		return
	}
	if scheme == "memory" {
		return
	}
	parts := strings.SplitN(url.Host, ":", 3)
	if len(parts) != 3 {
		err = fmt.Errorf("invalid connection url %s", connurl)
		return
	}
	host = parts[0]
	password = parts[1]

	if port, err = strconv.Atoi(parts[2]); err != nil {
		return
	}
	path := strings.Trim(url.Path, "/")
	if db, err = strconv.Atoi(path); err != nil {
		return
	}
	return
}
