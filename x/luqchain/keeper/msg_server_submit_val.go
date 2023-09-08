package keeper

import (
	"context"
	"fmt"

	"luqchain/x/luqchain/types"

	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SubmitVal(goCtx context.Context, msg *types.MsgSubmitVal) (*types.MsgSubmitValResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddr, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		ctx.Logger().Error("invalid account address", "address", accAddr)
		return nil, err
	}
	val := sdk.ValAddress(accAddr)

	if !k.staking.Validator(ctx, val).IsBonded() {
		return nil, fmt.Errorf("validator is not bonded")
	}
	// reporter has to be validator
	var report = types.Report{
		Qdata:     msg.Qdata,
		Value:     msg.Value,
		Timestamp: uint64(ctx.BlockTime().Unix()),
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ReportKey))
	addedReport := k.cdc.MustMarshal(&report)
	qid, err := HashQdata(msg.Qdata)
	if err != nil {
		ctx.Logger().Error("invalid qdata", "qdata", msg.Qdata)
		return nil, err
	}
	k.Logger(ctx).Error("submit report", "report", bytes.HexBytes(qid).String())
	tbytes := Uint64ToBytes(report.Timestamp)
	store.Set(append(qid, tbytes...), addedReport)
	return &types.MsgSubmitValResponse{}, nil
}
