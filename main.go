package main

import "github.com/yushengguo557/register/router"

func main() {
	r := router.NewRouter()
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
