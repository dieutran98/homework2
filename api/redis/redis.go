package redis

import (
	"caching/util"
	"context"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type Redis interface {
	HSet(ctx context.Context, key string, value ...any) error
	HGetAll(ctx context.Context, key string, val any) error
}

type redisCli struct {
	redis *redis.Client
}

var myRedisCli *redisCli
var defaultExpireTime = 30

func SetUpdateRedis() error {
	if myRedisCli == nil {
		if err := connectRedis(); err != nil {
			return errors.Wrap(err, "failed to connect redis")
		}
	}
	return nil
}

func GetRedis() (Redis, error) {
	if myRedisCli == nil {
		if err := connectRedis(); err != nil {
			return nil, errors.Wrap(err, "failed to connect redis")
		}
	}
	return myRedisCli, nil
}

func connectRedis() error {
	opts, err := buildOptions()
	if err != nil {
		return errors.Wrap(err, "failed to build options")
	}

	rCli := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

	defer cancel()
	if err := rCli.Ping(ctx).Err(); err != nil {
		return errors.Wrap(err, "ping failed")
	}

	myRedisCli = &redisCli{
		redis: rCli,
	}

	return nil
}

func buildOptions() (*redis.Options, error) {
	e, err := util.GetEnv()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get env")
	}

	options := redis.Options{}

	if err := mapstructure.Decode(e.Redis, &options); err != nil {
		return nil, errors.Wrap(err, "failed to decode options")
	}

	return &options, nil
}

func (r *redisCli) HSet(ctx context.Context, key string, value ...any) error {
	if err := r.redis.HSet(ctx, key, value...).Err(); err != nil {
		return errors.Wrapf(err, "failed to get key '%s'", key)
	}
	// set time out
	if err := r.redis.Expire(ctx, key, time.Second*time.Duration(defaultExpireTime)).Err(); err != nil {
		return errors.Wrapf(err, "failed to set expire time for key %s", key)
	}

	return nil
}

func (r *redisCli) HGetAll(ctx context.Context, key string, val any) error {
	if err := checkValidParam(val); err != nil {
		return errors.Wrap(err, "failed check params")
	}

	if err := r.redis.HGetAll(ctx, key).Scan(val); err != nil {
		return errors.Wrap(err, "failed to get value from redis")
	}
	return nil
}

func checkValidParam(params any) error {
	t := reflect.TypeOf(params)
	if t.Kind() != reflect.Ptr {
		return errors.New("params must be a pointer")
	}
	if t.Elem().Kind() != reflect.Struct {
		return errors.New("params must be pointer of struct")
	}
	return nil
}
