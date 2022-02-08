package main

import (
	"time"
	"fmt"
	"context"
	log "github.com/sirupsen/logrus"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/DomesticMoth/confer"
)

const DEFAULT_GLOBAL_PATH string = "/etc/syswatcher/config.toml"

type Conf struct{
	Delay uint64
	Addr string
	Database string
	Username string
	Password string
	Table string
}

func calcSingleCoreUsage(curr, prev linuxproc.CPUStat) float32 {
  PrevIdle := prev.Idle + prev.IOWait
  Idle := curr.Idle + curr.IOWait
  PrevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
  NonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal
  PrevTotal := PrevIdle + PrevNonIdle
  Total := Idle + NonIdle
  totald := Total - PrevTotal
  idled := Idle - PrevIdle
  CPU_Percentage := (float32(totald) - float32(idled)) / float32(totald)
  return CPU_Percentage
}

func getRam() uint64{
	mem := ReadMemoryStats()
	total := mem.MemTotal
	free := mem.MemAvailable
	return (total-free)/(total/100)
}

func getCpu() linuxproc.CPUStat {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil { log.Fatal(err) }
	return stat.CPUStatAll
}

func main(){
	var conf Conf
	conf.Delay = 1000 // Default 1 ceond

	err := confer.LoadConfig([]string{DEFAULT_GLOBAL_PATH}, &conf)	
	if err != nil { log.Fatal(err) }

	ctx := context.Background()
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{conf.Addr},
		Auth: clickhouse.Auth{
			Database: conf.Database,
			Username: conf.Username,
			Password: conf.Password,
		},
		DialTimeout:     time.Second,
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Hour,
	})
	if err != nil { log.Fatal(err) }
	
	last := getCpu()
	for {
		stat := getCpu()
		load := uint64(calcSingleCoreUsage(stat, last)*100)
		last = stat
		if load < 100 {
			timestamp := time.Now().Unix()
			query := fmt.Sprintf("INSERT INTO %s VALUES (%d,%d,%d,0)", conf.Table, timestamp, load, getRam())
			//log.Info(query)
			err := conn.AsyncInsert(ctx, query, false)
			if err != nil { log.Fatal(err) }
		}
		time.Sleep(time.Duration(conf.Delay) * time.Millisecond)
	}
}
