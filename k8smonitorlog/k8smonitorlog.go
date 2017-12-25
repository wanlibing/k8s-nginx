package k8smonitorlog

import (
	"log"
	"os"
)

type LogWrite struct {
	*log.Logger
}

func Openfile(filename string) *os.File{
	logFile,err  := os.OpenFile(filename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
//	defer logFile.Close()   //when execute this
	if err != nil {
		log.Fatalln("open file error !")
	}
	return logFile
}


