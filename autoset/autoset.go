package autoset

import (
	"log"
	"text/template"
	"os"
	"fmt"
	"strings"
)


/*
func InitConfig(path string){
	fmt.Println("init config") //操蛋，第三方包没有遍历方法
}
*/

type Inventory struct {
	Projectname string
	Servicename string
	Targetport string
	Podlabels string
}



func GetServiceName(key string) string {
	fmt.Println(key)
	configPath := strings.Split(key,"/")
	configNameSpace:=configPath[len(configPath)-2]
	configName := configPath[len(configPath)-1]
	serviceName :=  configName+"."+configNameSpace
	return serviceName

}

//判断文件是否存在，如不存在则将渲染后的内容输出到配置文件，并且reload nginx
func PathExists(path string)(bool, error){
	fmt.Println("input mopan into file")
	confpath := "confd/" + path+".conf"
	fmt.Println(confpath)
	_,err := os.Stat(confpath)
	if err != nil && os.IsNotExist(err){
		fmt.Println("file not exists")
		fmt.Printf("%q file not exits,begin create nginx confd file\n",confpath)
		return false,err

	}
	fmt.Println("file  exits")
	fmt.Printf("%q file is exits,continue minitor \n",confpath)
	return true,nil

}

func check(e error){
	if e != nil{
		fmt.Printf(e.Error())
		panic(e)
	}
}

func (sweaters Inventory)AutoSet(tempFile,fname string){
	//判断文件是否存在
	ifexistsfile,err := PathExists(sweaters.Servicename)  //true,nil key="reportjt-houtai-service.default"
	if err != nil || ! ifexistsfile {
		//fname := condir + key +".conf"
		fmt.Println(fname)
		servicefile,errf := os.Create(fname) //如果文件存在会被覆盖

		defer servicefile.Close()
		if errf != nil{
			fmt.Println("create file failed")
			fmt.Println(errf)
			fmt.Printf("create file %q failed ,err is %q \n",fname,errf.Error())

		}

		tmpl, err := template.ParseFiles(tempFile)
		if err != nil {
			fmt.Printf("log temlpate failed,err is %q",err.Error())
			log.Println("load temlpate file failed...")
			panic(err)
		}
		//模板中变量首字母必须为大写
		//sweaters = Inventory{strings.Split(sweaters.Servicename,"-")[0],sweaters.Servicename,""}
		log.Println(sweaters)
		err = tmpl.Execute(servicefile, sweaters)

	} else {
		//配置文件存在
		fmt.Println("file exits,exit ...")
	}
}

func DeleteConfig(key,fname string)  {
	ifexistsfile,err := PathExists(key)  //true,nil
	if err == nil &&  ifexistsfile {
		//fname := condir + key +".conf"
		err := os.Remove(fname)
		check(err)
		log.Println("file has been remove")
		fmt.Printf("nginx confd file has been delete \n",fname)
		}
}


