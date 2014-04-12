
package main

import (
	"fmt"
	"bufio"
	"os"
)

func getInput() func()string {
	sc := bufio.NewScanner(os.Stdin)
	return func() string {
		// TODO: errors
		sc.Scan()
		return sc.Text()
	}
}
func main() {
	fmt.Println("Started!")
	
	
	
	StartUI(getInput())
}










