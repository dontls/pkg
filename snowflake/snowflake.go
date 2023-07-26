package snowflake

import (
	"log"
	"time"

	"github.com/bwmarrin/snowflake"
)

var _snowflake *snowflake.Node = nil

// New 创建实例
func New(startTime string, machineID int64) {
	var st time.Time
	// 格式化 1月2号下午3时4分5秒  2006年
	st, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		log.Fatalln(err)
	}
	snowflake.Epoch = st.UnixNano() / 1e6
	if _snowflake, err = snowflake.NewNode(machineID); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	New("2021-12-03", 0x01)
}

func NextID() uint {
	return uint(_snowflake.Generate().Int64())
}
