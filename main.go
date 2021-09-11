package main

import (
	"fmt"

	"github.com/Hassall/transit/transit"
)

func main() {
	fmt.Println(transit.TimedHTTPRequest("google.com"))
}
