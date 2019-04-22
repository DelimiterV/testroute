// Models.go
package Models

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	//_ "github.com/go-sql-driver/mysql"
)

func makeTimestamp() int {
	var k time.Time = time.Now()
	r := (k.Year()-2000)*1000000 + k.Day()*1000000 + k.Hour()*10000 + k.Minute()*100 + k.Second()
	return r
}
func Start() {
	go ManagerDB(Cfg)
}

type Dbcfg struct {
	Kind       string
	Transport  string
	ServerIP   string
	ServerPort string
	DbName     string
	User       string
	Password   string
}

var Cfg Dbcfg

func ini_channels() {
	/*	LoginUAskChan = make(chan []string)
		LoginURezChan = make(chan AccM)
		AddUAskChan = make(chan AccM)
		AddURezChan = make(chan int)
		GetUAskChan = make(chan int)
		GetURezChan = make(chan AccM)
		CheckUAskChan = make(chan []string)
		CheckURezChan = make(chan int)
		IfExistAskChan = make(chan string)
		IfExistRezChan = make(chan int)
		UpdateUAskChan = make(chan Zzi)
		UpdateURezChan = make(chan int)
		SavePassAskChan = make(chan []string)
		SavePassRezChan = make(chan int)
		ClearContentAskChan = make(chan int)
		ClearContentRezChan = make(chan int)
		LoadCourseTextAskChan = make(chan L_D_T)
		LoadCourseTextRezChan = make(chan int)
		GetCourseTextAskChan = make(chan L_D_T)
		GetCourseTextRezChan = make(chan string) */
}

func abs(num int) int {
	if num < 0 {
		return num * (-1)
	} else {
		return num
	}
}

