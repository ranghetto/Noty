package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

/*
	STRUCT
*/

// Handler : helpers functions
type Handler struct {
	Session   *discordgo.Session
	ChannelID string
}

/*
	SETTERS
*/

// SetSession : set session
func (h *Handler) SetSession(session *discordgo.Session) {
	h.Session = session
}

// SetChannelID : set channel id
func (h *Handler) SetChannelID(channelID string) {
	h.ChannelID = channelID
}

/*
	PUBLIC
*/

// SendLogError : send a detailed error message to the channel
func (h Handler) SendLogError(e error) {
	if e != nil {
		h.deleteLastMessage()
		_, err := h.Session.ChannelMessageSend(h.ChannelID, e.Error())
		FatalError(err)
		LogError(e)
	}
}

// CreateValidChannel : create a valid channel to start a new TODO LIST
func (h Handler) CreateValidChannel() {
	messages, err := h.Session.ChannelMessages(h.ChannelID, 0, "", "", "")
	LogError(err)

	for i := 0; i < len(messages); i++ {
		go h.DeleteMessage(messages[i].ID)
	}
}

// UpdateChannelList : update the todo listin the channel
func (h Handler) UpdateChannelList(new bool) {
	if !new {
		h.deleteLastMessage()
	}
	message := h.TodoList()
	if len(message) == 0 {
		message = "Tip: **ADD** a new **TODO** using `/todo add <text>`"
	}
	_, err := h.Session.ChannelMessageSend(h.ChannelID, message)
	LogError(err)
}

// TodoAdd : add a new todo to the database and return true if added or false if not
func (h Handler) TodoAdd(text string) bool {

	query, err := database.Prepare("INSERT INTO todo (text, done, channel) VALUES (?, ?, ?)")
	if err != nil {
		h.SendLogError(err)
		return false
	}

	res, err := query.Exec(text, 0, h.ChannelID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	_, err = res.LastInsertId()
	if err != nil {
		h.SendLogError(err)
		return false
	}

	return true
}

// TodoList : return a string with a list of todos
func (h Handler) TodoList() string {

	query := "SELECT id, text, done FROM todo WHERE channel = ?"

	var listMessage string

	rows, err := database.Query(query, h.ChannelID)
	h.SendLogError(err)

	var id int
	var text string
	var done int

	/*if !rows.Next() {
		listMessage = "Create a new todo with `/todo add <text>` to populate this list"
	}*/
	for rows.Next() {

		var todo string

		err = rows.Scan(&id, &text, &done)
		h.SendLogError(err)

		if done == 1 {
			todo = fmt.Sprintf("~~%d ~ %s~~\n", id, text)
		} else {
			todo = fmt.Sprintf("%d ~ %s\n", id, text)
		}

		listMessage += todo
	}

	rows.Close()

	return listMessage
}

// TodoUpdate : update a todo and return true if updated or false if not
func (h Handler) TodoUpdate(text string, todoID string) bool {

	query, err := database.Prepare("UPDATE todo SET text = ? WHERE id = ?")
	if err != nil {
		h.SendLogError(err)
		return false
	}

	res, err := query.Exec(text, todoID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	rows, err := res.RowsAffected()
	if err != nil {
		h.SendLogError(err)
		return false
	}

	if rows != 1 {
		return false
	}

	return true
}

// TodoDone : mark a todo as done and return true if updated or false if not
func (h Handler) TodoDone(todoID string) bool {

	query, err := database.Prepare("UPDATE todo SET done = 1 WHERE id = ?")
	if err != nil {
		h.SendLogError(err)
		return false
	}

	res, err := query.Exec(todoID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	rows, err := res.RowsAffected()
	if err != nil {
		h.SendLogError(err)
		return false
	}

	if rows != 1 {
		return false
	}

	return true
}

// TodoNotDone : mark a todo as not done and return true if updated or false if not
func (h Handler) TodoNotDone(todoID string) bool {

	query, err := database.Prepare("UPDATE todo SET done = 0 WHERE id = ?")
	if err != nil {
		h.SendLogError(err)
		return false
	}

	res, err := query.Exec(todoID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	rows, err := res.RowsAffected()
	if err != nil {
		h.SendLogError(err)
		return false
	}

	if rows != 1 {
		return false
	}

	return true
}

// TodoDelete : delete a todo
func (h Handler) TodoDelete(todoID string) bool {

	query, err := database.Prepare("DELETE FROM todo WHERE channel = ? AND id = ?")
	if err != nil {
		h.SendLogError(err)
		return false
	}
	_, err = query.Exec(h.ChannelID, todoID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	return true
}

// TodoReset : reset todo table
func (h Handler) TodoReset() bool {

	query, err := database.Prepare("DELETE FROM todo WHERE channel = ?")
	if err != nil {
		h.SendLogError(err)
		return false
	}
	_, err = query.Exec(h.ChannelID)
	if err != nil {
		h.SendLogError(err)
		return false
	}

	return true
}

// DeleteMessage : delete given message
func (h Handler) DeleteMessage(messageID string) {
	err := h.Session.ChannelMessageDelete(h.ChannelID, messageID)
	LogError(err)
}

// ChannelIsValid : check list of valid channels in config
func (h Handler) ChannelIsValid(cf Config) bool {
	for i := 0; i < len(cf.Channels); i++ {
		if h.ChannelID == cf.Channels[i] {
			return true
		}
	}
	return false
}

/*
	PRIVATE
*/

func (h Handler) deleteLastMessage() {
	messages, err := h.Session.ChannelMessages(h.ChannelID, 1, "", "", "")
	LogError(err)

	h.DeleteMessage(messages[0].ID)
}
