package config

import (
	"bufio"
	"fmt"
//	"io"
	"os"
	"io"
	"strings"
)

const middle  = "=="

type Config struct {
	Mymap map[string]string
	item string
}

func (c *Config) InitConfig(path string)  {     //不用指针参数无法初始化实例
	c.Mymap =make(map[string]string)
	f,err := os.Open(path)
	if err != nil{
		fmt.Println("open file failed")
		panic(err)   //抛出异常，并终止程序运行
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for{
		b,_,err := r.ReadLine()   //返回的是字符的位数
		if err != nil{
			if err == io.EOF{
				break
			}
			panic(err)
		}
		s := strings.TrimSpace(string(b))    //字符串显示
		if strings.Index(s,"#") == 0{
		//	fmt.Println(s)
			continue
		}
		n1 := strings.Index(s,"[")
		n2 := strings.Index(s,"]")
		if n1 > -1 && n2 > -1 && n2 > n1 + 1{
			c.item = strings.TrimSpace(s[n1+1:n2])
			fmt.Println(s)
			continue
		}
		if len(c.item) == 0{
			continue
		}
		index := strings.Index(s,"=")
		if index ==0 || index == len(s) -1 {
			//fmt.Println(s)
			continue
		}
		key := strings.TrimSpace(s[:index])
		key = c.item + middle + key
		value := strings.TrimSpace(s[index+1:])
		//fmt.Println(key)
		//fmt.Println(value)
		c.Mymap[key]=strings.TrimSpace(value)
	}

}

//读取配置文件
func (c Config) Read(item ,key string) string{
	key = item + middle + key
	v,found := c.Mymap[key]
	if !found{
		return ""
	}
	return v
}

