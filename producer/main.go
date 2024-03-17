package producer

import (
	"fmt"

	"github.com/quickfixgo/quickfix"
)

// Producer implements quickfix.Application interface
type Producer struct {}

// OnCreate
func (p Producer) OnCreate(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnCreate]: Creating session: %s", sessionID))
}

// OnLogon
func (p Producer) OnLogon(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnLogon]: Starting session: %s", sessionID))
}

// OnLogout
func (p Producer) OnLogout(sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [OnLogout]: Terminating session: %s", sessionID))
}

// FromAdimn
func (p Producer) FromAdimn(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [FromApp]: Received: %s", msg.String()))
  return nil
}

// ToAdimn
func (p Producer) ToAdimn(msg *quickfix.Message, sessionID quickfix.SessionID) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [ToAdimn]: Sending: %s", msg.String()))
}

// ToApp
func (p Producer) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [ToApp]: Sending: %s", msg.String()))
  return
}

// FromApp
func (p Producer) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
  println(fmt.Sprintf(" [PRODUCER] [INFO] [FromApp]: Received: %s", msg.String()))
  return
}
