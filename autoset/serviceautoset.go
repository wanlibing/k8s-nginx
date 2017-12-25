package autoset

import (
	"fmt"
	"demo/config"
	"reflect"
	//"strconv"

	"strings"
)


func ServiceExists(svcname string) (bool,error){
	fmt.Println("service is exist",svcname)
	return false,nil
}

func Createservice(svcname ,portnum string) (bool,error){
	//生成k8s service yaml 文件，
	autosetconfig := new(config.Config)
	autosetconfig.InitConfig("etc/k8s-monitor")
	fmt.Println("test for portnum is:",portnum,reflect.TypeOf(portnum))
	podlabel := strings.Split(svcname,"_")[0]
	sweaters := Inventory{"",svcname,portnum,podlabel}
	sweaters.AutoSet(autosetconfig.Read("default","SERVICETEMPFILE"),svcname+".yaml")
	fmt.Println("svcname tmp file parse",svcname)
//	fnamepath:= svcname +".yaml"

	//kubectl apply -f servicename.yaml
	fmt.Println("kubelet create file ")
	return true,nil
}