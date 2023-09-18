package commuser

import "microsvc/proto/model/user"

func UserListToMap(list []*user.User) (umap map[int64]*user.User) {
	if len(list) > 0 {
		umap = make(map[int64]*user.User)
		for _, i := range list {
			umap[i.Uid] = i
		}
	}
	return
}

func MergeUserListToMap(umap map[int64]*user.User, list []*user.User) {
	for _, i := range list {
		umap[i.Uid] = i
	}
}
