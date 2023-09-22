package app

import (
	"fmt"
	luqchaintypes "luqchain/x/luqchain/types"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
)

var EmptyResponse = abci.ResponsePrepareProposal{Txs: [][]byte{}}

// EncodeMsgsIntoTxBytes encodes the given msgs into a single transaction.
func EncodeMsgsIntoTxBytes(txConfig client.TxConfig, msgs ...sdk.Msg) ([]byte, error) {
	txBuilder := txConfig.NewTxBuilder()
	err := txBuilder.SetMsgs(msgs...)
	if err != nil {
		return nil, err
	}
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return nil, err
	}

	return txBytes, nil
}

func PrepareProposalHandler(
	txConfig client.TxConfig,
	txVerifier baseapp.ProposalTxVerifier,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		fmt.Println("Start PrepareProposalHandler!")
		decoder := txConfig.TxDecoder()
		// submitValuetxs := [][]byte{}
		// list of transactions
		// for each transaction, decode it and if msg name == MsgSubmitVal, then
		// add to mapping of qdata -> list of values to be averaged
		mapping := make(map[string][]uint64)
		// print num txs in req
		fmt.Println("num txs in req:", len(req.Txs))
		for _, tx := range req.Txs {

			txDec, err := decoder(tx)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("TxDecoder error: %v", err))
				return EmptyResponse
			}
			_, err = txVerifier.PrepareProposalVerifyTx(txDec)
			if err != nil {
				panic(err)
			}
			validateErr := txDec.ValidateBasic()
			if validateErr != nil {
				fmt.Println("TxDecoder ValidateBasic error: ")
				fmt.Println(validateErr)
				ctx.Logger().Error(fmt.Sprintf("TxDecoder ValidateBasic error: %v", validateErr))
				return EmptyResponse
			}

			msgs := txDec.GetMsgs()
			for _, msg := range msgs {
				funcName := proto.MessageName(msg)
				if funcName == "luqchain.luqchain.MsgSubmitVal" {
					// submitValuetxs = append(submitValuetxs, tx)
					msgSubmitVal := msg.(*luqchaintypes.MsgSubmitVal)
					fmt.Println("creator:", msgSubmitVal.Creator)
					fmt.Println("qdata:", msgSubmitVal.Qdata)
					fmt.Println("value:", msgSubmitVal.Value)
					mapping[msgSubmitVal.Qdata] = append(mapping[msgSubmitVal.Qdata], msgSubmitVal.Value)
				}
			}
		}

		fmt.Println(mapping)
		// for each key in mapping, assemble a MsgSubmitVal with the median value
		// and add to transactions to be returned.
		txsToReturn := [][]byte{}
		for _, values := range mapping {
			median := values[len(values)/2]
			addr := sdk.AccAddress(req.ProposerAddress)
			fmt.Println("proposer address:", addr)
			msgSubmitVal := luqchaintypes.NewMsgSubmitVal(addr.String(), "spot", median) //change key from qdata to avoid temp
			if err := msgSubmitVal.ValidateBasic(); err != nil {
				fmt.Println("MsgSubmitVal ValidateBasic error: ")
				fmt.Println(err)
				ctx.Logger().Error(fmt.Sprintf("MsgSubmitVal ValidateBasic error: %v", err))
				return EmptyResponse
			}
			txBytes, err := EncodeMsgsIntoTxBytes(txConfig, msgSubmitVal)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("EncodeMsgsIntoTxBytes error: %v", err))
				return EmptyResponse
			}
			txDec, err := decoder(txBytes)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("TxDecoder error: %v", err))
				return EmptyResponse
			}
			// transaction will fail to verify if it is not signed
			_, err = txVerifier.PrepareProposalVerifyTx(txDec)
			if err != nil {
				ctx.Logger().Error(fmt.Sprintf("PrepareProposalVerifyTx error: %v", err))
				return EmptyResponse
			}
			txsToReturn = append(txsToReturn, txBytes)
		}
		fmt.Println("End PrepareProposalHandler")
		fmt.Println(len(txsToReturn))
		return abci.ResponsePrepareProposal{Txs: txsToReturn}
	}
}
