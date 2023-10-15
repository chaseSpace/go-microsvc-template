package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"microsvc/proto/model/user"
	"microsvc/service/user/dao"
)

type UserMap map[int64]*user.User

func GetUser(uid ...int64) (umap UserMap, err error) {
	var keys []string
	for _, u := range uid {
		keys = append(keys, fmt.Sprintf(UserInfoCacheKey, u))
	}
	reply := user.R.MGet(context.Background(), keys...)
	if reply.Err() != nil {
		return nil, reply.Err()
	}
	umap = make(UserMap, len(uid))

	var cacheMissUids []int64
	for i, v := range reply.Val() {
		if v == nil {
			cacheMissUids = append(cacheMissUids, uid[i])
			continue
		}
		u := new(user.User)
		_ = json.Unmarshal([]byte(v.(string)), u)
		umap[u.Uid] = u
	}

	if len(cacheMissUids) > 0 {
		list, _, err := dao.GetUser(cacheMissUids...)
		if err != nil {
			return nil, err
		}
		lo.ForEach(list, func(item *user.User, index int) {
			umap[item.Uid] = item
		})
	}
	return
}
