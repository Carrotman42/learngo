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
func StartUI(input func()string) {
	menuStack := []UIScreen{MainMenu{}}
	
	var last string
	clear := true
	for {
		if clear {
			ClearScreen()
			clear = false
		}
		lll := len(menuStack)-1
		cur := menuStack[lll]
		ch := cur.Choices(wr)
		
		if len(ch) == 0 {
			// This means the ui screen wants to be popped off the stack
			fmt.Print("\nPress enter...")
			input()
			
			// After waiting for them, do that
			menuStack[lll] = nil // for gc
			menuStack = menuStack[:lll]
			
			clear = true
			continue
		}
		
		for i,v := range ch {
			wr(i, ": ", v)
		}
		
		fmt.Print("\nWhat do? ")
		choice := input()
		if choice == "" {
			if last == "" {
				continue
			}
			choice = last
		} else {
			last = choice
		}
		
		fmt.Println()
		if i,err := strconv.Atoi(choice); err != nil || i < 0 || i > len(ch) {
			wr("\tOops: invalid choice:", err)
		} else if n, err := cur.Choose(wr, i); err != nil {
			wr(err)
		} else {
			if n != nil { // if n==nil it means stay on page
				// Else navigate further into the menus
				menuStack = append(menuStack, n)
			}
			// Successful, clear the screen!
			clear = true
		}
		fmt.Println()
	}
}

type NotImplemented string
func (n NotImplemented) Choices(o Out) []string {
	o(string(n), " is not yet implemented, sorry")
	return nil
}
func (n NotImplemented) Choose(o Out, i int) (UIScreen, error) {
	return nil, nil
}

type Out func(...interface{})

type UIScreen interface {
	Choices(o Out) []string
	// May block as long as you want
	Choose(o Out, i int) (UIScreen, error)
}

type MainMenu struct {}

func (MainMenu) Choices(o Out) []string {
	o("Welcome to Kevin's Learning Thing")
	return []string {
		"Problems",
		"Stats",
		"Settings",
	}
}
func (MainMenu) Choose(o Out, i int) (UIScreen, error) {
	switch i {
		case 0: return &ProblemScreen{-1,false}, nil
		case 1: return NotImplemented("Stats"), nil
		case 2: return NotImplemented("Settings"), nil
	}
	return nil, nil
}

type ProblemScreen struct {
	Prob int
	JustSolved bool
}
func (p ProblemScreen) Choices(o Out) []string {
	// Write out all of the possible problems
	pid := p.Prob
	if pid == -1 {
		o("Here's a list of problems to look at:")
		ret := make([]string, len(Probs))
		for i,v := range Probs {
			str := v.Name
			if IsSolved(i) {
				str += " (SOLVED)"
			}
			ret[i] = str
		}
		return ret
	} else {
		pr := Probs[pid]
		o(`Problem "`, pr.Name, `"`)
		
		var ret []string
		if p.JustSolved {
			o("\tYou've completed the problem, good job!")
		} else if IsSolved(pid) {
			o(`You've solved this one already!`)
		} else {
			ret = []string{
				"Open problem",
				"Run all tests",
				"Start problem over",
			}
		}
		return ret
	}
}
func (p*ProblemScreen) Choose(o Out, i int) (UIScreen, error) {
	pid := p.Prob
	if pid == -1 {
		// Problem submenu
		return &ProblemScreen{i,false}, nil
	} else if IsSolved(pid) {
		
	}
	
	var err error
	switch i {
		case 0: // Pop up editor
			err = Edit(pid)
		case 1: // Run tests
			err = Test(o, pid) //Blocks until tests are done
			if err == nil {
				// Mark the problem as solved
				MarkSolved(pid)
				p.JustSolved = true
			}
		case 2:
			err = WriteOut(pid)
	}
	return nil, err
}





