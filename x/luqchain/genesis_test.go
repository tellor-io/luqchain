package luqchain_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	keepertest "luqchain/testutil/keeper"
	"luqchain/testutil/nullify"
	"luqchain/x/luqchain"
	"luqchain/x/luqchain/types"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LuqchainKeeper(t)
	luqchain.InitGenesis(ctx, *k, genesisState)
	got := luqchain.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
