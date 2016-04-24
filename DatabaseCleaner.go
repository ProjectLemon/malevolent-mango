package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

//SessionCleaner wakes up every ten minutes and
//removes inactive sessions from database
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

//ImageCleaner wakes up every 24h and
//removes unused images from the server
func ImageCleaner(quit chan bool) {
	for {
		select {
		case <-quit:
			return
		default:
			time.Sleep(time.Second * 10)
			profileHeaders, err := filepath.Glob("www/img/profile-headers/*")
			if err != nil {
				fmt.Println(err)
				return
			}

			profileIcons, err := filepath.Glob("www/img/profile-icons/*")
			if err != nil {
				fmt.Println(err)
				return
			}

			for i := 0; i < len(profileHeaders); i++ {
				if isImg(profileHeaders[i]) {
					inDB, err := db.UniversalLookup(profileHeaders[i])
					if err != nil {
						return
					}
					if !inDB {
						os.Remove(profileHeaders[i])
					}
				}
			}
			for i := 0; i < len(profileIcons); i++ {
				if isImg(profileIcons[i]) {
					inDB, err := db.UniversalLookup(profileIcons[i])
					if err != nil {
						return
					}
					if !inDB {
						os.Remove(profileIcons[i])
					}
				}
			}
		}
	}
}

func isImg(fileName string) bool {
	extension := fileName[len(fileName)-4 : len(fileName)]
	return extension == ".jpg" || extension == ".png" || extension == ".gif"
}
