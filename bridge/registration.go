package bridge

import (
	"context"
	. "luqchain/bridge/types"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/cosmos-sdk/client"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// to check queryServer implements ServiceServer
var _ BridgeServer = bridgeServer{}

// queryServer implements ServiceServer
type bridgeServer struct {
	clientCtx client.Context
}

// NewQueryServer returns new queryServer from provided client.Context
func NewQueryServer(clientCtx client.Context) BridgeServer {
	return bridgeServer{
		clientCtx: clientCtx,
	}
}

func RegisterHeaderService(clientCtx client.Context, server gogogrpc.Server) {
	RegisterBridgeServer(server, NewQueryServer(clientCtx))
}

// RegisterGRPCGatewayRoutes mounts the node gRPC service's GRPC-gateway routes
// on the given mux object.
func RegisterGRPCGatewayRoutes(clientConn gogogrpc.ClientConn, mux *runtime.ServeMux) {
	RegisterBridgeHandlerClient(context.Background(), mux, NewBridgeClient(clientConn))
}

func (s bridgeServer) getCommit(height int64) (*coretypes.ResultCommit, error) {
	node, err := s.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}
	var h *int64

	if height != 0 {
		h = &height
	}
	commit, err := node.Commit(context.Background(), h)
	return commit, err
}
