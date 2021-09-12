package main

import (
	"fmt"

	"github.com/Hassall/transit/pkg/http"
)

func main() {
	fmt.Println(http.TimedHTTPRequest("google.com"))
}
