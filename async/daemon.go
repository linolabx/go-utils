package async

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Daemon struct {
	Interval    time.Duration
	ExitSignals []os.Signal
	PreTask     func() error
	PostTask    func()
	Logger      func(msg string)

	lock sync.Mutex
}

func (this *Daemon) initialize() {
	if this.Interval == 0 {
		this.Interval = time.Microsecond
	}

	if len(this.ExitSignals) == 0 {
		this.ExitSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
	}

	if this.PreTask == nil {
		this.PreTask = func() error { return nil }
	}

	if this.PostTask == nil {
		this.PostTask = func() {}
	}

	if this.Logger == nil {
		this.Logger = func(msg string) {}
	}
}

func (this *Daemon) Exec(task func()) error {
	if lock := this.lock.TryLock(); lock == false {
		panic("each daemon can only run once until process end")
	}

	this.initialize()

	this.Logger("run post task")
	if err := this.PreTask(); err != nil {
		return err
	}

	this.Logger("start main loop...")

	ticker := time.NewTicker(this.Interval)

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, this.ExitSignals...)

	mu := &sync.Mutex{}
	exitTimes := 0

	for {
		select {
		case <-quit:
			if exitTimes > 0 {
				this.Logger("force exit")
				// exit 0 to avoid container restart
				os.Exit(0)
				return nil
			}
			exitTimes++

			go func() {
				this.Logger("graceful shutdown")

				this.Logger("stop tinker...")
				ticker.Stop()

				this.Logger("wait for running task...")
				mu.Lock()

				this.Logger("run post task")
				this.PostTask()

				this.Logger("done")
				os.Exit(0)
				return
			}()
		case <-ticker.C:
			if mu.TryLock() == false {
				this.Logger("prev task takes too long, skip this round")
				continue
			}

			go func() {
				defer mu.Unlock()
				task()
			}()
		}
	}
}
