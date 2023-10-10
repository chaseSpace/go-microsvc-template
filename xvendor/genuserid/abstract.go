package genuserid

import "context"

// UidGeneratorApi
/*
支持特性：
- 指定起始id
- 跳过指定uid（自定义逻辑，可使用正则等方式）
*/
type UidGeneratorApi interface {
	UpdateStartUid(uint64)
	GenUid(ctx context.Context) (uint64, error) // 通过ctx设置timeout
}
