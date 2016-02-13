package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/tj/go-spin"
	"github.com/ttacon/chalk"
)

const (
	HideCursor = "\033[?25l"
	ShowCursor = "\033[?25h"
	Gray       = "\033[90m"
	ClearLine  = "\r\033[0K"
)

type Progress struct {
	Msg string
}

func (p *Progress) Done(msg string) {
	p.Update(msg + "\n")
}

func (p *Progress) Update(msg string) {
	fmt.Printf("\r\033[0K\033[90m%s %s", p.Msg, chalk.Cyan.Color(msg))
}

func NewProgress(msg string, fn func(func(str string)) error) error {

	p := &Progress{
		Msg: msg,
	}
	os.Stdout.Write([]byte(HideCursor))
	err := fn(p.Update)

	if err != nil {
		p.Done(chalk.Red.Color("error"))
	} else {
		p.Done(chalk.Green.Color("ok"))
	}
	os.Stdout.Write([]byte(ShowCursor))
	return err
}

type Process struct {
	Msg  string
	done chan bool
}

func (p *Process) Start() {
	os.Stdout.Write([]byte(HideCursor))
	defer os.Stdout.Write([]byte(ShowCursor))
	p.done = make(chan bool)

	ticker := time.NewTicker(100 * time.Millisecond)
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
	fmt.Printf("\r%s%s %s\r", Gray, p.Msg, chalk.Cyan.Color(msg))

}

func (p *Process) Done(msg string) {
	p.done <- true
	os.Stdout.Write([]byte(ShowCursor))
	fmt.Printf("\r%s%s %s\n", Gray, p.Msg, msg)

}

func NewProcess(msg string, fn func() error) error {

	p := &Process{Msg: msg}

	p.Start()

	err := fn()

	if err != nil {
		p.Done(chalk.Red.Color("error"))
	} else {
		p.Done(chalk.Green.Color("ok"))
	}
	return err
}
