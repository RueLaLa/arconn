// Copyright 2018 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package websocketutil contains methods for interacting with websocket connections.
package websocketutil

import (
	"errors"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// IWebsocketUtil is the interface for the websocketutil.
type IWebsocketUtil interface {
	OpenConnection(url string) (*websocket.Conn, error)
	CloseConnection(ws websocket.Conn) error
}

// WebsocketUtil struct provides functionality around creating and maintaining websockets.
type WebsocketUtil struct {
	dialer *websocket.Dialer
}

// NewWebsocketUtil is the factory function for websocketutil.
func NewWebsocketUtil(dialerInput *websocket.Dialer) *WebsocketUtil {

	var websocketUtil *WebsocketUtil

	if dialerInput == nil {
		websocketUtil = &WebsocketUtil{
			dialer: websocket.DefaultDialer,
		}
	} else {
		websocketUtil = &WebsocketUtil{
			dialer: dialerInput,
		}
	}

	return websocketUtil
}

// OpenConnection opens a websocket connection provided an input url.
func (u *WebsocketUtil) OpenConnection(url string) (*websocket.Conn, error) {

	log.Infof("Opening websocket connection to: ", url)

	conn, _, err := u.dialer.Dial(url, nil)
	if err != nil {
		log.Errorf("Failed to dial websocket: %s", err.Error())
		return nil, err
	}

	log.Infof("Successfully opened websocket connection to: ", url)

	return conn, err
}

// CloseConnection closes a websocket connection given the Conn object as input.
func (u *WebsocketUtil) CloseConnection(ws *websocket.Conn) error {

	if ws == nil {
		return errors.New("websocket conn object is nil")
	}

	log.Debugf("Closing websocket connection to:", ws.RemoteAddr().String())

	err := ws.Close()
	if err != nil {
		log.Errorf("Failed to close websocket: %s", err.Error())
		return err
	}

	log.Debugf("Successfully closed websocket connection to:", ws.RemoteAddr().String())

	return nil
}
