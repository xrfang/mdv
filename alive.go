package main

import (
	"os"
	"sync"
	"time"
)

var (
	mx    sync.Mutex
	ticks int
)

func refreshTickCount() {
	mx.Lock()
	defer mx.Unlock()
	ticks = 0
}

func init() {
	go func() {
		for {
			time.Sleep(time.Second)
			func() {
				mx.Lock()
				defer mx.Unlock()
				ticks++
				if ticks > 99 {
					os.Exit(0)
				}
			}()
		}
	}()
}
