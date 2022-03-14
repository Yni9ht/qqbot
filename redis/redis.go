package redis

import (
	"context"
	"fmt"
	"github.com/Yni9ht/qqbot/vars"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	// userSignKey 用户签到位图 key：userId、年月
	userSignKey = "user:sign:%s:%s"
	// userSignScoreKey 用户签到积分 key：userId
	userSignScoreKey = "user:sign:score:%s"
)

type redisSvc struct {
	con *redis.Client
}

func (r *redisSvc) Close() {
	_ = r.con.Close()
}

func (r *redisSvc) CheckUserClockIn(ctx context.Context, userID string, checkTime time.Time) (bool, error) {
	yearmonth := checkTime.Format("200601")
	key := fmt.Sprintf(userSignKey, userID, yearmonth)

	day := int64(checkTime.Day())
	value, err := r.con.GetBit(ctx, key, day).Result()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

func (r *redisSvc) CreateUserClockIn(ctx context.Context, userID string, clockInTime time.Time) error {
	yearmonth := clockInTime.Format("200601")
	key := fmt.Sprintf(userSignKey, userID, yearmonth)

	day := int64(clockInTime.Day())
	_, err := r.con.SetBit(ctx, key, day, 1).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisSvc) IncrUserSignScore(ctx context.Context, userID string, signScore int) (int64, error) {
	key := fmt.Sprintf(userSignScoreKey, userID)

	currentScoreTotal, err := r.con.IncrBy(ctx, key, int64(signScore)).Result()
	if err != nil {
		return 0, err
	}
	return currentScoreTotal, nil
}

// CancelUserClockIn 删除用户签到记录
func (r *redisSvc) CancelUserClockIn(ctx context.Context, userID string, clockInTime time.Time) error {
	yearmonth := clockInTime.Format("200601")
	key := fmt.Sprintf(userSignKey, userID, yearmonth)

	day := int64(clockInTime.Day())
	_, err := r.con.SetBit(ctx, key, day, 0).Result()
	if err != nil {
		return err
	}
	return nil
}

// CancelUserSignScore 删除用户签到积分
func (r *redisSvc) CancelUserSignScore(ctx context.Context, userID string, signScore int) error {
	key := fmt.Sprintf(userSignScoreKey, userID)

	_, err := r.con.IncrBy(ctx, key, int64(-signScore)).Result()
	if err != nil {
		return err
	}
	return nil
}

func NewRedisSvc() *redisSvc {
	redisCon := redis.NewClient(&redis.Options{
		Addr:     vars.Config.GetRedisAddr(),
		Password: vars.Config.Redis.Password,
		DB:       vars.Config.Redis.DB,
	})
	return &redisSvc{
		con: redisCon,
	}
}
