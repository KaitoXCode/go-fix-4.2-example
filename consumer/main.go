package consumer

import (
	"fmt"

	"github.com/quickfixgo/quickfix"
)

// Consumer implements quickfix.Application interface
type Consumer struct {}

// OnCreate
func (p Consumer) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnCreate]: Creating session: %s", sessionID))
}

// OnLogon
func (p Consumer) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogon]: Starting session: %s", sessionID))
}

// OnLogout
func (p Consumer) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [OnLogout]: Terminating session: %s", sessionID))
}

// FromAdimn
func (p Consumer) FromAdimn(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromApp]: Received: %s", msg.String()))
  return nil
}

// ToAdimn
func (p Consumer) ToAdimn(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToAdimn]: Sending: %s", msg.String()))
}

// ToApp
func (p Consumer) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [ToApp]: Sending: %s", msg.String()))
  return
}

// FromApp
func (p Consumer) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [CONSUMER] [INFO] [FromApp]: Received: %s", msg.String()))
  return
}
