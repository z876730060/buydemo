package main

import (
	"fmt"

	"github.com/z876730060/buydemo/config"
	"github.com/z876730060/buydemo/database"
	"github.com/z876730060/buydemo/middlewares"
	"github.com/z876730060/buydemo/router"
)

func main() {
	cfg := config.Load()

	// Init JWT
	middlewares.InitJWT(cfg.JWTSecret)

	// Init database
	database.Init(cfg)

	// Setup router
	r := router.Setup()

	fmt.Printf("ERP系统已启动，监听端口 %s\n", cfg.ServerPort)
	fmt.Printf("访问地址: http://localhost%s\n", cfg.ServerPort)
	fmt.Printf("默认账号: %s / %s\n", cfg.AdminUser, cfg.AdminPass)

	if err := r.Run(cfg.ServerPort); err != nil {
		panic(err)
	}
}
