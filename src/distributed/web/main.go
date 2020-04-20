package main

import (
	"net/http"

	"github.com/sophiedebenedetto/power_plant/src/distributed/web/controller"
)

func main() {
	controller.Initialize()
	http.ListenAndServe(":3000", nil)
}
