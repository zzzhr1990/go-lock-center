package main

import (
	"log"
	"time"

	"github.com/zzzhr1990/go-lock-center/locker"
)

func main() {
	//
	cfg := &locker.Config{
		RedisAddress: []string{
			"redis://:papaya@127.0.0.1:6379/0",
		},
	}
	ec, err := locker.CreateNew(cfg)
	if err != nil {
		log.Fatalf("Init locker failed %v", err)
	} else {
		_, err = ec.LockForKeyWithNoRetry("papayax", time.Second*10)
		if err != nil {
			log.Fatalf("Cannot set 'papaya' lock %v", err)
			return
		}
		log.Printf("TRY_GET_LOCK_1")
		_, err = ec.LockForKeyWithNoRetry("papayax", time.Second)
		log.Printf("RESPONSE_GET_LOCK_1")
		if err != nil {
			log.Printf("Yes you can't set 'papayax' lock wait for 20 sec")
			time.Sleep(time.Second * 20)
			log.Printf("TRY_GET_LOCK_2")
			lcx, err := ec.LockForKeyWithNoRetry("papayax", time.Second)
			log.Printf("TRY_GET_LOCK_RESPONSE")
			if err != nil {
				log.Fatalf("Cannot set 'papaya' lock 2 %v", err)
				return
			}
			log.Printf("Lock expired, try next")

			lcx.Unlock()
			lcxx, err := ec.LockForKeyWithNoRetry("papayax", time.Second)
			if err != nil {
				log.Fatalf("Cannot set 'papaya' lock 3 %v", err)
				return
			}
			lcxx.Unlock()
			log.Printf("All test ok")
			return
		}
		log.Fatalln("'papaya' lock not exists")

	}
}
