
package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"gopkg.in/v1/yaml"
	"encoding/gob"
	"os"
)

const UserFile = "user.dat"
const ProblemFile = "probs.yml"


// Saves all user data
func Save() error {
	if f, err := os.Create("user.dat"); err != nil {
		return err
	} else {
		err = gob.NewEncoder(f).Encode(solved)
		f.Close()
		return err
	}
}
// Loads all user data
func Load() {
	var ok bool
	if f, err := os.Open("user.dat"); err == nil {
		if err = gob.NewDecoder(f).Decode(&solved); err == nil {
			ok = true
		}
		f.Close()
	}
	if !ok {
		solved = make(map[int]bool)
	}
}


func init() {
	Load() // Ignore if it fails
}


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

var solved map[int]bool
func IsSolved(pid int) bool {
	return solved[pid]
}
func MarkSolved(pid int) {
	solved[pid] = true
	
	if err := Save(); err != nil {
		fmt.Println("Error saving:", err)
	}
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


func sl(s...string) []string {return s}

var Probs = LoadProblems()
func LoadProblems() []Problem {
	f, err := os.Open(ProblemFile)
	if err == nil {
		var ret []Problem
		var all []byte
		all, err = ioutil.ReadAll(f)
		if err = yaml.Unmarshal(all, &ret); err == nil {
			// Success!
			return ret
		}
	}
	panic(err)
}








