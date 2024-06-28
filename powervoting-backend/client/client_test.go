package client_test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	pb "powervoting-server/client/proto"
	"testing"

	"github.com/stretchr/testify/assert"
	clientPkg "powervoting-server/client"
)

const bufSize = 1024 * 1024 // 1 MB buffer

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func setup() (*grpc.ClientConn, pb.SnapshotClient, func()) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterSnapshotServer(s, &mockSnapshotServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial bufnet: %v", err)
	}
	client := pb.NewSnapshotClient(conn)

	return conn, client, func() {
		err := conn.Close()
		if err != nil {
			return
		}
		err = lis.Close()
		if err != nil {
			return
		}
		s.Stop()
	}
}

type mockSnapshotServer struct {
	pb.UnimplementedSnapshotServer
}

func (s *mockSnapshotServer) GetAddressPower(_ context.Context, req *pb.AddressPowerRequest) (*pb.AddressPowerResponse, error) {
	return &pb.AddressPowerResponse{
		Address:          req.Address,
		SpPower:          "100",
		ClientPower:      "50",
		TokenHolderPower: "75",
		DeveloperPower:   "25",
		BlockHeight:      100000,
		DateStr:          "2024-06-19",
	}, nil
}

func TestMockGetAddressPower(t *testing.T) {
	_, client, teardown := setup()
	defer teardown()

	// Prepare the request.
	req := &pb.AddressPowerRequest{
		NetId:     314159,
		Address:   "0xBc27ca842D22cD5BdBC41B27A571EC1FbB559307",
		RandomNum: 1,
	}

	// Call the GetAddressPower function directly.
	res, err := client.GetAddressPower(context.Background(), req)
	if err != nil {
		t.Fatalf("GetAddressPower returned error: %v", err)
	}

	// Assert the expected response values.
	assert.Equal(t, req.Address, res.Address, "Address does not match")
	assert.Equal(t, "100", res.SpPower, "SpPower does not match")
	assert.Equal(t, "50", res.ClientPower, "ClientPower does not match")
	assert.Equal(t, "75", res.TokenHolderPower, "TokenHolderPower does not match")
	assert.Equal(t, "25", res.DeveloperPower, "DeveloperPower does not match")
	assert.Equal(t, int64(100000), res.BlockHeight, "BlockHeight does not match")
	assert.Equal(t, "2024-06-19", res.DateStr, "DateStr does not match")
}

func TestParseAddressPowerResponse(t *testing.T) {
	// Mock AddressPowerResponse
	res := &pb.AddressPowerResponse{
		SpPower:          "100",
		ClientPower:      "50",
		TokenHolderPower: "75",
		DeveloperPower:   "25",
		BlockHeight:      1000,
	}

	// Call parseAddressPowerResponse
	power, err := clientPkg.ParseAddressPowerResponse(res)
	// Assert expected results
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, power.SpPower, "Expected SpPower to be set")
	assert.NotNil(t, power.ClientPower, "Expected ClientPower to be set")
	assert.NotNil(t, power.TokenHolderPower, "Expected TokenHolderPower to be set")
	assert.NotNil(t, power.DeveloperPower, "Expected DeveloperPower to be set")
	assert.Equal(t, int64(1000), power.BlockHeight.Int64(), "Expected BlockHeight to be 1000")
}