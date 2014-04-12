package main

import (
	"os/exec"
	"os"
	"fmt"
)

func ClearScreen() {
    cmd := exec.Command("cmd", "/c", "cls")
    cmd.Stdout = os.Stdout
    cmd.Run()
}

type CompileError struct {
	er error
	out []byte
}
func (c CompileError) Error() string {
	return fmt.Sprint("Compiler error: ", c.er, "\n\n", string(c.out))
}
type RunError struct {
	er error
	out []byte
}
func (c RunError) Error() string {
	if c.er == nil {
		return string(c.out)
	}
	return fmt.Sprint("Run error: ", c.er, "\n\n", string(c.out))
}

func Test(o Out, pid int) error {
	// First compile it
	output, err := exec.Command("go", "build", "-o", "test.exe", GetFile(pid)).CombinedOutput()
	if err != nil {
		return CompileError{err, output}
	}
	
	// Run it!
	output, err = exec.Command("test.exe").CombinedOutput()
	if err != nil || len(output) > 0 {
		return RunError{err, output}
	}
	return nil
}

func WriteOut(pid int) error {
	fi := GetFile(pid)
	if f, err := os.Create(fi); err != nil {
		return err
	} else {
		WriteDefault(pid, f)
		f.Close()
		if err = exec.Command("go", "fmt", fi).Run(); err != nil {
			return err
		}
	}
	return nil
}

func Edit(pid int) error {
	// Check to see if it exists
	fi := GetFile(pid)
	if _, err := os.Stat(fi); err != nil {
		fmt.Println("STAT err: ", err)
		
		if err := WriteOut(pid); err != nil {
			return err
		}
	}
	return exec.Command("notepad++.exe", fi).Run()
}
















