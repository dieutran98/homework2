package internal

import (
	"caching/redis"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type userService struct {
	redisSvc redis.Redis
}

type User struct {
	Name     string `json:"name" redis:"name"`
	Age      int    `json:"age" redis:"age"`
	CacheHit bool   `json:"-"`
}

//go:embed *.json
var f embed.FS

var testData []User

func getTestData() ([]User, error) {
	if testData == nil {
		file, err := f.Open("MOCK_DATA.json")
		if err != nil {
			return nil, errors.Wrap(err, "failed to read file")
		}
		defer file.Close()
		if err := json.NewDecoder(file).Decode(&testData); err != nil {
			return nil, errors.Wrap(err, "failed to decode file")
		}
	}
	return testData, nil
}

func NewService() (*userService, error) {
	r, err := redis.GetRedis()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get redis")
	}
	return &userService{
		redisSvc: r,
	}, nil
}

func (s *userService) GetUsers(ctx context.Context, id string) (*User, error) {
	userId, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to ")
	}

	var u User
	// get data from caching
	if err := s.redisSvc.HGetAll(ctx, makeUserKey(userId), &u); err != nil {
		return nil, errors.Wrap(err, "failed to get user from cache")
	}

	// found
	if u != (User{}) {
		u.CacheHit = true
		return &u, nil
	}

	// get dat from db (testdata)
	data, err := getTestData()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user data")
	}
	if len(data)-1 < userId || userId < 0 {
		return nil, errors.New("id not found")
	}

	// save data to cache
	go s.redisSvc.HSet(ctx, makeUserKey(userId), data[userId])

	return &data[userId], nil
}

func makeUserKey(userId int) string {
	return fmt.Sprintf("user:%d", userId)
}