func ManagerDB(cfg Dbcfg) {

	var ex int = 0
	var ex2 int = 0
	//var i, j int

	var TimeStamp time.Time

	ini_channels()

	for ex == 0 {

		strConnect := "user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.DbName + " sslmode=disable"
		fmt.Println(strConnect)
		db, err := sql.Open(cfg.Kind, strConnect)
		if err == nil {
			err = db.Ping()
			if err == nil {
				//*************************************************
				_, err = db.Exec("DROP TABLE IF EXISTS newstable ;")
				if err != nil {
					fmt.Println(err)
				}
				_, err = db.Exec("CREATE TABLE newstable " +
					`("id" SERIAL PRIMARY KEY,` +
					`"date" varchar(50), "news" varchar(250))`)
				if err != nil {
					fmt.Println(err)
				}
				TimeStamp = time.Now().UTC()
				ex2 = 0
				for ex2 == 0 {
					select {
					/* case u := <-GetCourseTextAskChan:
						go func() {
							rez := ""
							rows, err := db.Query("SELECT Content FROM course WHERE idCourse=? AND lang=?", u.mid, u.lang)
							if err == nil {
								for rows.Next() {
									err := rows.Scan(&rez)
									if err != nil {
										fmt.Println(err)
									}
								}
								rows.Close()
							} else {
								rez = "-2"
							}
							GetCourseTextRezChan <- rez
						}()
					case u := <-LoadCourseTextAskChan:
						r := 0
						stmt, err := db.Prepare("INSERT course SET idCourse=?,lang=?,content=?")
						if err == nil {
							res, err := stmt.Exec(u.mid, u.lang, u.text)
							if err == nil {
								id, errp := res.LastInsertId()
								if errp == nil {
									r = int(id)
								} else {
									r = -1
								}
							} else {
								r = -2
							}
							stmt.Close()
						} else {
							r = -3
						}
						LoadCourseTextRezChan <- r
					case u := <-ClearContentAskChan:
						fmt.Println("Models:", u)
						r := 0
						stmt, err1 := db.Prepare("DELETE FROM course ;")
						if err1 == nil {
							_, e := stmt.Exec()
							if e == nil {
								r = 1
							} else {
								fmt.Println(e)
								r = -1
							}
							stmt.Close()
						} else {
							fmt.Println(err1)
							r = -2
						}
						ClearContentRezChan <- r
					case u := <-AddUAskChan:
						go func() {
							u.DatFrom = ""
							u.DatTo = ""
							r := 0
							stmt, err := db.Prepare("INSERT users SET Nicname=?,Name=?,SecName=?,Email=?,Password=?,TTL=?,Kind=?,Token=?,Policy=?,FA=?,DateFrom=?,DateTo=?")
							if err == nil {
								res, err := stmt.Exec(u.NickName, u.Name, u.Family, u.Email, u.Password, u.Telephone, u.Ttl, u.Kind, u.Token, u.Policy, u.Fa, u.DatFrom, u.DatTo)
								if err == nil {
									id, err := res.LastInsertId()
									if err == nil {
										r = int(id)
									} else {
										r = -1
									}
								} else {
									fmt.Println(err)
									r = -2
								}
								stmt.Close()
							} else {
								fmt.Println(err)
								r = -3
							}
							AddURezChan <- r
						}()
					case u := <-SavePassAskChan:
						go func() {
							var ru int
							if len(u) > 1 {
								str := "UPDATE users SET "
								str += "Password='" + u[1] + "' "
								str += "WHERE Email ='" + u[0] + "' ;"
								stm, err := db.Prepare(str)
								if err == nil {
									_, err2 := stm.Exec()
									if err2 == nil {
										ru = 1
									} else {
										ru = -1
									}
									stm.Close()
								} else {
									ru = -2
								}
							}
							SavePassRezChan <- ru
						}()
					case u := <-CheckUAskChan:
						go func() {
							if len(u) > 1 {
								rows, err := db.Query("SELECT id FROM users WHERE Email=? AND Password=?", u[0], u[1])
								if err == nil {
									id := -1
									for rows.Next() && id == -1 {
										err := rows.Scan(&id)
										if err != nil {
											fmt.Println(err)
										}
									}
									CheckURezChan <- id
									rows.Close()
								} else {
									CheckURezChan <- -2
								}
							} else {
								CheckURezChan <- -3
							}
						}()
					case u := <-GetUAskChan:
						go func() {
							var rez AccM
							rez.Id = -1
							rows, err := db.Query("SELECT * FROM users WHERE id=?", u)
							if err == nil {

								for rows.Next() {
									err := rows.Scan(&rez.Id, &rez.NickName, &rez.Name, &rez.Family, &rez.Email, &rez.Password, &rez.Telephone, &rez.Ttl, &rez.Kind, &rez.Token, &rez.Policy, &rez.Fa, &rez.DatFrom, &rez.DatTo)
									if err != nil {
										fmt.Println(err)
									}
								}
								GetURezChan <- rez
								rows.Close()
							} else {
								rez.Id = -2
								GetURezChan <- rez
							}
						}()
					case r := <-UpdateUAskChan:
						go func() {
							var ru int = 0
							//UPDATE user SET name='xxx',family='yyyy' WHERE id =7000;
							str := "UPDATE users SET "
							str += "Nicname='" + r.nickname + "',"
							str += "Name='" + r.name + "',"
							str += "SecName='" + r.family + "',"
							str += "telephon='" + r.telephone + "',"
							str += "DateFrom='" + r.dateFrom + "',"
							str += "DateTo='" + r.dateTo + "' "
							str += "WHERE id =" + strconv.Itoa(r.id) + " ;"
							stm, err := db.Prepare(str)
							if err == nil {
								_, err2 := stm.Exec()
								if err2 == nil {
									ru = 1
								} else {
									ru = -1
								}
								stm.Close()
							} else {
								ru = -2
							}

							UpdateURezChan <- ru
						}()
					case u := <-IfExistAskChan:
						go func() {
							var id int = 0
							rows, err := db.Query("SELECT id FROM users WHERE email=? AND FA=1", u)
							if err == nil {
								for rows.Next() && id == 0 {
									err := rows.Scan(&id)
									if err != nil {
										fmt.Println(err)
									}
								}
								rows.Close()
							}

							IfExistRezChan <- id
						}()
					case x := <-LoginUAskChan:
						var u AccM
						u.Id = -1
						z := 0
						rows, err := db.Query("SELECT * FROM users WHERE Email='" + x[0] + "' AND Password='" + x[1] + "' ;")
						if err == nil {
							for rows.Next() && z == 0 {
								ert := rows.Scan(&u.Id, &u.NickName, &u.Name, &u.Family, &u.Email, &u.Password, &u.Telephone, &u.Ttl, &u.Kind, &u.Token, &u.Policy, &u.Fa, &u.DatFrom, &u.DatTo)
								if ert != nil {
									fmt.Println(ert.Error())
								}
								z = 1
							}
							rows.Close()
						} else {
							fmt.Printf(err.Error())
						}
						LoginURezChan <- u    */
					default:
						ddur := time.Since(TimeStamp)
						if int(ddur.Seconds()) >= 2 {
							err = db.Ping()
							if err == nil {
								TimeStamp = time.Now().UTC()
								//fmt.Println("+")
							} else {
								ex2 = 1
							}
						}
						time.Sleep(20 * time.Nanosecond)
					}
				}
				db.Close()
			} else {
				db.Close()
			}

		} else {
			fmt.Println("some bd error")
			fmt.Println(err.Error())
		}

	}
}

//=*Models*

//-*Models*
