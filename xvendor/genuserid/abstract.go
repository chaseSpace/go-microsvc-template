package genuserid

import "context"

// UIDGeneratorApi
/*
支持特性：
- 指定起始id
- 跳过指定uid（自定义逻辑，可使用正则等方式）
- 递增
- 使用号池模式，支持高并发(见 TestConcurrencyGenUID )
*/
type UIDGeneratorApi interface {
	GenUid(ctx context.Context) (uint64, error) // 通过ctx设置timeout
}
