package main

import (
	"fmt"
	"time"
	"unicode"

	. "github.com/immesys/nb"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

const Interval = 200 * time.Millisecond
const Version = "1.0"

func main() {
	fmt.Println("Version", Version)
	NB("nbtop.event", "type", "started")
	defer func() {
		NB("nbtop.event", "type", "exit")
		NBClose()
	}()
	then := time.Now()
	for {
		time.Sleep(then.Add(Interval).Sub(time.Now()))
		then = time.Now()
		doMemory()
		doCPU()
		doDocker()
		doDisk()
		doNetwork()
	}
}

func doMemory() {
	t, err := mem.VirtualMemory()
	if err != nil {
		NB("nbtop.error", "error", err.Error(), "loc", "mem")
		panic(err)
	}
	NB("nbtop.mem",
		"ts_available", float64(t.Available),
		"ts_used", float64(t.Used),
		"ts_free", float64(t.Free),
	)
}

func doCPU() {
	t, err := cpu.Times(false)
	if err != nil {
		NB("nbtop.error", "error", err.Error(), "loc", "cpu")
		panic(err)
	}
	NB("nbtop.cpu",
		"cpu", t[0].CPU,
		"ts_user", t[0].User,
		"ts_system", t[0].System,
		"ts_idle", t[0].Idle,
		"ts_nice", t[0].Nice,
		"ts_iowait", t[0].Iowait,
		"ts_irq", t[0].Irq,
		"ts_softirq", t[0].Softirq,
		"ts_steal", t[0].Steal,
		"ts_guest", t[0].Guest,
		"ts_guestnice", t[0].GuestNice,
		"ts_stolen", t[0].Stolen,
		"ts_total", t[0].Total(),
	)
}

func doDocker() {

}

func doDisk() {
	t, err := disk.IOCounters()
	if err != nil {
		NB("nbtop.error", "error", err.Error(), "loc", "disk")
		panic(err)
	}

	for disk, stats := range t {
		if unicode.IsNumber(rune(disk[len(disk)-1])) {
			continue
		}
		args := []interface{}{}
		args = append(args, "mt_disk", disk)
		args = append(args, "ts_readcount", float64(stats.ReadCount))
		args = append(args, "ts_mergedreadcount", float64(stats.MergedReadCount))
		args = append(args, "ts_writecount", float64(stats.WriteCount))
		args = append(args, "ts_mergedwritecount", float64(stats.MergedWriteCount))
		args = append(args, "ts_readbytes", float64(stats.ReadBytes))
		args = append(args, "ts_writebytes", float64(stats.WriteBytes))
		args = append(args, "ts_readtime", float64(stats.ReadTime))
		args = append(args, "ts_writetime", float64(stats.WriteTime))
		args = append(args, "ts_iopsinprogress", float64(stats.IopsInProgress))
		args = append(args, "ts_iotime", float64(stats.IoTime))
		NB("nbtop.disk", args...)
	}
}

func doNetwork() {
	t, err := net.IOCounters(true)
	if err != nil {
		NB("nbtop.error", "error", err.Error(), "loc", "net")
		panic(err)
	}
	for _, stat := range t {
		NB("nbtop.net",
			"mt_nic", stat.Name,
			"ts_sentbytes", stat.BytesSent,
			"ts_recvbytes", stat.BytesRecv,
			"ts_sentpackets", stat.PacketsSent,
			"ts_recvpackets", stat.PacketsRecv,
			"ts_errin", stat.Errin,
			"ts_errout", stat.Errout,
			"ts_dropin", stat.Dropin,
			"ts_dropout", stat.Dropout,
		)
	}
}
