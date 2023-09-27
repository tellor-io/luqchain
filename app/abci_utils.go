package app

import (
	"fmt"
	luqchaintypes "luqchain/x/luqchain/types"
	"math"
	"sort"

	cosmosmath "cosmossdk.io/math"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
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

type MsgData struct {
	signer string
	value  uint64
}

func calcWeightedMean(values []uint64, weights []uint64) uint64 {
	var sum, weightsSum uint64
	for i, value := range values {
		sum += value * weights[i]
		weightsSum += weights[i]
	}
	return sum / weightsSum
}

func calcWeightedVariance(values []uint64, weights []uint64, weightedMean uint64) float64 {
	var sum, weightsSum float64
	for i, value := range values {
		sum += float64(weights[i]) * math.Pow(float64(value-weightedMean), 2)
		weightsSum += float64(weights[i])
	}
	return sum / weightsSum
}

func calcWeightedStandardDeviation(values []uint64, weights []uint64, weightedMean uint64) uint64 {
	return uint64(math.Sqrt(calcWeightedVariance(values, weights, weightedMean)))
}

func calcWeightedMedian(values []uint64, weights []uint64) uint64 {
	if len(values) == 0 || len(weights) == 0 || len(values) != len(weights) {
		return 0 // err instead
	}

	type pair struct {
		value  uint64
		weight uint64
	}

	pairs := make([]pair, len(values))
	for i := range values {
		pairs[i] = pair{values[i], weights[i]}
	}

	sort.Slice(pairs, func(i, j int) bool { return pairs[i].value < pairs[j].value })

	totalWeight := 0
	for _, weight := range weights {
		totalWeight += int(weight)
	}

	accumulatedWeight := 0
	for _, p := range pairs {
		accumulatedWeight += int(p.weight)
		if accumulatedWeight >= totalWeight/2 {
			return p.value
		}
	}

	return 0 // should not reach this line
}

func PrepareProposalHandler(
	txConfig client.TxConfig,
	txVerifier baseapp.ProposalTxVerifier,
	k *stakingkeeper.Keeper,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		fmt.Println("Start PrepareProposalHandler!")
		decoder := txConfig.TxDecoder()

		// submitValuetxs := [][]byte{}
		// list of transactions
		// for each transaction, decode it and if msg name == MsgSubmitVal, then
		// add to mapping of qdata -> list of values to be averaged
		mapping := make(map[string][]MsgData)
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
					mapping[msgSubmitVal.Qdata] = append(
						mapping[msgSubmitVal.Qdata],
						MsgData{
							signer: msgSubmitVal.Creator,
							value:  msgSubmitVal.Value,
						},
					)
				}
			}
		}

		fmt.Println(mapping)
		// for each key (todo: make query id) in mapping, assemble a MsgSubmitVal with the median value
		// and add to transactions to be returned.
		txsToReturn := [][]byte{}
		for _, values := range mapping {
			weights := make([]uint64, len(values))
			vals := make([]uint64, len(values))
			for i, msgData := range values {
				fmt.Println("signer:", msgData.signer)
				signerAccAddress, _ := sdk.AccAddressFromBech32(msgData.signer)
				// fmt.Println("signerAccAddress:", signerAccAddress)
				fmt.Println("validator address:", sdk.ValAddress(signerAccAddress.Bytes()))
				validator := k.Validator(ctx, sdk.ValAddress(signerAccAddress.Bytes()))
				// fmt.Println("validator bonded status:", validator.IsBonded())
				fmt.Println("validator bonded tokens amount:", validator.GetBondedTokens())
				// print current block height
				// fmt.Println("current block height:", cosmosmath.NewInt(ctx.BlockHeight()))
				power := validator.GetConsensusPower(cosmosmath.NewInt(ctx.BlockHeight()))
				fmt.Println("power:", power)
				weights[i] = uint64(power)
				vals[i] = (msgData.value)
			}
			median := calcWeightedMedian(vals, weights)
			fmt.Println("weighted median:", median)
			stdDev := calcWeightedStandardDeviation(vals, weights, median)
			fmt.Println("weighted standard deviation:", stdDev)
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
