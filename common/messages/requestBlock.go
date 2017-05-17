// Copyright 2017 Factom Foundation
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package messages

import (
	"fmt"

	"github.com/FactomProject/factomd/common/constants"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
)

//A placeholder structure for messages
type RequestBlock struct {
	MessageBase
	Timestamp interfaces.Timestamp

	//TODO: figure whether this should be signed or not?

	//Not marshalled
	hash interfaces.IHash
}

var _ interfaces.IMsg = (*RequestBlock)(nil)

func (a *RequestBlock) IsSameAs(b *RequestBlock) bool {
	if b == nil {
		return false
	}
	if a.Timestamp.GetTimeMilli() != b.Timestamp.GetTimeMilli() {
		return false
	}

	//TODO: expand

	return true
}

func (m *RequestBlock) Process(uint32, interfaces.IState) bool { return true }

func (m *RequestBlock) GetRepeatHash() interfaces.IHash {
	return m.GetMsgHash()
}

func (m *RequestBlock) GetHash() interfaces.IHash {
	if m.hash == nil {
		data, err := m.MarshalForSignature()
		if err != nil {
			panic(fmt.Sprintf("Error in RequestBlock.GetHash(): %s", err.Error()))
		}
		m.hash = primitives.Sha(data)
	}
	return m.hash
}

func (m *RequestBlock) GetMsgHash() interfaces.IHash {
	if m.MsgHash == nil {
		data, err := m.MarshalBinary()
		if err != nil {
			return nil
		}
		m.MsgHash = primitives.Sha(data)
	}
	return m.MsgHash
}

func (m *RequestBlock) GetTimestamp() interfaces.Timestamp {
	return m.Timestamp
}

func (m *RequestBlock) Type() byte {
	return constants.REQUEST_BLOCK_MSG
}

func (m *RequestBlock) UnmarshalBinaryData(data []byte) ([]byte, error) {
	buf := primitives.NewBuffer(data)

	t, err := buf.PopByte()
	if err != nil {
		return nil, err
	}
	if t != m.Type() {
		return nil, fmt.Errorf("Invalid Message type")
	}

	m.Timestamp = new(primitives.Timestamp)
	err = buf.PopBinaryMarshallable(m.Timestamp)
	if err != nil {
		return nil, err
	}

	//TODO: expand

	return buf.DeepCopyBytes(), nil
}

func (m *RequestBlock) UnmarshalBinary(data []byte) error {
	_, err := m.UnmarshalBinaryData(data)
	return err
}

func (m *RequestBlock) MarshalForSignature() ([]byte, error) {
	buf := primitives.NewBuffer(nil)
	err := buf.PushByte(m.Type())
	if err != nil {
		return nil, err
	}
	err = buf.PushBinaryMarshallable(m.Timestamp)
	if err != nil {
		return nil, err
	}

	//TODO: expand

	return buf.DeepCopyBytes(), nil
}

func (m *RequestBlock) MarshalBinary() (data []byte, err error) {
	//TODO: sign or delete
	return m.MarshalForSignature()
}

func (m *RequestBlock) String() string {
	return "Request Block"
}

func (m *RequestBlock) DBHeight() int {
	return 0
}

func (m *RequestBlock) ChainID() []byte {
	return nil
}

func (m *RequestBlock) ListHeight() int {
	return 0
}

func (m *RequestBlock) SerialHash() []byte {
	return nil
}

func (m *RequestBlock) Signature() []byte {
	return nil
}

// Validate the message, given the state.  Three possible results:
//  < 0 -- Message is invalid.  Discard
//  0   -- Cannot tell if message is Valid
//  1   -- Message is valid
func (m *RequestBlock) Validate(state interfaces.IState) int {
	return 0
}

func (m *RequestBlock) ComputeVMIndex(state interfaces.IState) {
}

func (m *RequestBlock) LeaderExecute(state interfaces.IState) {
}

func (m *RequestBlock) FollowerExecute(interfaces.IState) {
}

func (e *RequestBlock) JSONByte() ([]byte, error) {
	return primitives.EncodeJSON(e)
}

func (e *RequestBlock) JSONString() (string, error) {
	return primitives.EncodeJSONString(e)
}
