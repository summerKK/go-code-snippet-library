package scheduler

import (
	"fmt"
	"sync"
)

const (
	// 未初始化
	SCHED_STATUS_UINIT = iota
	// 初始化中
	SCHED_STATUS_INITING
	// 已初始化
	SCHED_STATUS_INITED
	//  启动中
	SCHED_STATUS_STARTING
	// 已启动
	SCHED_STATUS_STARTED
	// 停止中
	SCHED_STATUS_STOPPING
	// 已停止
	SCHED_STATUS_STOPPED
)

func checkStatus(oldStatus, wantedStatus Status, locker sync.Locker) (err error) {
	if locker != nil {
		locker.Lock()
		defer locker.Unlock()
	}
	// 处于正在初始化、正在启动或正在停止状态时,不能从外部改变状态
	switch oldStatus {
	case SCHED_STATUS_INITING:
		err = genError("the scheduler is being initialized!")
	case SCHED_STATUS_STARTING:
		err = genError("the scheduelr is being started!")
	case SCHED_STATUS_STOPPING:
		err = genError("the scheduelr is being stoped!")
	}
	if err != nil {
		return
	}
	if oldStatus == SCHED_STATUS_UINIT && (wantedStatus == SCHED_STATUS_STARTING || wantedStatus == SCHED_STATUS_STOPPING) {
		err = genError("the scheduler has not yet been initizled!")
		return
	}
	switch wantedStatus {
	case SCHED_STATUS_INITING:
		switch oldStatus {
		case SCHED_STATUS_STARTED:
			err = genError("the scheduler has been started!")
		}
	case SCHED_STATUS_STARTING:
		switch oldStatus {
		case SCHED_STATUS_UINIT:
			err = genError("the scheduler has not been initialized!")
		case SCHED_STATUS_STARTED:
			err = genError("the scheduler has been started!")
		}
	case SCHED_STATUS_STOPPING:
		if oldStatus != SCHED_STATUS_STARTED {
			err = genError("the scheduler has not been started!")
		}
	default:
		errMsg :=
			fmt.Sprintf("unsupported wanted status for check! (wantedStatus: %d)",
				wantedStatus)
		err = genError(errMsg)
	}
	return
}

// GetStatusDescription 用于获取状态的文字描述。
func GetStatusDescription(status Status) string {
	switch status {
	case SCHED_STATUS_UINIT:
		return "uninitialized"
	case SCHED_STATUS_INITING:
		return "initializing"
	case SCHED_STATUS_INITED:
		return "initialized"
	case SCHED_STATUS_STARTING:
		return "starting"
	case SCHED_STATUS_STARTED:
		return "started"
	case SCHED_STATUS_STOPPING:
		return "stopping"
	case SCHED_STATUS_STOPPED:
		return "stopped"
	default:
		return "unknown"
	}
}
