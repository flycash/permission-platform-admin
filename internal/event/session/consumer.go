package session

import (
	"encoding/json"
	"fmt"
	"gitee.com/flycash/permission-platform-admin/internal/domain"
	"gitee.com/flycash/permission-platform-admin/internal/pkg/mqx"
	"time"

	"github.com/ecodeclub/ekit/slice"

	"github.com/gotomicro/ego/core/elog"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

const (
	number36       = 36
	defaultTimeout = 5 * time.Second
	permissionName = "permission"
)

type Consumer struct {
	client   redis.Cmdable
	consumer mqx.Consumer
	logger   *elog.Component
}

func NewConsumer(client redis.Cmdable, consumer mqx.Consumer) *Consumer {
	return &Consumer{
		client:   client,
		consumer: consumer,
		logger:   elog.DefaultLogger,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
			c.Consume(ctx)
			cancel()
		}
	}()
}

func (c *Consumer) Consume(ctx context.Context) {
	msg, err := c.consumer.ReadMessage(-1)
	if err != nil {
		return
	}
	var evt UserPermissionEvent
	err = json.Unmarshal(msg.Value, &evt)
	if err != nil {
		c.logger.Error("解析消息失败",
			elog.FieldErr(err),
			elog.Any("msg", msg))
		return
	}
	vals := make([]any, 0, number36)
	sessionIDMap, err := c.getSessionIDMap(ctx, evt)
	if err != nil {
		c.logger.Error("获取sessioniID失败",
			elog.FieldErr(err),
			elog.Any("msg", msg))
		return
	}

	pipeline := c.client.Pipeline()
	for uid := range evt.Permissions {
		sessionID, ok := sessionIDMap[fmt.Sprintf("%d", uid)]
		if !ok {
			continue
		}
		key := c.key(sessionID)
		userPermission := evt.Permissions[uid]
		domainPermissions := slice.Map(userPermission.Permissions, func(_ int, src Permission) domain.UserPermission {
			return domain.UserPermission{
				BizID:  userPermission.BizID,
				UserID: uid,
				Permission: domain.Permission{
					Resource: domain.Resource{
						Type: src.Resource.Type,
						Key:  src.Resource.Key,
					},
					Action: src.Action,
				},
				Effect: domain.Effect(src.Effect),
			}
		})
		permissionByte, err := json.Marshal(domainPermissions)
		if err != nil {
			c.logger.Error("序列化权限消息失败",
				elog.FieldErr(err),
				elog.Int64("uid", uid),
				elog.Any("permissions", domainPermissions))
			return
		}
		vals = append(vals, key, string(permissionByte))
		pipeline.HSet(ctx, c.key(sessionID), key, string(permissionByte))
	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		c.logger.Error("保存到redis失败",
			elog.FieldErr(err),
			elog.Any("msg", msg))
	}
}

func (c *Consumer) getSessionIDMap(ctx context.Context, evt UserPermissionEvent) (map[string]string, error) {
	uids := make([]string, 0, len(evt.Permissions))
	for uid := range evt.Permissions {
		uids = append(uids, fmt.Sprintf("%d", uid))
	}
	sliceRes, err := c.client.MGet(ctx, uids...).Result()
	if err != nil {
		return nil, err
	}
	res := make(map[string]string, len(uids))
	for idx := range uids {
		v, ok := sliceRes[idx].(string)
		if ok {
			res[uids[idx]] = v
		}
	}
	return res, nil
}

func (c *Consumer) key(sessionID string) string {
	return fmt.Sprintf("session:%s", sessionID)
}
