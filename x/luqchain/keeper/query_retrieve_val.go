package keeper

import (
	"context"
	"encoding/hex"

	"luqchain/x/luqchain/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RetrieveVal(goCtx context.Context, req *types.QueryRetrieveValRequest) (*types.QueryRetrieveValResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	var report types.Report
	ctx := sdk.UnwrapSDKContext(goCtx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReportKey))
	qid, err := hex.DecodeString(req.Qid)
	if err != nil {
		ctx.Logger().Error("invalid qid", "qid", req.Qid)
		return nil, err
	}
	tbytes := Uint64ToBytes(req.Timestamp)
	b := store.Get(append(qid, tbytes...))
	if b == nil {
		return nil, sdkerrors.ErrKeyNotFound
	}
	k.cdc.MustUnmarshal(b, &report)

	return &types.QueryRetrieveValResponse{Report: &report}, nil
}
