//package etcdclient

package etcdclient

import (
	"time"
	"fmt"
	"reflect"
	"github.com/coreos/etcd/client"
	"demo/config"
	"log"
	"demo/autoset"
	"context"
	"demo/k8smonitorlog"
	"strings"
)

func EtcdConn(autosetconfig *config.Config) client.KeysAPI  {
	cfg := client.Config{
		Endpoints:               []string{autosetconfig.Read("etcd","ETCD_EP")},
		Transport:               client.DefaultTransport,
		// set timeout per request to fail fast when the target endpoint is unavailable
		HeaderTimeoutPerRequest: 5*time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		fmt.Printf("conntect etcdsever failed %q \n",err.Error())
	}
	kapi := client.NewKeysAPI(c)
	fmt.Println("kapi type is",reflect.TypeOf(kapi))
	return kapi
}



func StartAutoSet() {

	autosetconfig := new(config.Config)
	autosetconfig.InitConfig("etc/k8s-monitor")

	//初始化日志实例
	filename := autosetconfig.Read("default","LOG_DIR")
	logfile := k8smonitorlog.Openfile(filename)
	defer logfile.Close()

	//infoLog := log.New(logfile,"[INFO]",log.LstdFlags+log.Llongfile)
	//infoLogOutput := k8smonitorlog.LogWrite{infoLog}

	//watch key by wanlb at 2017-12-07
	kapi := EtcdConn(autosetconfig)
	watcher := kapi.Watcher(autosetconfig.Read("etcd","ETCD_PREFIX"), &client.WatcherOptions{
		Recursive:true,
	})
	for {
		res,err := watcher.Next(context.Background())
		if err != nil{
			log.Println("err watch workers",err)
			break
		}

		//程序启动时候初始化配置目录
		fmt.Println(res.Action)    //打印k8s对etcd的动作
		if res.Action == "set" || res.Action == "update" || res.Action=="create" || res.Action=="compareAndSwap"{
			fmt.Println("key is change")
			fmt.Printf("%q has change or create,begin to monitor nginx conffig file \n",res.Node.Key)

			key := res.Node.Key
			configName :=  autoset.GetServiceName(key)
			fnamepath := autosetconfig.Read("default","NGINXCONFD") + configName + ".conf"
			sweaters := autoset.Inventory{strings.Split(configName,"-")[0],configName,"",""}
			sweaters.AutoSet(autosetconfig.Read("default","NGINXTEMPFILE"),fnamepath)
			//autoset.AutoSet(configName,autosetconfig.Read("default","NGINXTEMPFILE"),fnamepath)
		}
		if res.Action == "delete" || res.Action == "compareAndDelete"{
			fmt.Printf("%q has been deleted,begin to delete nginx config file\n",res.Node.Key)
			key := res.Node.Key
			configName :=  autoset.GetServiceName(key)
			fnamepath := autosetconfig.Read("default","NGINXCONFD") + configName + ".conf"
			//fmt.Println(fnamepath)
			autoset.DeleteConfig(configName,fnamepath)

		}
	}
}

/*
func main()  {
	StartAutoSet()
}
*/