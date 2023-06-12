package main

import (
	analyse "LogAnalyse/app/shared/kitex_gen/analyse/analyseservice"
	"log"
)

func main() {
	svr := analyse.NewServer(new(AnalyseServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
