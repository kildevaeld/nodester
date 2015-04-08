package main

import (
	"fmt"
	"github.com/tj/go-spin"
	"time"
)

type Process struct {
	Msg  string
	done chan bool
}

func (p *Process) Start() {

	p.done = make(chan bool)

	ticker := time.NewTicker(1000 * time.Microsecond)
	s := spin.New()

	go func() {
	loop:
		for {

			select {
			case <-p.done:
				ticker.Stop()
				break loop
			case <-ticker.C:
				p.update(s.Next())
			}

		}

		close(p.done)

	}()

}

func (p *Process) update(msg string) {
	fmt.Printf("\r%s %s\r", p.Msg, msg)

}

func (p *Process) Done(msg string) {
	p.done <- true
	fmt.Printf("\r%s %s\n", p.Msg, msg)

}

func NewProcess(msg string, fn func() error) error {

	p := &Process{Msg: msg}

	p.Start()

	err := fn()

	if err != nil {
		p.Done("error")
	} else {
		p.Done("ok")
	}
	return err
}

type Progress struct {
	*Process
}

func (p *Progress) Start() {
	/*p.stop = make(chan bool)
	p.Progress = make(chan string)
	go func() {
	loop:
		for {
			select {
			case <-p.stop:
				break loop
			case str := <-p.Progress:
				p.Process.update(str)
			}
		}

		close(p.Progress)
		close(p.stop)
	}()*/

}

func (p *Progress) Done(msg string) {
	fmt.Printf("\033[2K\r%s %s\r\n", p.Msg, msg)
}

func (p *Progress) Progress(msg string) {
	p.Process.update(msg)
}

func NewProgress(msg string, fn func(p *Progress) error) error {

	p := &Progress{
		Process: &Process{
			Msg: msg,
		},
	}

	p.Start()

	err := fn(p)

	if err != nil {
		p.Done("error")
	} else {
		p.Done("ok")
	}
	return err
}
