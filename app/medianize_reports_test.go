package app_test

import (
	"fmt"
	"luqchain/app"
	"luqchain/app/params"
	"luqchain/x/luqchain/types"
	"testing"

	dbm "github.com/cometbft/cometbft-db"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/libs/log"

	abci "github.com/cometbft/cometbft/abci/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/server"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	simcli "github.com/cosmos/cosmos-sdk/x/simulation/client/cli"
	"github.com/cosmos/gogoproto/proto"

	"errors"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

type MyProposalTxVerifier struct {
	encoder client.TxConfig
}

// NewMyProposalTxVerifier creates a new instance of MyProposalTxVerifier
func NewMyProposalTxVerifier(encoder client.TxConfig) baseapp.ProposalTxVerifier {
	return &MyProposalTxVerifier{
		encoder: encoder,
	}
}

func (m *MyProposalTxVerifier) PrepareProposalVerifyTx(tx sdktypes.Tx) ([]byte, error) {
	return m.encoder.TxEncoder()(tx)
}

func (m *MyProposalTxVerifier) ProcessProposalVerifyTx(txBz []byte) (sdktypes.Tx, error) {
	return m.encoder.TxDecoder()(txBz)
}

func CreateSignedTxs(
	msgs []sdktypes.Msg,
	privKeys []cryptotypes.PrivKey,
	app *app.App,
	encCfg params.EncodingConfig,
	accNums []uint64,
	accSeqs []uint64,
	chainID string,
) ([][]byte, error) {
	if len(msgs) != len(privKeys) || len(privKeys) != len(accNums) || len(accNums) != len(accSeqs) {
		return nil, errors.New("mismatched lengths of msgs, privKeys, accNums, and accSeqs")
	}

	var txBytes [][]byte

	for i, msg := range msgs {
		txBuilder := app.TxConfig().NewTxBuilder()

		// Setting the message
		err := txBuilder.SetMsgs(msg)
		if err != nil {
			return nil, err
		}

		// Gathering all signer infos
		var sigsV2 []signing.SignatureV2
		for j, priv := range privKeys {
			sigV2 := signing.SignatureV2{
				PubKey: priv.PubKey(),
				Data: &signing.SingleSignatureData{
					SignMode:  app.TxConfig().SignModeHandler().DefaultMode(),
					Signature: nil,
				},
				Sequence: accSeqs[j],
			}

			sigsV2 = append(sigsV2, sigV2)
		}

		err = txBuilder.SetSignatures(sigsV2...)
		if err != nil {
			return nil, err
		}

		// sign tx with priv key at index i
		sigsV2 = []signing.SignatureV2{}
		signerData := xauthsigning.SignerData{
			ChainID:       chainID,
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}

		sigV2, err := tx.SignWithPrivKey(
			encCfg.TxConfig.SignModeHandler().DefaultMode(),
			signerData,
			txBuilder,
			privKeys[i],
			encCfg.TxConfig,
			accSeqs[i],
		)
		if err != nil {
			return nil, err
		}

		sigsV2 = append(sigsV2, sigV2)

		// Setting the signatures
		err = txBuilder.SetSignatures(sigsV2...)
		if err != nil {
			return nil, err
		}

		// get tx bytes
		encodedTx, err := app.TxConfig().TxEncoder()(txBuilder.GetTx())
		if err != nil {
			return nil, err
		}
		// Appending tx bytes
		txBytes = append(txBytes, encodedTx)
	}

	return txBytes, nil
}

func printSubmitValTxs(txs [][]byte, txConfig client.TxConfig) {
	for _, tx := range txs {
		txDec, err := txConfig.TxDecoder()(tx)
		if err != nil {
			fmt.Println("TxDecoder error: ", err)
		}
		msgs := txDec.GetMsgs()
		for _, msg := range msgs {
			funcName := proto.MessageName(msg)
			if funcName == "luqchain.luqchain.MsgSubmitVal" {
				fmt.Println("MsgSubmitVal signer: ", msg.GetSigners())
				msgSubmitVal := msg.(*types.MsgSubmitVal)
				fmt.Println("creator:", msgSubmitVal.Creator)
				fmt.Println("qdata:", msgSubmitVal.Qdata)
				fmt.Println("value:", msgSubmitVal.Value)
			}
		}
	}
}

func TestPrepareProposalHandler(t *testing.T) {
	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = app.DefaultNodeHome
	appOptions[server.FlagInvCheckPeriod] = simcli.FlagPeriodValue
	testApp := app.New(
		log.NewNopLogger(),
		dbm.NewMemDB(), nil, true, map[int64]bool{}, app.DefaultNodeHome, 0,
		app.MakeEncodingConfig(),
		appOptions)
	ctxA := testApp.NewContext(true, tmproto.Header{Height: testApp.LastBlockHeight()})

	txConfig := testApp.TxConfig()
	txVerifier := NewMyProposalTxVerifier(txConfig)

	priv1, _, addr1 := testdata.KeyTestPubAddr()
	priv2, _, addr2 := testdata.KeyTestPubAddr()
	priv3, _, addr3 := testdata.KeyTestPubAddr()

	msg1 := types.NewMsgSubmitVal(addr1.String(), "spot", 100)
	msg2 := types.NewMsgSubmitVal(addr2.String(), "spot", 200)
	msg3 := types.NewMsgSubmitVal(addr3.String(), "spot", 300)

	encodingConfig := app.MakeEncodingConfig()
	txs, err := CreateSignedTxs(
		[]sdktypes.Msg{msg1, msg2, msg3},
		[]cryptotypes.PrivKey{priv1, priv2, priv3},
		testApp,
		encodingConfig,
		[]uint64{0, 0, 0},
		[]uint64{0, 0, 0},
		"test-chain",
	)
	require.NoError(t, err)

	// printSubmitValTxs(txs, txConfig)

	handler := app.PrepareProposalHandler(txConfig, txVerifier)
	req := abci.RequestPrepareProposal{
		Txs:             txs,
		ProposerAddress: addr1.Bytes(),
	}
	resp := handler(ctxA, req)

	emptyResponse := abci.ResponsePrepareProposal{Txs: [][]byte{}}
	require.NotEqual(t, emptyResponse, resp, "Expected non-empty response from PrepareProposalHandler.")

	fmt.Println("medianized tx:")
	printSubmitValTxs(resp.Txs, txConfig)
}
