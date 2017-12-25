package main

import (
	"fmt"
//	"html"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"time"
//	"encoding/json"
	"demo/autoset"
	"demo/etcdclient"

)

type Todo struct{
	Name string   `json:"jsname"`
	Completed bool  `json:"completed"`
	Due time.Time 	`json:"due"`	  //相当于注解
}


type Todos []Todo

//var testgo string

//monitor service and create service use k8s client
func serviceMonitor(w http.ResponseWriter,r *http.Request){

	//获取GET请求参数
	r.ParseForm()
	servicename:=""
	if len(r.Form["servicename"])>0{
		servicename=r.Form["servicename"][0]
	}
	portnum := r.Form["port"][0]
	fmt.Printf("myname is:%q ,%q\n",servicename,portnum)
	//monitor k8s service if exists
	resp := ""
	flag,err := autoset.ServiceExists(servicename)   //查看etcd是否有KEY
	if err != nil{
		fmt.Println(err)
	}
	//if k8s service not exists , parse service tmpl file ,create k8s service

	if !flag{
		f,_ :=autoset.Createservice(servicename,portnum)    //1、创建service文件2、调用k8s接口
		if f{
			resp = servicename + "create success"
		}
		resp = servicename + " create failed"
	} else {
		resp = servicename + " service is exist.."
	}
	fmt.Fprintln(w,resp)
}


func MonitorStart() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/serviceMonitor.json",serviceMonitor)
	log.Fatal(http.ListenAndServe(":8080",router))

}


func main(){

	go MonitorStart()
	go etcdclient.StartAutoSet()
	for{
		time.Sleep(time.Second*100)
	}

}