package bridge

import (
	"context"
	"encoding/hex"
	"encoding/json"
	fmt "fmt"
	. "luqchain/bridge/types"
	"luqchain/x/luqchain/keeper"
	"luqchain/x/luqchain/types"
	"math/big"

	"github.com/cometbft/cometbft/libs/bytes"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

var (
	relayArguments  abi.Arguments
	verifyArguments abi.Arguments
)

var pl = fmt.Println

func convertHeader(h BlockHeaderMerkle) BlockHeaderEvm {
	return BlockHeaderEvm{
		VersionAndChainIdHash:  common.BytesToHash(h.VersionChainidHash),
		Height:                 h.Height,
		TimeSecond:             h.TimeSecond,
		TimeNanoSecondFraction: h.TimeNanosecond,
		LastBlockIdCommitMerkleHash: common.BytesToHash(
			h.LastblockidCommitHash,
		),
		NextValidatorConsensusMerkleHash: common.BytesToHash(
			h.NextvalidatorConsensusHash,
		),
		LastResultsHash:            common.BytesToHash(h.LastresultsHash),
		EvidenceProposerMerkleHash: common.BytesToHash(h.EvidenceProposerHash),
	}
}
func convertStore(h MultiStoreTreeFields) MultistoreEvm {
	return MultistoreEvm{
		LuqchainIAVLStateHash:            common.BytesToHash(h.LuqchainIavlStateHash),
		MintStoreMerkleHash:              common.BytesToHash(h.MintStoreMerkleHash),
		IcacontrollerToIcahostMerkleHash: common.BytesToHash(h.IcacontrollerToIcahostMerkleHash),
		FeegrantToIbcMerkleHash:          common.BytesToHash(h.FeegrantToIbcMerkleHash),
		AccToEvidenceMerkleHash:          common.BytesToHash(h.AccToEvidenceMerkleHash),
		ParamsToVestingMerkleHash:        common.BytesToHash(h.ParamsToVestingMerkleHash),
	}
}
func convertSigs(sigs []TmSig) []TmSigEvm {
	var ret []TmSigEvm
	for _, sig := range sigs {
		ret = append(ret, TmSigEvm{
			R:                common.BytesToHash(sig.R),
			S:                common.BytesToHash(sig.S),
			V:                uint8(sig.V),
			EncodedTimestamp: sig.EncodedTimestamp,
		})
	}
	return ret
}
func convertIavl(iavl []IAVLMerklePath) []IavlEvm {
	var ret []IavlEvm
	for _, i := range iavl {
		ret = append(ret, IavlEvm{
			IsDataOnRight:  i.IsDataOnRight,
			SubtreeHeight:  uint8(i.SubtreeHeight),
			SubtreeSize:    big.NewInt(int64(i.SubtreeSize)),
			SubtreeVersion: big.NewInt(int64(i.SubtreeVersion)),
			SiblingHash:    common.BytesToHash(i.SiblingHash),
		})
	}
	return ret
}
func convertIavlToString(iavl []IAVLMerklePath) []IAVLMerklePathHex {
	var ret []IAVLMerklePathHex
	for _, i := range iavl {
		ret = append(ret, IAVLMerklePathHex{
			IsDataOnRight:  i.IsDataOnRight,
			SubtreeHeight:  i.SubtreeHeight,
			SubtreeSize:    i.SubtreeSize,
			SubtreeVersion: i.SubtreeVersion,
			SiblingHash:    common.BytesToHash(i.SiblingHash).Hex(),
		})
	}
	return ret
}
func convertMultistoretoHex(h MultiStoreTreeFields) MulitstoreHex {
	return MulitstoreHex{
		LuqchainIavlStateHash:            common.BytesToHash(h.LuqchainIavlStateHash).Hex(),
		MintStoreMerkleHash:              common.BytesToHash(h.MintStoreMerkleHash).Hex(),
		IcacontrollerToIcahostMerkleHash: common.BytesToHash(h.IcacontrollerToIcahostMerkleHash).Hex(),
		FeegrantToIbcMerkleHash:          common.BytesToHash(h.FeegrantToIbcMerkleHash).Hex(),
		AccToEvidenceMerkleHash:          common.BytesToHash(h.AccToEvidenceMerkleHash).Hex(),
		ParamsToVestingMerkleHash:        common.BytesToHash(h.ParamsToVestingMerkleHash).Hex(),
	}
}
func convertHeadertoHex(h BlockHeaderMerkle) BlockHeaderHex {
	return BlockHeaderHex{
		VersionChainidHash:         common.BytesToHash(h.VersionChainidHash).Hex(),
		Height:                     h.Height,
		TimeSecond:                 h.TimeSecond,
		TimeNanosecond:             h.TimeNanosecond,
		LastblockidCommitHash:      common.BytesToHash(h.LastblockidCommitHash).Hex(),
		NextvalidatorConsensusHash: common.BytesToHash(h.NextvalidatorConsensusHash).Hex(),
		LastresultsHash:            common.BytesToHash(h.LastresultsHash).Hex(),
		EvidenceProposerHash:       common.BytesToHash(h.EvidenceProposerHash).Hex(),
	}
}

func convertCommontoHex(c CommonEncodedVotePart) CommonHex {
	return CommonHex{
		SignedDataPrefix: bytes.HexBytes(c.SignedDataPrefix).String(),
		SignedDataSuffix: bytes.HexBytes(c.SignedDataSuffix).String(),
	}
}
func convertSigtoHex(s []TmSig) []TmSigHex {
	var ret []TmSigHex
	for _, sig := range s {
		ret = append(ret, TmSigHex{
			R:                common.BytesToHash(sig.R).Hex(),
			S:                common.BytesToHash(sig.S).Hex(),
			V:                sig.V,
			EncodedTimestamp: bytes.HexBytes(sig.EncodedTimestamp).String(),
		})
	}
	return ret
}

// Proof returns a proof from provided request ID and block height
func (s bridgeServer) Proof(ctx context.Context, req *QueryMultistoreRequest) (*QueryProofResponse, error) {
	resp, err := s.MultistoreTree(ctx, req)
	if err != nil {
		return nil, err
	}
	sigRequest := QueryTmRequest{
		Height: req.Height,
	}
	n, err := s.clientCtx.GetNode()
	if err != nil {
		return nil, err
	}
	block, err := n.Block(ctx, &req.Height)
	if err != nil {
		return nil, err
	}
	sigs, err := s.TmSig(ctx, &sigRequest)
	if err != nil {
		return nil, err
	}

	headerRequest := QueryBlockheaderMerkleRequest{
		Height: req.Height,
	}
	commonVote := append(sigs.Common.SignedDataPrefix, block.BlockID.Hash...)
	commonVote = append(commonVote, sigs.Common.SignedDataSuffix...)
	header, err := s.BlockheaderMerkle(ctx, &headerRequest)
	blockRelay := BlockRelayProof{
		MultistoreProof:        convertMultistoretoHex(resp.MultiStoreTree),
		BlockHeaderMerkleParts: convertHeadertoHex(header.BlockheaderMerkle),
		CommonEncodedVotePart:  convertCommontoHex(sigs.Common),
		Signatures:             convertSigtoHex(sigs.TmSig),
		AppHash:                common.BytesToHash(block.Block.AppHash).Hex(),
		HeaderHash:             common.BytesToHash(block.BlockID.Hash).Hex(),
		CommonVote:             bytes.HexBytes(commonVote).String(),
	}

	err = json.Unmarshal(relayFormat, &relayArguments)
	if err != nil {
		panic(err)
	}
	blockRelayBytes, err := relayArguments.Pack(
		convertStore(resp.MultiStoreTree),
		convertHeader(header.BlockheaderMerkle),
		sigs.Common,
		convertSigs(sigs.TmSig),
	)
	if err != nil {
		panic(err)
	}
	// assemble key to report
	qid, err := hex.DecodeString(req.Qid)
	tbytes := keeper.Uint64ToBytes(req.Timestamp)
	key := append(types.KeyPrefix(types.ReportKey), append(qid, tbytes...)...)

	reportData := ReportDataProof{
		DataKey:     bytes.HexBytes(key).String(),
		Report:      resp.Report,
		Version:     resp.Version,
		MerklePaths: convertIavlToString(resp.Iavl),
	}
	err = json.Unmarshal(verifyFormat, &verifyArguments)
	if err != nil {
		panic(err)
	}
	reportDataBytes, err := verifyArguments.Pack(
		big.NewInt(req.Height),
		resp.Report,
		big.NewInt(int64(resp.Version)),
		convertIavl(resp.Iavl),
	)
	if err != nil {
		panic(err)
	}

	var relayAndVerifyArguments abi.Arguments
	format := `[{"type":"bytes"},{"type":"bytes"}]`
	err = json.Unmarshal([]byte(format), &relayAndVerifyArguments)
	if err != nil {
		panic(err)
	}

	evmProofBytes, err := relayAndVerifyArguments.Pack(blockRelayBytes, reportDataBytes)
	if err != nil {
		return nil, err
	}
	return &QueryProofResponse{
		Height: req.Height,
		Result: &Proof{
			BlockHeight:     uint64(req.Height - 1),
			ReportDataProof: &reportData,
			BlockRelayProof: &blockRelay,
		},
		EvmProofBytes: "0x" + hex.EncodeToString(evmProofBytes),
	}, nil
}
