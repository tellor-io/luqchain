package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	testkeeper "luqchain/testutil/keeper"
	"luqchain/x/luqchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LuqchainKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
