package transport

// xlattice_go/mock_address.go

type MockAddress struct {
	Address string
}

func NewMockAddress( s string ) (AddressI) {
	return &MockAddress{ s }
}
func (m *MockAddress) Clone() (AddressI, error) {
	m2 := MockAddress{ m.Address }
	return &m2, nil
}
func (m *MockAddress) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == m {
		return true
	}
	switch v := any.(type) {
	case *MockAddress:
		_ = v
	default:
		return false
	}
	other := any.(*MockAddress)
	return m.Address == other.Address
}
func (m *MockAddress) String() string {
	return m.Address
}
