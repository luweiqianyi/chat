package api

import "fmt"

const (
	PrefixRoomsKey = "com.company.app.service:rooms"

	PrefixGroupKey     = "com.company.app.service" // eg.  com.alibaba.TMall.chat TODO replace in runtime
	FormatGroup        = "%v:group:%v"
	FieldGroupBaseInfo = "GroupBaseInfo"
)

// RedisCreateRoomData one user can only create one live chat room
type RedisCreateRoomData struct {
	RoomID          string `json:"roomID"`
	CreatorID       string `json:"creatorID"`
	CreateTimestamp int64  `json:"createTimestamp"`
}

func GenRedisKey(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func RoomsRedisKey() string {
	return PrefixRoomsKey
}
