
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
		defer f.Close()
		enc := gob.NewEncoder(f)
		if err = enc.Encode(U); err != nil {
			panic(err)
			return err
		}
		return nil
	}
}
// Loads all user data
func Load() {
	var ok bool
	if f, err := os.Open("user.dat"); err == nil {
		dec := gob.NewDecoder(f)
		if err = dec.Decode(&U); err != nil {
			panic(err)
		}
		f.Close()
		ok = true
	}
	if !ok {
		U.Probs = make(map[int]ProblemStatus)
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

type ProblemStatus uint
const (
	Solved ProblemStatus = 1 << iota
	HintUnlocked
)
type User struct {
	Probs map[int]ProblemStatus
	Points int
}

// Single user for now
var U User

func (u User) IsSolved(pid int) bool {
	return u.Probs[pid] & Solved != 0
}
func (u*User) MarkSolved(pid int) {
	if u.IsSolved(pid) { // double check that they don't just farm for points here
		return
	}
	u.Probs[pid] = u.Probs[pid] | Solved
	p := Probs[pid]
	u.Points += p.Difficulty*2 + 5
	
	if err := Save(); err != nil {
		fmt.Println("Error saving:", err)
	}
}

type Problem struct {
	Name string
	Difficulty int
	Help, Hint string
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








