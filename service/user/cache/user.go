package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"microsvc/bizcomm"
	"microsvc/proto/model/user"
	"microsvc/service/user/dao"
)

func GetUser(uid ...int64) (umap map[int64]*user.User, err error) {
	var keys []string
	for _, u := range uid {
		keys = append(keys, fmt.Sprintf(UserCacheKey, u))
	}
	reply := user.R.MGet(context.Background(), keys...)
	if reply.Err() != nil {
		return nil, err
	}
	umap = make(map[int64]*user.User, len(uid))
	var queryDBUids []int64
	for i, v := range reply.Val() {
		if v == nil {
			queryDBUids = append(queryDBUids, uid[i])
			continue
		}
		u := new(user.User)
		_ = json.Unmarshal([]byte(v.(string)), u)
		umap[u.Uid] = u
	}

	if len(queryDBUids) > 0 {
		list, _, err := dao.GetUser(queryDBUids...)
		if err != nil {
			return nil, err
		}
		bizcomm.MergeUserListToMap(umap, list)
	}
	return
}
