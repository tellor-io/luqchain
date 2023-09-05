package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"luqchain/x/luqchain/types"
)

var _ = strconv.Itoa(0)

func CmdSubmitVal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-val [qdata] [value]",
		Short: "params are query data in hex string and value as uint",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argQdata := args[0]
			argValue, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgSubmitVal(
				clientCtx.GetFromAddress().String(),
				argQdata,
				argValue,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
