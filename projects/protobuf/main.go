package main

import (
	"fmt"

	simplepb "github.com/Wallruzz9114/protobuf/simple"
)

func main() {
	doSimple()
}

func doSimple() {
	sm := simplepb.Simple{
		Id:         12345,
		IsSimple:   true,
		Name:       "My Simple Message",
		SampleList: []int32{1, 4, 7, 8},
	}

	fmt.Println(sm)
}
