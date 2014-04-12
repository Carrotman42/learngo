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
		if len(menuStack) > 1 {
			wr(len(ch), ": Back")
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
		if i,err := strconv.Atoi(choice); err != nil || i < 0 || i >= len(ch) {
			if i == len(ch) {
				menuStack = menuStack[:lll]
				clear = true
				continue
			}
			if err == nil {
				err = fmt.Errorf("out of range: %v", i)
			}
			wr("\tOops: invalid choice: ", err)
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

type DoubleCheck struct {
	msg string
	action func()error
	done bool
}
func (d DoubleCheck) Choices(o Out) []string {
	if d.done { //already called this action
		return nil
	}
	o("Are you sure you want to ", d.msg, "?")
	return []string{"Yes, of course"}
}
func (d*DoubleCheck) Choose(o Out, i int) (UIScreen, error) {
	// We'll only get in here if they choose to confirm
	d.done = true
	return nil, d.action()
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
			if U.IsSolved(i) {
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
		} else {
			ret = []string{
				"Open problem",
				"Run all tests",
			}
			if len(pr.Help) > 0 {
				//ret = append(ret, "Help with this problem")
			}
			if len(pr.Hint) > 0 {
				//ret = append(ret, "Read hint")
			}
			
			if U.IsSolved(pid) {
				o("\nNote: You've solved this one already!\n")
			} else {
				// Only let them start over if they haven't solved it already
				//    (I don't want to deal with removing the solved status)
				ret = append(ret, "Start problem over")
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
	}
	
	var err error
	switch i {
		case 0: // Pop up editor
			err = Edit(pid)
		case 1: // Run tests
			err = Test(o, pid) //Blocks until tests are done
			if err == nil {
				// Mark the problem as solved
				U.MarkSolved(pid)
				p.JustSolved = true
			}
		case 2:
			return &DoubleCheck{"delete your work and start this problem over from scratch",
				func() error {
					return WriteOut(pid)
				}, false,
			}, nil
	}
	return nil, err
}





