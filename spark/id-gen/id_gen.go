package id_gen

import (
	"errors"
	"github.com/sony/sonyflake"
)

var (
	flake *sonyflake.Sonyflake
)

func Init(machineId uint16) {
	settings := sonyflake.Settings{}
	settings.MachineID = func() (u uint16, err error) {
		return machineId, nil
	}
	flake = sonyflake.NewSonyflake(settings)
}

func GenId() (id uint64, err error) {
	if flake == nil {
		err = errors.New("sony flake not init")
		return
	}
	return flake.NextID()
}
