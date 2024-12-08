package smp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/session-manager-plugin/pkg/log"
	"github.com/coder/websocket"
)

func StartSession(sessionInfo, target_json string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	startSessionRequest := make(map[string]interface{})
	json.Unmarshal([]byte(sessionInfo), &startSessionRequest)
	c, resp, err := websocket.Dial(ctx, startSessionRequest["StreamUrl"].(string), nil)
	log.Always(resp.Status)
	if err != nil {
		log.Error(err.Error())
	}
	defer c.CloseNow()
	mtype, message, _ := c.Read(context.Background())
	log.Alwaysf("%v", mtype)
	log.Always(string(message))
}
