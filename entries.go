package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var dcount int
var lmut sync.Mutex

func GetNewsByID(ind int, idnews string) string {
	var r string = ""
	lmut.Lock()
	dstr := strconv.Itoa(dcount)
	dcount++
	lmut.Unlock()
	nc.Publish("id."+dstr, []byte(idnews))
	sub, err := nc.SubscribeSync("i." + dstr)
	if err == nil {
		m, er := sub.NextMsg(10 * time.Millisecond)
		if er == nil {
			r = m.Reply
		}
	}
	return r
}

func AddNews(myid int, content string) {
	lmut.Lock()
	dstr := strconv.Itoa(dcount)
	dcount++
	lmut.Unlock()
	nc.Publish("new."+dstr, []byte(content))
	sub, err := nc.SubscribeSync("n." + dstr)
	if err == nil {
		m, er := sub.NextMsg(10 * time.Millisecond)
		if er == nil {
			fmt.Println(m.Reply)
		}
	}

}
