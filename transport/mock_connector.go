package transport

// xlattice_go/transport/mock_connector.go

import ()

//
type MockConnector struct {
	FarEnd *MockEndPoint
}

func NewMockConnector(farEnd EndPointI) (ctor *MockConnector, err error) {
	switch v := farEnd.(type) {
	case *MockEndPoint:
		_ = v
	default:
		return nil, NotMockEndPoint
	}
	mockFarEnd := farEnd.(*MockEndPoint)

	// copy the far end
	ep2, err := mockFarEnd.Clone()
	if err == nil {
		ctor = &MockConnector{ep2.(*MockEndPoint)}
	}
	return
}

// Establish a Connection with another entity using the transport
// and address in the EndPoint.
//
// @param nearEnd  local end point to use for connection
// @param blocking whether the new Connection is to be blocking
//
func (c *MockConnector) Connect(nearEnd EndPointI) (
	cnx ConnectionI, err error) {

	var (
		mockNearEnd *MockEndPoint
	)
	if nearEnd == nil {
		mockNearEnd = NewMockEndPoint("T", "A").(*MockEndPoint)
	} else {
		mockNearEnd = nearEnd.(*MockEndPoint)
	}
	if err == nil {
		cnx, err = NewMockConnection(mockNearEnd, c.FarEnd)
	}
	return
}

// return the Acceptor EndPoint that this Connector is used to
//          establish connections to
//
func (c *MockConnector) GetFarEnd() EndPointI {
	// XXX Should return copy
	return c.FarEnd
}

func (c *MockConnector) String() string {
	// FarEnd serialization begins with "MockEndPoint: "
	return "MockConnector: " + c.FarEnd.String()[13:]
}
