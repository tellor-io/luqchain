package keeper

import (
	"context"

	"luqchain/x/luqchain/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) RetrieveAll(goCtx context.Context, req *types.QueryRetrieveAllRequest) (*types.QueryRetrieveAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	var reports []types.Report
	store := ctx.KVStore(k.storeKey)
	reportStore := prefix.NewStore(store, types.KeyPrefix(types.ReportKey))
	pageRes, err := query.Paginate(reportStore, req.Pagination, func(key []byte, value []byte) error {
		var report types.Report
		if err := k.cdc.Unmarshal(value, &report); err != nil {
			return err
		}
		reports = append(reports, report)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryRetrieveAllResponse{Report: reports, Pagination: pageRes}, nil
}
