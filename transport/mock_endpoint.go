package transport

import (
	"fmt"
)

// Something with the EndPointI interface - for use in testing.

type MockEndPoint struct {
	T    string   // transport
	Addr AddressI // usually a MockAddress
}

func NewMockEndPoint(transport, addr string) EndPointI {
	mockA := NewMockAddress(addr)
	return &MockEndPoint{transport, mockA}
}

func (m *MockEndPoint) Clone() (EndPointI, error) {
	klone := MockEndPoint{m.T, NewMockAddress(m.Addr.String())}
	return &klone, nil
}
func (m *MockEndPoint) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == m {
		return true
	}
	switch v := any.(type) {
	case *MockEndPoint:
		_ = v
	default:
		return false
	}
	other := any.(*MockEndPoint)
	return other.T == other.T &&
		other.Addr.String() == other.Addr.String()
}
func (m *MockEndPoint) Address() AddressI {
	return m.Addr
}
func (m *MockEndPoint) Transport() string {
	return m.T
}
func (m *MockEndPoint) String() string {
	return fmt.Sprintf("MockEndPoint: %s, %s", m.T, m.Addr.String())
}
