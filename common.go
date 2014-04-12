
package main

type SepList struct {
	Sep string
	cur string
}
func (l*SepList) Append(s string) {
	if l.cur == "" {
		l.cur = s
	} else {
		l.cur += l.Sep + s
	}
}
func (l SepList) String() string {
	return l.cur
}


type Out func(...interface{})