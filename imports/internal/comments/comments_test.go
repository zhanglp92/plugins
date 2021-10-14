package comments

import (
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {

	var body = []byte(`
package main


import "github.com/zhanglp92/plugins/imports"
import "fmt"
import "gitlab.myteksi.net/gophers/go/food/food-search/common"
import "google.xxx/g"

func A() {}





func B() {}

func main() {
	fmt.Println("")
	imports.Process()
	common.A()	
	g.A()
}
`)

	res, err := Process(body)

	fmt.Println("err", err)
	fmt.Println(string(res))
}
