package snowflake

import (
	"time"

	"github.com/busyfree/shorturl-go/util/conf"
	"github.com/sony/sonyflake"
	"github.com/spf13/cast"
)

var (
	sf        *sonyflake.Sonyflake
	startTime time.Time
	machineID uint16
)

func init() {
	var st sonyflake.Settings
	st.MachineID = GetMachineId
	startTime = conf.GetTime("SNOWFLAKE_START_TIME")
	if !startTime.IsZero() {
		st.StartTime = startTime
	}
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		panic("sonyflake not created")
	}
}

func NextId() (nextId int64, err error) {
	var newId uint64
	newId, err = sf.NextID()
	if err != nil {
		return
	}
	nextId = int64(newId)
	return
}

//UInt16 - [0 : 65535]
func GetMachineId() (uint16, error) {
	machineID = cast.ToUint16(conf.GetInt("SNOWFLAKE_MACHINE_ID"))
	return machineID, nil
}

func CheckMachineID(mchId uint16) bool {
	if mchId != machineID {
		return false
	}
	return true
}
