package main

import (
	"union/realapi/api"
	"net/http"

	log "github.com/alecthomas/log4go"
	"github.com/gocraft/web"
	"time"
)

func main() {

	rootRouter := web.New(api.RootContext{}).
		Middleware(web.StaticMiddleware("template", web.StaticOption{Prefix: "/template/"})).
		Middleware((*api.RootContext).Acao)
	//		Middleware(web.LoggerMiddleware).
	//		Middleware(web.ShowErrorsMiddleware).

	//***********
	//**apn******
	//***********
	apnRouter := rootRouter.Subrouter(api.Apn{}, "/apn")
	//add route
	//signup,login
	apnRouter.Post("/signup", (*api.Apn).SignUp)
	apnRouter.Post("/signup/code", (*api.Apn).GetVerificationCode)
	apnRouter.Post("/login", (*api.Apn).Login)
	apnRouter.Post("/password/reset", (*api.Apn).ResetPassword)

	apnRouter.Post("/media/getAll", (*api.Apn).GetAllMeida)
	apnRouter.Post("/media/add", (*api.Apn).AddOneMeida)


	//***********
	//**api******
	//***********
	apiRouter := rootRouter.Subrouter(api.Api{}, "/api")
	apiRouter.Middleware((*api.Api).ApiSessionMiddleware)
	apiRouter.Get("/test/get", (*api.Api).TestGet)
	apiRouter.Post("/test/post", (*api.Api).TestPost)



	ch := make(chan (error), 1)

	go func() {
		//http1 server
		ch <- http.ListenAndServe(":7777", rootRouter)
	}()

	err := <-ch
	log.Debug("realapi err:", err)

	<-time.After(10 * time.Second)
}
