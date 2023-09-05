package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"luqchain/x/luqchain/types"
)

var _ = strconv.Itoa(0)

func CmdRetrieveVal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "retrieve-val [qid] [timestamp]",
		Short: "params are query id (hash of query data) and timestamp of when the report was submitted",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reqQid := args[0]
			reqTimestamp, err := cast.ToUint64E(args[1])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRetrieveValRequest{

				Qid:       reqQid,
				Timestamp: reqTimestamp,
			}

			res, err := queryClient.RetrieveVal(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
