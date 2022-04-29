package main

import (
	"time"
)

func (f *Frontend) radiuscheck() {

	time.Sleep(time.Second * 30)

	for {
		rdcheck = ""
		go func() {
			time.Sleep(time.Second * 5)
			if rdcheck == "" {
				Logger.Println("[HC-LOG]Radius check failed")
				time.Sleep(time.Second * 5)
				if rdcheck == "" {
					Logger.Println("[HC-LOG]server is stopping")
					f.server.Listener.Close()
					rdDown = true
					Logger.Println("[HC-LOG]rdDown is now true")
				}
			}
		}()
		rdcheck = beckend(rdcuser, rdcrealm, true)

		if rdDown && rdcheck != "" {
			Logger.Println("[HC-LOG]Radius check Recovered")
			rdDown = false
			Logger.Println("[HC-LOG]rdDown is now false")
			go func() {
				Logger.Println("[HC-LOG]server is starting")
				err := f.server.ListenAndServe(f.serverAddr)
				Logger.Println(err)
			}()
		}

		time.Sleep(time.Second * 15)
	}
}
