package main

import (
	"encoding/json"
	"log"
)

type ChatResponse struct {
	From    string `json:"from"`
	Message string `json:"message"`
}

func (c *ChatResponse) String() string {
	jsonByte, errJSON := json.Marshal(c)
	if errJSON != nil {
		log.Println("unable to json marshall: " + errJSON.Error())
		return ""
	}
	return string(jsonByte[:])
}

func NewServerChatResponse(message string) *ChatResponse {
	return &ChatResponse{From: "ChatServer", Message: message}
}
func NewBotChatResponse(message string) *ChatResponse {
	return &ChatResponse{From: "Bot", Message: message}
}

type BotChatWriter struct {
	chain *Chain
}

func NewBotChatWriter(chain *Chain) *BotChatWriter {
	return &BotChatWriter{chain}
}

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
