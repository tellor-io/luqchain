package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	s "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	GetAccount(ctx sdk.Context, addr sdk.AccAddress) types.AccountI
	// Methods imported from account should be defined here
}

// BankKeeper defines the expected interface needed to retrieve account balances.
type BankKeeper interface {
	SpendableCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
	// Methods imported from bank should be defined here
}

type StakingKeeper interface {
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) s.ValidatorI
	GetAllValidators(ctx sdk.Context) []s.Validator
	TotalBondedTokens(ctx sdk.Context) sdk.Int
	Validator(ctx sdk.Context, addr sdk.ValAddress) s.ValidatorI
	// Methods imported from bank should be defined here
}
