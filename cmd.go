package velvet

import "time"

// Cmd is a command that returns a message
type Cmd func() Msg

// Batch executes multiple commands concurrently
func Batch(cmds ...Cmd) Cmd {
	return func() Msg {
		return batchMsg(cmds)
	}
}

type batchMsg []Cmd

// Tick returns a command that sends a TickMsg after the given duration
func Tick(d time.Duration) Cmd {
	return func() Msg {
		time.Sleep(d)
		return TickMsg{Time: time.Now()}
	}
}

// Quit returns a command that quits the program
func Quit() Msg {
	return QuitMsg{}
}

// Sequence executes commands in order
func Sequence(cmds ...Cmd) Cmd {
	return func() Msg {
		return sequenceMsg(cmds)
	}
}

type sequenceMsg []Cmd

// Every returns a command that sends a TickMsg on a regular interval
func Every(d time.Duration, fn func(time.Time) Msg) Cmd {
	return func() Msg {
		time.Sleep(d)
		return fn(time.Now())
	}
}
