package main

import (
	"regexp"
	"strings"
)

// Command : command functions
type Command struct {
	Handler Handler
	Cmd     string
}

/*
	SETTERS
*/

// SetCommand : set command
func (c *Command) SetCommand(cmd string) {
	c.Cmd = cmd
}

// SetHandler : set handler
func (c *Command) SetHandler(handler Handler) {
	c.Handler = handler
}

/*
	PUBLIC
*/

// WaitAndHandleCommands : handle commands functions
func (c Command) WaitAndHandleCommands() {

	if c.commandIsValid() {

		// Creo un array di stringhe utilizzando come separatore lo spazio
		command := strings.Fields(c.Cmd)

		cmdLen := len(command)
		var cmd1 string
		var cmd2 string

		if cmdLen == 1 {
			cmd1 = strings.TrimPrefix(command[0], "/")
		}
		if cmdLen > 1 {
			cmd1 = strings.TrimPrefix(command[0], "/")
			cmd2 = command[1]
		}

		if cmd1 == "start" {
			h.CreateValidChannel()
			h.UpdateChannelList(true)
		}

		if cmd1 == "todo" {

			if cmd2 == "reset" {
				if h.TodoReset() {
					h.UpdateChannelList(false)
				}
			}

			if cmd2 == "add" {
				text := strings.Join(command[2:], " ")
				if h.TodoAdd(text) {
					h.UpdateChannelList(false)
				}
			}

			if cmd2 == "done" {
				todoID := command[2]
				if h.TodoDone(todoID) {
					h.UpdateChannelList(false)
				}
			}
			if cmd2 == "notdone" {
				todoID := command[2]
				if h.TodoNotDone(todoID) {
					h.UpdateChannelList(false)
				}
			}
			if cmd2 == "delete" {
				todoID := command[2]
				if h.TodoDelete(todoID) {
					h.UpdateChannelList(false)
				}
			}
			if cmd2 == "update" {
				todoID := command[2]
				text := strings.Join(command[3:], " ")
				if h.TodoUpdate(text, todoID) {
					h.UpdateChannelList(false)
				}
			}
		}
	}
}

/*
	PRIVATE
*/

// Controllo che il messaggio sia un comando formattato correttamente
func (c Command) commandIsValid() bool {
	match, err := regexp.MatchString(`/\w*`, c.Cmd)
	LogError(err)
	return match
}
