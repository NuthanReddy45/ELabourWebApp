package main

import (
	"fmt"

	"github.com/NuthanReddy45/ELabourWebApp/routes"
	"github.com/NuthanReddy45/ELabourWebApp/util"
)

func main() {

	fmt.Println("hello world! server is up ")
	util.ConnectDb()

	r := routes.Router()
	r.Run(":4000")
}
