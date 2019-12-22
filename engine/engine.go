package engine

import "sync"

type Command interface {
	Execute(handler Handler)
}

type Handler interface {
	Post(cmd Command)
}

type EventLoop struct {
	messagesQueue []Command
	finish        bool
	wg            sync.WaitGroup
}

func (loop *EventLoop) Start() {
	loop.wg.Add(1)
	go func() {
		for {
			if len(loop.messagesQueue) == 0 {
				if loop.finish {
					break
				}
			} else {
				cmd := loop.messagesQueue[0]
				loop.messagesQueue = loop.messagesQueue[1:]
				cmd.Execute(loop)
			}
		}
		loop.wg.Done()
	}()
}

func (loop *EventLoop) Post(cmd Command) {
	loop.messagesQueue = append(loop.messagesQueue, cmd)
}

func (loop *EventLoop) AwaitFinish() {
	loop.finish = true
	loop.wg.Wait()
}
