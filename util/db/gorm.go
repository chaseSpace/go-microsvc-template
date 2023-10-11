package db

import (
	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func IsMysqlErr(err error) bool {
	return err != nil && err != gorm.ErrRecordNotFound
}
func IsRedisErr(err error) bool {
	return err != nil && err != redis.Nil
}

/* mysql 常见错误码
Access Denied (1045): 这是一个常见的错误，表示连接到 MySQL 服务器时权限被拒绝。错误代码为 1045。你可以在 MySQL 官方文档中查找这个错误的信息，了解如何解决权限问题。

Table doesn't exist (1146): 当你尝试查询或操作一个不存在的表时，会遇到这个错误。错误代码为 1146。通常需要检查表名是否正确或确保表已经创建。

Duplicate entry (1062): 这个错误表示尝试插入重复的唯一键值。错误代码为 1062。你可以检查你的数据，或者使用 INSERT IGNORE 或 INSERT ... ON DUPLICATE KEY UPDATE 来处理这种情况。

Syntax error (1064): 这个错误表示 SQL 语法错误。错误代码为 1064。你需要检查 SQL 查询或语句的语法，确保它是有效的。

Lock wait timeout exceeded (1205): 当某个事务等待获取锁的时间超过设置的超时时间时，会发生这个错误。错误代码为 1205。你可以尝试增加超时时间，优化查询，或查看锁的情况。

Lost connection to MySQL server (2013): 当与 MySQL 服务器的连接丢失时，会发生这个错误。错误代码为 2013。这可能是由于网络问题或服务器崩溃引起的。你可以检查网络连接或服务器状态。

Data too long for column (1406): 当尝试插入的数据超过了列的最大长度时，会发生这个错误。错误代码为 1406。你需要检查数据的长度，并根据需要调整列的长度。
*/

func IsMysqlDuplicateErr(err error) bool {
	if err == nil {
		return false
	}
	if err, ok := err.(*mysql.MySQLError); ok && err.Number == 1062 {
		return true
	}
	return false
}
