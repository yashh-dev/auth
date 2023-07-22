package main

import (
	"miauw.social/auth/handlers"
)

func main() {
	var forever chan struct{}
	go Serve("auth.initial", handlers.UserCreate)
	go Serve("auth.login", handlers.UserLogin)
	go Serve("auth.verify", handlers.UserVerify)
	go Serve("auth.sessions.get", handlers.GetUserSession)
	go Serve("auth.sessions.exists", handlers.ExistsUserSession)
	<-forever
}
