package resourcemonitor

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"time"
)

const (
	sheet1 = "Sheet1"
)

var (
	titles = []string{
		"Time",

		"%Cpu",
		"UsedMem",

		// CPU百分比
		"%Cpu_rabbitmq-server",
		"%Cpu_redis-server",
		"%Cpu_mysqld",
		"%Cpu_mongod",

		// 内存使用量
		"Mem_rabbitmq-server",
		"Mem_redis-server",
		"Mem_mysqld",
		"Mem_mongod",
	}

	services = map[string]int{
		"rabbitmq-server": 1,
		"redis-server":    2,
		"mysqld":          3,
		"mongod":          4,
	}
)

type ResourceMonitor struct {
	Times int
}

func (r *ResourceMonitor) Monitor() {
	processList, _ := process.Processes()

	excel := excelize.NewFile()
	for i, v := range titles {
		if 65+i > 90 {
			excel.SetCellValue(sheet1, fmt.Sprintf("%c%c%d", 65, 39+i, 1), v)

		} else {
			excel.SetCellStr(sheet1, fmt.Sprintf("%c%d", 65+i, 1), v)
		}
	}

	ticker := time.NewTicker(time.Second * 1)
	defer ticker.Stop()
	line := 2
	for line < r.Times {
		select {
		case <-ticker.C:
			// A列时间
			now := time.Now()
			excel.SetCellValue(sheet1, fmt.Sprintf("%v%d", "A", line), now.Format("2006-01-02 15:04:05"))

			// B列CPU总占比
			totalRate, err := getCpu()
			if err != nil {
				fmt.Println(err)
			}
			excel.SetCellValue(sheet1, fmt.Sprintf("%v%d", "B", line), totalRate)

			// C列内存总使用量
			usedMem, err := getMem()
			if err != nil {
				fmt.Println(err)
			}
			excel.SetCellValue(sheet1, fmt.Sprintf("%v%d", "C", line), usedMem)

			// D-O列各服务CPU占比；P-AI列各服务内存使用量
			for _, proc := range processList {
				pName, _ := proc.Name()
				if idx, ok := services[pName]; ok {
					cpuRate, err := proc.CPUPercent()
					if err != nil {
						fmt.Println(err)
					}
					excel.SetCellValue(sheet1, fmt.Sprintf("%c%d", 67+idx, line), cpuRate)

					memory, err := proc.MemoryInfo()
					if err != nil {
						fmt.Println(err)
					}
					if memory != nil {
						if 70+idx > 90 {
							// 超出Z列，从AA列开始
							excel.SetCellValue(sheet1, fmt.Sprintf("%c%c%d", 65, 57+idx, line), memory.RSS/(1<<20))
						} else {
							excel.SetCellValue(sheet1, fmt.Sprintf("%c%d", 83+idx, line), memory.RSS/(1<<20))
						}
					}
				}
			}
			line++
		}
	}
	_ = excel.SaveAs("/home/user/gaowei/zdb_resources/top" + time.Now().Format("20060102-15") + ".xlsx")
}

func getCpu() (float64, error) {
	rates, err := cpu.Percent(time.Second*1, false)
	if err != nil {
		return 0, err
	}
	return rates[0], nil
}

func getMem() (uint64, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return vmem.Used / (1 << 20), err
}
