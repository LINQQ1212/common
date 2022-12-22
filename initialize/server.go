package initialize

import (
	"fmt"
	admin "github.com/LINQQ1212/common/admin/router"
	"github.com/LINQQ1212/common/global"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type server interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

func RunWebServer(fc func()) {
	Router := Routers()
	Router.StaticFS(global.CONFIG.Local.Path, http.Dir(global.CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址

	//frontend
	//frontend.InitRouter()
	fc()
	address := fmt.Sprintf(":%d", global.CONFIG.System.Addr)
	srv := initServer(address, http.DefaultServeMux)

	// admin
	AdmRouter := Routers()
	admin.InitRouter(AdmRouter)
	adminAddr := fmt.Sprintf(":%d", global.CONFIG.System.AdminAddr)
	admSrv := initServer(adminAddr, AdmRouter)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	go func() {
		if err := admSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()

	fmt.Printf(`
	默认地址:http://127.0.0.1%s
	admin地址:http://127.0.0.1%s
`, address, adminAddr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Server exit")
	global.Close()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	if err := admSrv.Shutdown(ctx); err != nil {
		log.Fatal("admin Server Shutdown:", err)
	}

	log.Println("Server exiting")
}
