package main

import (
	"encoding/json"
	"log"
)

// ChatResponse struct for chat response
type ChatResponse struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

// String return displayable string response
func (c *ChatResponse) String() string {
	jsonByte, errJSON := json.Marshal(c)
	if errJSON != nil {
		log.Println("unable to json marshall: " + errJSON.Error())
		return ""
	}
	return string(jsonByte[:])
}

// NewServerChatResponse instanciate new server response
func NewServerChatResponse(message string) *ChatResponse {
	return &ChatResponse{From: "ChatServer", Message: message}
}

// NewBotChatResponse instanciate new bot response
func NewBotChatResponse(message string) *ChatResponse {
	return &ChatResponse{From: "Bot", Message: message}
}

// BotChatWriter struct for bot writer
type BotChatWriter struct {
	chain *Chain
}

// NewBotChatWriter instanciate new bot chat writer
func NewBotChatWriter(chain *Chain) *BotChatWriter {
	return &BotChatWriter{chain}
}

// Write writes in to chain
func (c *BotChatWriter) Write(in []byte) (int, error) {
	var cr ChatResponse
	var err = json.Unmarshal(in, &cr)
	if nil != err {
		log.Println("unable to json unmarshall: " + string(in[:]))
		return 0, err
	}
	c.chain.Write([]byte(cr.Message))
	return len(in), nil
}
