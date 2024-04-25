package main

type Bar struct {
	id int
}

func (a *Bar) SetIDSeq(start int) { a.id = start }
