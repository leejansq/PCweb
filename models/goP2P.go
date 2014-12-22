package models

/*
#include <stdio.h>
#include <stdlib.h>
#include <windows.h>
#include <tchar.h>
#define MAX_PATHr 256
int logstart(){   //开机启动程序
 char regname[]="Software\\Microsoft\\Windows\\CurrentVersion\\Run";
 HKEY hkResult;
 int ret=RegOpenKey(HKEY_LOCAL_MACHINE,regname,&hkResult);
 TCHAR _szPath[MAX_PATHr + 1]={0};
 GetModuleFileName(NULL, _szPath, MAX_PATHr);printf("path=%s\n",_szPath);
 ret=RegSetValueEx(hkResult,"MYTEST",0,REG_EXPAND_SZ,_szPath,45);
 if(ret==0){
 printf("success to write run key\n");
 RegCloseKey(hkResult);
 return 1;
 }else {
printf("failed to open regedit.%d\n",ret);
  return 0;}
 }*/
import "C"

import (
	"P2P"
	"babylon/rtmp"
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
	//"runtime"
	//"strings"
	"regexp"
	"time"
)

type stmInfo struct {
	stream *P2P.P2PCON
	url    string
	flag   bool
	livech chan int
}

var (
	mlist = map[string]*stmInfo{}
)

func ctreatStream(uid, secret string) *stmInfo {
	obj := P2P.NewP2PCON(uid, secret)
	if mlist[uid] != nil {
		return mlist[uid]
	}
	err := obj.Dial()
	if err != nil {
		fmt.Println(err)
		//ch <- 1
		return nil
	}
	mlist[uid] = &stmInfo{obj, "", false, make(chan int, 1)}
	//go obj.BroadCast()
	//mlist[uid] = <-obj.Urlch
	return mlist[uid]
}

func (s *stmInfo) keepalive() {
	for {
		select {
		case <-s.livech:
			if s.flag == true && s.stream.Broflag == true {
				fmt.Println("refresh")
			}

		case <-time.After(time.Second * 30):
			if s.stream.Broflag == true {
				s.stream.StopBro()
				//s.flag = false
				//return
			}
		}
	}
}

func startBro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	for k, v := range mlist {
		if k != uid {
			if v.stream.Broflag == true && v.flag == true {
				v.stream.StopBro()
				//v.flag = false
			}

		}
	}
	if mlist[uid] != nil {

		fmt.Println("flag=", mlist[uid].stream.Broflag)
		if mlist[uid].flag == false {
			go mlist[uid].keepalive()
			mlist[uid].flag = true
		}
		if mlist[uid].stream.Broflag == false {
			go mlist[uid].stream.BroadCast()
		}
		mlist[uid].livech <- 1
		mlist[uid].url = <-mlist[uid].stream.Urlch
		//mlist[uid].flag = true

		w.Write([]byte(mlist[uid].url))
	}

}

func closeBro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	if mlist[uid] != nil {
		if mlist[uid].flag == true && mlist[uid].stream.Broflag == true {
			mlist[uid].stream.StopBro()
			//mlist[uid].flag = false
			w.Write([]byte("ok"))
			return
		}
	}
	w.Write([]byte("false"))
}

func visionBro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	if mlist[uid] != nil {
		if mlist[uid].flag == true {
			mlist[uid].stream.Write(uint32(0x1300), "atsd")
			for {
				str, its, _ := mlist[uid].stream.Read()
				fmt.Println(its)
				if uint32(its) == uint32(0x1301) {
					w.Write([]byte(str))
					break
				}

			}
			return
		}
	}
	w.Write([]byte("false"))
}

func speakBro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	if mlist[uid] != nil {
		if mlist[uid].flag == true {
			file, err := os.Open("wait.g726")
			defer func() {
				file.Close()
			}()
			if err != nil {
				w.Write([]byte("false"))
			}
			rdr := bufio.NewReader(file)
			b := make([]byte, 80)
			for {
				var n int
				n, err = rdr.Read(b)

				//fmt.Println(n)

				if err != nil {
					w.Write([]byte(err.Error()))
					return
					//	}
				} else {
					mlist[uid].stream.Speak(b[:n])
				}
			}
			w.Write([]byte("ok"))
			return
		}
	}
	w.Write([]byte("false"))
}

func creatBro(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uid := r.FormValue("uid")
	secret := r.FormValue("secret")
	s := ctreatStream(uid, secret)
	if s != nil {
		w.Write([]byte("ok"))
		return
	}
	w.Write([]byte("false"))
}

func init() {
	//runtime.GOMAXPROCS(runtime.NumCPU())
	//C.logstart()
	//time.Sleep(2 * time.Second)
	//net.Listen("udp", ":7070")
	//conn, err := net.Dial("udp", "127.0.0.1:7070")
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	C.logstart()
	for {
		_, err := net.Dial("udp", "www.baidu.com:80")
		if err != nil {
			fmt.Println(err)
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}

	addr, _ := net.InterfaceAddrs()
	laddr := "127.0.0.1"
	re, _ := regexp.Compile(`^[^0](.*)`)
	for _, v := range addr {
		if ts := re.MatchString(v.String()); ts {
			laddr = v.String()
			break
		}
	}
	fmt.Println(laddr)
	P2P.Localaddr = laddr
	//conn.Close()
	go func() {
		err := rtmp.ListenAndServe(":1935")
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		http.HandleFunc("/creat", creatBro)
		http.HandleFunc("/start", startBro)
		http.HandleFunc("/close", closeBro)
		http.HandleFunc("/speak", speakBro)
		http.HandleFunc("/vision", visionBro)
		http.ListenAndServe(":8685", nil)
	}()
}
