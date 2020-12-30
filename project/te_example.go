package project

import (
	"TopEngine/common"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"time"
)

//全局设定
var netSpeedTX float64 = 1000 / 8 * 1024 * 1024 //下行速率 单位bps
var netSpeedRX float64 = 1000 / 8 * 1024 * 1024 //上行速率 单位bps

var cpuSegment = make([]int64, 3600)
var memSegment = make([]int64, 3600)
var hddSegment = make([]int64, 3600)
var diskIoSegment = make([]int64, 3600)
var netIoRxSegment = make([]int64, 3600)
var netIoTxSegment = make([]int64, 3600)

type resPerfJSON struct {
	GopsutilData []int64 `json:"gopsutil_data"`
}

func perfDot() {
	for {
		//cpuDot
		cpuMethod, _ := cpu.Percent(1*time.Second, false)
		cpuSegment = append(cpuSegment, int64(cpuMethod[0]))
		cpuSegment = cpuSegment[1:]

		//memDot
		memMethod, _ := mem.VirtualMemory()
		memSegment = append(memSegment, int64(memMethod.UsedPercent))
		memSegment = memSegment[1:]

		//diskDot
		hddMethod, _ := disk.Usage("/")
		hddSegment = append(hddSegment, int64(hddMethod.UsedPercent))
		hddSegment = hddSegment[1:]
	}
}

func diskDot() {
	for {
		start := make(map[string]uint64)
		//手动填写一下要聚合数据的磁盘标签
		diskIo, _ := disk.IOCounters()
		for name, value := range diskIo {
			start[name] = value.IoTime
		}

		time.Sleep(1 * time.Second)

		//手动填写一下要聚合数据的磁盘标签
		diskIo, _ = disk.IOCounters()
		var valTmp int64
		for name, value := range diskIo {
			valTmp += int64(value.IoTime - start[name])
		}
		valTmp /= int64(len(diskIo))

		diskIoSegment = append(diskIoSegment, int64(float64(valTmp)/1000.00*100))
		diskIoSegment = diskIoSegment[1:]
	}
}

func HandleRoute(dr *common.DynamicRoute) {
	var err error
	go perfDot()
	go diskDot()
	dr.AddRoute("/platform/monitor/cpu", []string{"GET"}, func(req common.Request) (res common.Response) {
		res.Header = make(map[string]string)
		res.Header["Content-Type"] = "application/json"
		resJSON := &resPerfJSON{
			GopsutilData: cpuSegment,
		}
		if res.Body, err = json.Marshal(resJSON); err != nil {
			fmt.Println("ERROR")
		}
		return
	})

	dr.AddRoute("/platform/monitor/mem", []string{"GET"}, func(req common.Request) (res common.Response) {
		res.Header = make(map[string]string)
		res.Header["Content-Type"] = "application/json"
		resJSON := &resPerfJSON{
			GopsutilData: memSegment,
		}
		if res.Body, err = json.Marshal(resJSON); err != nil {
			fmt.Println("ERROR")
		}
		return
	})

	dr.AddRoute("/platform/monitor/disk", []string{"GET"}, func(req common.Request) (res common.Response) {
		res.Header = make(map[string]string)
		res.Header["Content-Type"] = "application/json"
		resJSON := &resPerfJSON{
			GopsutilData: hddSegment,
		}
		if res.Body, err = json.Marshal(resJSON); err != nil {
			fmt.Println("ERROR")
		}
		return
	})

	dr.AddRoute("/platform/monitor/disk_io", []string{"GET"}, func(req common.Request) (res common.Response) {
		res.Header = make(map[string]string)
		res.Header["Content-Type"] = "application/json"
		resJSON := &resPerfJSON{
			GopsutilData: diskIoSegment,
		}
		if res.Body, err = json.Marshal(resJSON); err != nil {
			fmt.Println("ERROR")
		}
		return
	})
}
