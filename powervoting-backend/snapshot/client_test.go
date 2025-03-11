package snapshot_test

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"powervoting-server/config"
	"powervoting-server/snapshot"
	pb "powervoting-server/snapshot/proto"
	"powervoting-server/utils"
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
	power, err := snapshot.ParseAddressPowerResponse(res)
	// Assert expected results
	assert.NoError(t, err, "Expected no error")
	assert.NotNil(t, power.SpPower, "Expected SpPower to be set")
	assert.NotNil(t, power.ClientPower, "Expected ClientPower to be set")
	assert.NotNil(t, power.TokenHolderPower, "Expected TokenHolderPower to be set")
	assert.NotNil(t, power.DeveloperPower, "Expected DeveloperPower to be set")
	assert.Equal(t, int64(1000), power.BlockHeight.Int64(), "Expected BlockHeight to be 1000")
}

// TestGetAddressPowerByDay is a test function to verify the functionality of GetAddressPowerByDay method.
func TestGetAddressPowerByDay(t *testing.T) {
	// Initialize the logger for logging purposes.
	config.InitLogger()
	// Initialize the configuration by loading the config file from the specified path.
	config.InitConfig("../")

	// Set the ABI (Application Binary Interface) file paths for PowerVoting and Oracle contracts.
	config.Client.ABIPath.PowerVotingAbi = "../abi/power-voting.json"
	config.Client.ABIPath.OracleAbi = "../abi/oracle.json"

	// Call the GetAddressPowerByDay function with specified parameters:
	// - Chain ID: 314159
	// - Address: "0x763D410594a24048537990dde6ca81c38CfF566a"
	// - Date: "20250224"
	res, err := snapshot.GetAddressPowerByDay(314159, "0x763D410594a24048537990dde6ca81c38CfF566a", "20250105")
	// Assert that there is no error returned from the function call.
	assert.NoError(t, err)
	// Assert that the result is not nil, ensuring that the function returned a valid response.
	assert.NotNil(t, res)
	// Print the powers obtained from the result, converting them from a base of 10^18 to a more readable format.
	fmt.Printf("sp power: %s\t client power: %s\t token holder power: %s\t developer power: %s\n",
		utils.DividedBy10To18(res.SpPower, 4), // Convert SpPower to a readable format
		utils.DividedBy10To18(res.ClientPower, 4),
		utils.DividedBy10To18(res.TokenHolderPower, 4), // Convert TokenHolderPower to a readable format
		utils.DividedBy10To18(res.DeveloperPower, 4),   // Convert DeveloperPower to a readable format
	)
}
