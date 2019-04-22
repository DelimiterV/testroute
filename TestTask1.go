package main

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	nats "github.com/nats-io/go-nats"
)

var nc *nats.Conn

func init() {
	var err error
	nc, err = nats.Connect("nats://192.168.99.100:4222", nats.PingInterval(20*time.Second))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Nats is here!")
	}
}

//====================================================
var ppu string = `Server: Delimiter-lab
User-Agent: Technolab/1.0 Highload/1.0
Connection: close
Content-Type: application/json`
var ppu2 string = `Content-Length:`
var ppu3 string = "\n\n"

func CreatePacket(kod string, mcontent string) (rez string) {
	rez = "HTTP/1.1 " + kod + "\n"
	rez += ppu
	g := len(mcontent)
	if g > 0 {
		rez += "\n" + ppu2
		rez += strconv.Itoa(g)
		rez += ppu3
		rez += mcontent
	} else {
		rez += ppu3
	}
	return
}

//----------------------------------------------------
const SENDLIMIT = 8000
const RECVLIMIT = 8000
const WIDTH = 2000

func handleConnection(ind int) {
	var i, j, z, n, m, u, zu, ku, rlen int
	var err error
	var rstr string
	var content string
	var id, params string

	ex := 0

	for ex != 1 {
		select {
		case <-Nch[ind].mych:
			rlen, err = Nch[ind].conn.Read(Nch[ind].bufR)

			if err == nil && rlen > 0 {
				//fmt.Println(string(Nch[ind].bufR[:rlen-1]))
				zu, ku, z, u = 0, 0, 0, 0
				id = ""
				params = ""
				if rlen > 0 && ex != 1 {
					for j = 0; j < rlen && Nch[ind].bufR[j] != ' ' && Nch[ind].bufR[j] != '\n'; j++ {
						zu += int(Nch[ind].bufR[j]) << uint(j)
					}
					j++
					z = 0
					for i = j; i < rlen && Nch[ind].bufR[i] != ' ' && Nch[ind].bufR[i] != '\n'; i++ {
						if Nch[ind].bufR[i] == '?' {
							z = 1
						}
						if z == 0 {
							if Nch[ind].bufR[i] < '0' || Nch[ind].bufR[i] > '9' {
								u++
								ku += int(Nch[ind].bufR[i]) << uint((u)%4)
							} else {
								id += string(Nch[ind].bufR[i : i+1])
							}
						} else {
							params += string(Nch[ind].bufR[i : i+1])
						}
					}
					fmt.Println(zu, ku)
					switch zu {
					case 545: //GET
						switch ku {
						case 2819: //news/id
							z := GetNewsByID(ind, id)
							_, _ = Nch[ind].conn.Write([]byte(z))
							//							AccountsFilter(Nch[ind].conn, ind, params)
						default:
							rstr = CreatePacket("404 Not Found", "")
							Nch[ind].bufS = []byte(rstr)
							_, _ = Nch[ind].conn.Write(Nch[ind].bufS[:len(Nch[ind].bufS)-1])

						}
						rstr = CreatePacket("200 OK", "{}")
						Nch[ind].bufS = []byte(rstr)
						_, _ = Nch[ind].conn.Write(Nch[ind].bufS[:len(Nch[ind].bufS)-1])

					case 1242: //POST
						content = ""
						m = 0
						for n = i; n < rlen-1 && m == 0; n++ {

							if Nch[ind].bufR[n] == '{' {
								m = 1
								content = string(Nch[ind].bufR[n:rlen])
							}
						}
						fmt.Println(content)
						switch ku {
						case 3073: //add/news

							//							AccountsNew(Nch[ind].conn, ind, content)
							AddNews(ind, content)
						default:
							rstr = CreatePacket("404 Not Found", "")
							Nch[ind].bufS = []byte(rstr)
							_, _ = Nch[ind].conn.Write(Nch[ind].bufS[:len(Nch[ind].bufS)-1])

						}

					default:
						rstr = CreatePacket("405 Method Not Allowed ", "")
						Nch[ind].bufS = []byte(rstr)
						_, _ = Nch[ind].conn.Write(Nch[ind].bufS[:])

					}

				}

			}

		}
		Nch[ind].conn.Close()
		Nch[ind].mflag.Lock()
		Nch[ind].flag = 0
		Nch[ind].mflag.Unlock()
	}

}

const PORT = ":80"
const NLINE = 40

type lnnn struct {
	mych   chan int
	fch    chan uint8
	conn   net.Conn
	flag   int
	mflag  sync.Mutex
	dbmess []byte
	bufR   []byte
	bufS   []byte
}

var Nch [NLINE]lnnn

func main() {
	var counter int
	var i int

	var ex int = 0
	//	var tmm time.Time
	for i = 0; i < NLINE; i++ {
		Nch[i].mych = make(chan int)
		Nch[i].fch = make(chan uint8)
		Nch[i].flag = 0
		Nch[i].dbmess = make([]byte, WIDTH)
		Nch[i].bufR = make([]byte, RECVLIMIT)
		Nch[i].bufS = make([]byte, SENDLIMIT)
		go handleConnection(i)
	}
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println("Bind error:", err.Error())
	} else {
		for {
			conn, err := ln.Accept()
			if err != nil {
				fmt.Println("Accept error:", err.Error())
			} else {
				ex = 0
				for ex == 0 {
					if counter < NLINE-2 {
						counter++
					} else {
						counter = 0
					}
					Nch[counter].mflag.Lock()
					if Nch[counter].flag == 0 {
						//fmt.Print(".")
						Nch[counter].flag = 1
						Nch[counter].conn = conn
						Nch[counter].mych <- counter
						ex = 1
					}
					Nch[counter].mflag.Unlock()
				}
			}
		}
	}
	ln.Close()
}
