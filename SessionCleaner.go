package main

import "time"

func SessionCleaner(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Minute * 10)
			db.CleanUserSession()
		}
	}
}
