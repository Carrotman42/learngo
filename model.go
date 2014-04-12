
package main

import (
	"fmt"
	"io"
	"gopkg.in/v1/yaml"
)
var _ = yaml.Marshal


const WorkDir = "workspace"

type TestCase struct {
	Input []string
	Output []string
}

func (t TestCase) Write(o Out) {
	invals := SepList{Sep: ","}
	for _,v := range t.Input {
		invals.Append(v)
	}
	
	outvars := SepList{Sep: ","}
	outvals := SepList{Sep: ","}
	test := SepList{Sep: "||"}
	for i,v := range t.Output {
		n := fmt.Sprint("o", i)
		outvars.Append(n)
		test.Append("(" + n + ")!=(" + v + ")")
		outvals.Append(v)
	}
	
	o("{")
	o(outvars, ":=solve(", invals, ")")
	// Note: We'll always have at least one outvar, or else it's not really a test case
	o("if ", test, ` {
fmt.Println("Error in test case:")
fmt.Println("   inputs   : `, invals, `")
fmt.Println("   expected : ",`, outvals, `)
fmt.Println("   actual   : ", `, outvars, `)
}`)
	
	o("}")
}

var solved = make(map[int]bool)
func IsSolved(pid int) bool {
	return solved[pid]
}
func MarkSolved(pid int) {
	solved[pid] = true
}

type Problem struct {
	Name string
	Difficulty int
	Parts []string
	Tests []TestCase
}
func WriteDefault(pid int, dest io.Writer) {
	o := func(v...interface{}) {
		fmt.Fprintln(dest, v...)
	}
	p := Probs[pid]
	for _,v := range p.Parts {
		o(v)
	}
	
	o("func main() {")
	for _,v := range p.Tests {
		v.Write(o)
	}
	o("}")
}
func GetFile(pid int) string {
	return fmt.Sprint(WorkDir, "/", pid, ".go")
}


var raw = `name: Hello World
difficulty: 0
parts:
- |
   package main
   import "fmt"
   
   // Should return a string that acknowledges the world
   func solve() string {
   	return ""
   }
tests:
- input:
  output:
  - '"Hello World"'
`

func sl(s...string) []string {return s}

var Probs = LoadProblems()
func LoadProblems() []Problem {
	type Hey struct {
		A string
	}
	h := Hey{`"Hello World"`}
	a,b:=yaml.Marshal(h)
	fmt.Println(string(a), b)
	var ret Problem
	if err := yaml.Unmarshal(([]byte)(raw), &ret); err != nil {
		panic(err)
	}
	return []Problem{ret}
}








