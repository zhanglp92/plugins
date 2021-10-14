package internal

import (
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {
	var body = []byte(`
// jksfhdjksahfkads
/*
fdjsfkadls
fdsajmfkdjsk


fdsajkf;dasj


*/
package main


import "github.com/zhanglp92/plugins/imports"
import "fmt"
import "gitlab.myteksi.net/gophers/go/food/food-search/common"
import "google.xxx/g"


type Aaaa struct {
  Aa int // Aa ...
  Bb string // xxxx
}

func (m*Aaaa) Check() {}

func A() {}




// Bbb
func B() {}

func main() {
	Ax := func() {}


	fmt.Println("")
	imports.Process()
	common.A()	
	g.A()
}
`)

	res, err := Process(body, true)

	fmt.Println("err", err)
	fmt.Println(string(res))
}
