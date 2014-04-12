package main

// Deals with the UI for the program. It probably will stay a console-based program though
//   because I personally feel that coders should be familiar with a terminal.


import (
	"strconv"
	"fmt"
)

func wr(s...interface{}) {
	for _,v := range s {
		fmt.Print(v)
	}
	fmt.Println()
}
type MenuStack []UIScreen
func (m MenuStack) Cur() UIScreen {
	return m[len(m)-1]
}
func (m*MenuStack) Push(u UIScreen) {
	*m = append(*m, u)
}
func (m*MenuStack) Pop() {
	*m = (*m)[:len(*m)-1]
}

type Action struct {
	Name string
	Act func(o Out) (UIScreen, error)
}
type UIScreen interface {
	Choices(o Out) []Action
}

var DontClear = dontClear{}
type dontClear struct{}
func (dontClear) Error() string {
	return ""
}
var PopParent = popParent{}
type popParent struct{}
func (popParent) Choices(o Out) []Action { return nil }


type NotImplemented string
func (n NotImplemented) Choices(o Out) []Action {
	o(string(n), " is not yet implemented, sorry")
	return nil
}
type DoubleCheck struct {
	msg string
	thing func(Out) error
}
func (d DoubleCheck) Choices(o Out) []Action {
	o("Are you sure you want to ", d.msg, "?")
	return []Action{
		{"Yes, of course", func(o Out) (UIScreen, error) {
			return PopParent, d.thing(o)}},
	}
}


func StartUI(input func()string) {
	var menus MenuStack
	var last string
	if err := Load(); err != nil {
		wr("There was an error with loading user data: ", err)
		wr("If this is your first time running the program, ignore this!")
		wr("Otherwise, existing data will be overwritten on next save")
	}
	
	if U.DoneTutorial {
		menus.Push(MainUI{})
	} else {
		menus.Push(TutUI{})
		U.DoneTutorial = true
	}
	
	var lastErr error
	
	for len(menus) > 0{
		if lastErr != nil {
			if lastErr != DontClear {
				wr(lastErr)
			}
			lastErr = nil
		} else {
			ClearScreen()
		}
		chs := menus.Cur().Choices(wr)
		if len(chs) == 0 {
			if chs == nil {
				fmt.Print("\n\nPress enter...")
				input()
			}
			menus.Pop()
			continue
		}
		wr()
		for i,v := range chs {
			wr(i, ": ", v.Name)
		}
		if len(menus) > 1 {
			wr(len(chs), ": Back")
		}
		wr()
		
		fmt.Print("Do what? ")
		line := input()
		if line == "" {
			if last == "" {
				continue
			}
			line = last
		} else {
			last = line
		}
		
		if ind, err := strconv.Atoi(line); err != nil {
			lastErr = err
		} else if ind < 0 || ind > len(chs) {
			lastErr = fmt.Errorf("index out of bounds: %v", ind)
		} else if ind == len(chs) {
			menus.Pop()
		} else {
			var next UIScreen
			next, lastErr = chs[ind].Act(wr)
			// Save lastErr until the next loop iteration
			if next == PopParent {
				menus.Pop()
			} else if next != nil {
				menus.Push(next)
			}
		}
	}
}
type TutUI struct{}
func (TutUI) Choices(o Out) []Action {
	o("Welcome to Kevin's Learning Thing!")
	o("This program helps you learn the Go programming language by throwing you straight into it and asking questions.")
	o("Many problems have hints you can unlock and some even have links to online information.")
	o("When in doubt, Google it! A search engine is can be a programmer's most helpful resource.\n")
	U.DoneTutorial = true
	Save()
	return []Action {
		{ "Press '0' and then 'Enter' or 'Return' to continue",
			func(Out) (UIScreen, error) { return MainUI{}, nil } },
	}
}

func Choice(u UIScreen) func(Out) (UIScreen, error) {
	return func(Out) (UIScreen, error) {
		return u, nil
	}
} 

type MainUI struct {}
func (MainUI) Choices(o Out) []Action {
	o("Welcome to Kevin's Learning Thing!")
	return []Action {
		{ "Problems", Choice(ProblemList{})},
		{ "Stats", Choice(Stats{})},
		{ "Settings", Choice(NotImplemented("Settings"))},
	}
}

// Todo: pagination
type ProblemList struct {}
func (ProblemList) Choices(o Out) []Action {
	o("Choose a problem!")
	ret := make([]Action, len(Probs))
	for i,v := range Probs {
		i := i
		var str = v.Name
		if U.IsSolved(i) {
			str += " (SOLVED)"
		}
		ret[i] = Action{str, Choice(ProblemMenu{i})}
	}
	return ret
}

type ProblemMenu struct {
	pid int
}
func (p ProblemMenu) Choices(o Out) []Action {
	pid := p.pid
	pr := Probs[pid]
	o(`Problem "`, pr.Name, `"`)
	ret := []Action{
		{"Open problem", func(o Out) (UIScreen, error){
			return nil, Edit(pid)}},
		{"Run all tests", func(o Out) (UIScreen, error){
			if err := Test(o, pid); err != nil {
				return nil, err
			}
			return ProblemSolved{pid}, nil}},
	}
	if len(pr.Help) > 0 {
		//ret = append(ret, "Help with this problem")
	}
	if len(pr.Hint) > 0 {
		//ret = append(ret, "Read hint")
	}
	
	// Only let them start over if they haven't solved it already
	//    (I don't want to deal with removing the solved status)
	if U.IsSolved(pid) {
		o("\nNote: You've solved this one already!\n")
	} else {
		ret = append(ret, Action{"Start problem over", Choice(DoubleCheck{
			"delete your work and start this problem over from scratch",
			func(Out) error {
				return WriteOut(pid)
			}})})
	}
	return ret
}
type ProblemSolved struct {
	pid int
}
func (p ProblemSolved) Choices(o Out) []Action {
	pid := p.pid
	pr := Probs[pid]
	o(`Problem "`, pr.Name, `" correctly solved!`)
	U.MarkSolved(pid)
	o("   Current points: ", U.Points)
	
	return nil
}
type Stats struct{}
func (Stats) Choices(o Out) []Action {
	o("User stats:")
	o("\tPoints: ", U.Points)
	
	o("\nProblem Name (Difficulty): SolvedStatus")
	for i,v := range Probs {
		var solv string
		if U.IsSolved(i) {
			solv = "Solved"
		} else {
			solv = "Unsolved"
		}
		o("\t", v.Name, " (", v.Difficulty, "): ", solv)
	}
	
	return nil
}
