package main

import (
	"time"
)

func secondarycheck() {

	time.Sleep(time.Second * 30)

	for {
		rdcheck = ""
		go func() {
			time.Sleep(time.Second * 5)
			if rdcheck == "" {
				Logger.Println("[HC-LOG-SEC]Radius check failed")
				time.Sleep(time.Second * 5)
				if rdcheck == "" {
					rdDown = true
					Logger.Println("[HC-LOG-SEC]rdDown is now true")
				}
			}
		}()
		rdcheck = beckend(rdcuser, rdcrealm, true)
		//Logger.Println(rdcheck)

		if rdDown && rdcheck != "" {
			Logger.Println("[HC-LOG-SEC]Radius check Recovered")
			rdDown = false
			Logger.Println("[HC-LOG-SEC]rdDown is now false")

		}

		time.Sleep(time.Second * 15)

		//Logger.Println("[HC-LOG-SEC]recheck process started")
	}
}
