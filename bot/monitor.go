package bot

import "time"

type Monitor struct {
}

func (m *Monitor) start() {
	for {
		logs("Monitor loop tick. ⏱️")

		time.Sleep(time.Second * MonitorTick)
	}
}

func initMonitor() *Monitor {
	m := &Monitor{}
	go m.start()
	return m
}
