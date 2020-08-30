package main

import (
	"fmt"
	"net/http"


	"github.com/seaung/gin-blog/pkg/setting"
	"github.com/seaung/gin-blog/routers"
)

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr: fmt.Sprintf(":%d", setting.HttpPort),
		Handler: router,
		ReadTimeout: setting.ReadTime,
		WriteTimeout: setting.WriteTime,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()
}
