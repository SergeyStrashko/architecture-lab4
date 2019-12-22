package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/SergeyStrashko/architecture-lab4/engine"
)

type printCommand struct {
	arg string
}

func (p *printCommand) Execute(loop engine.Handler) {
	fmt.Println(p.arg)
}

type splitCommand struct {
	str, separator string
}

func (spl *splitCommand) Execute(loop engine.Handler) {
	strArr := strings.Split(string(spl.str), string(spl.separator))

	for _, s := range strArr {
		loop.Post(&printCommand{arg: s})
	}
}

func parse(commandline string) engine.Command {
	parts := strings.Fields(commandline)
	if len(parts) == 0 {
		return &printCommand{arg: "SYNTAX ERROR: Command not found"}
	}
	switch parts[0] {
	case "print":
		if len(parts) == 1 {
			return &printCommand{arg: "SYNTAX ERROR: Arg not found"}
		}
		return &printCommand{arg: parts[1]}
	case "split":
		if len(parts) == 1 {
			return &printCommand{arg: "SYNTAX ERROR: Args not found"}
		}
		if len(parts) == 2 {
			return &printCommand{arg: "SYNTAX ERROR: Separator not found"}
		}
		return &splitCommand{str: parts[1], separator: parts[2]}
	default:
		return &printCommand{arg: "SYNTAX ERROR: Unexpected command"}
	}
}

func main() {
	eventLoop := new(engine.EventLoop)
	eventLoop.Start()

	if input, err := os.Open("./commands.txt"); err == nil {
		defer input.Close()

		scanner := bufio.NewScanner(input)
		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := parse(commandLine)
			eventLoop.Post(cmd)
		}
	}

	eventLoop.AwaitFinish()
}
