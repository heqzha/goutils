package concurrency

import(
	"log"
	"time"
)

func TaskRunPeriodic(f func() time.Duration, task_name string, defaultInterval time.Duration) {
	if defaultInterval < time.Second {
		defaultInterval = time.Second
	}

	go func() {
		for {
			func() {
				defer func() {
					i := recover()
					log.Println("task panic: ", i, task_name)
				}()
				for {
					if interval := f(); interval > 0 {
						time.Sleep(interval)
					} else {
						time.Sleep(defaultInterval)
					}
				}
			}()
			time.Sleep(defaultInterval)
		}
	}()
}
