package bridge

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	fmt "fmt"
	. "luqchain/bridge/types"
	"luqchain/x/luqchain/keeper"
	"luqchain/x/luqchain/types"

	"github.com/cometbft/cometbft/libs/bytes"
	cometclient "github.com/cometbft/cometbft/rpc/client"
	ics23 "github.com/confio/ics23/go"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func (s bridgeServer) MultistoreTree(ctx context.Context, req *QueryMultistoreRequest) (*QueryMultistoreResponse, error) {
	var h *int64
	if req.Height != 0 {
		h = &req.Height
	}
	qid, err := hex.DecodeString(req.Qid)
	if err != nil {
		return nil, err
	}
	tbytes := keeper.Uint64ToBytes(req.Timestamp)
	resp, err := s.clientCtx.Client.ABCIQueryWithOptions(
		context.Background(),
		"/store/luqchain/key",
		append(types.KeyPrefix(types.ReportKey), append(qid, tbytes...)...),
		cometclient.ABCIQueryOptions{Height: *h - 1, Prove: true},
	)
	if err != nil {
		return nil, err
	}
	var report Report
	types.ModuleCdc.MustUnmarshal(resp.Response.GetValue(), &report)

	proof := resp.Response.GetProofOps()
	if proof == nil {
		return nil, nil
	}
	ops := proof.GetOps()
	if ops == nil {
		return nil, nil
	}

	var multistoreProof *ics23.ExistenceProof
	var iavlProof *ics23.ExistenceProof

	for _, op := range ops {
		switch op.GetType() {
		case storetypes.ProofOpIAVLCommitment:
			proof := &ics23.CommitmentProof{}
			err := proof.Unmarshal(op.Data)
			if err != nil {
				panic(err)
			}
			iavlCOps := storetypes.NewIavlCommitmentOp(op.Key, proof)
			iavlProof = iavlCOps.Proof.GetExist()
			if iavlProof == nil {
				return nil, nil
			}
		case storetypes.ProofOpSimpleMerkleCommitment:
			proof := &ics23.CommitmentProof{}
			err := proof.Unmarshal(op.Data)
			if err != nil {
				panic(err)
			}
			multiStoreOps := storetypes.NewSimpleMerkleCommitmentOp(op.Key, proof)
			multistoreProof = multiStoreOps.Proof.GetExist()
			if multistoreProof == nil {
				return nil, nil
			}
			appHash, err := multistoreProof.Calculate()
			fmt.Println("appHash", bytes.HexBytes(appHash))

		default:
			fmt.Println("Defaulting to nothing found")
			return nil, nil
		}
	}
	paths := GetMerklePaths(iavlProof)

	return &QueryMultistoreResponse{
		MultiStoreTree: MultiStoreTreeFields{
			LuqchainIavlStateHash:            multistoreProof.Value,
			MintStoreMerkleHash:              multistoreProof.Path[0].Suffix,
			IcacontrollerToIcahostMerkleHash: multistoreProof.Path[1].Prefix[1:],
			FeegrantToIbcMerkleHash:          multistoreProof.Path[2].Prefix[1:],
			AccToEvidenceMerkleHash:          multistoreProof.Path[3].Prefix[1:],
			ParamsToVestingMerkleHash:        multistoreProof.Path[4].Suffix,
		},
		Iavl:    paths,
		Version: decodeIAVLLeafPrefix(iavlProof.Leaf.Prefix),
		Report:  report,
	}, nil
}

func GetMerklePaths(iavlEp *ics23.ExistenceProof) []IAVLMerklePath {
	paths := make([]IAVLMerklePath, 0)
	for _, step := range iavlEp.Path {
		if step.Hash != ics23.HashOp_SHA256 {
			// Tendermint v0.34.9 is using SHA256 only.
			panic("Expect HashOp_SHA256")
		}
		imp := IAVLMerklePath{}

		// decode IAVL inner prefix
		// ref: https://github.com/cosmos/iavl/blob/master/proof_ics23.go#L96
		subtreeHeight, n1 := binary.Varint(step.Prefix)
		subtreeSize, n2 := binary.Varint(step.Prefix[n1:])
		subtreeVersion, n3 := binary.Varint(step.Prefix[n1+n2:])

		imp.SubtreeHeight = uint32(subtreeHeight)
		imp.SubtreeSize = uint64(subtreeSize)
		imp.SubtreeVersion = uint64(subtreeVersion)

		prefixLength := n1 + n2 + n3 + 1
		if prefixLength != len(step.Prefix) {
			imp.IsDataOnRight = true
			imp.SiblingHash = step.Prefix[prefixLength : len(step.Prefix)-1] // remove 0x20
		} else {
			imp.IsDataOnRight = false
			imp.SiblingHash = step.Suffix[1:] // remove 0x20
		}
		paths = append(paths, imp)
	}
	return paths
}

func decodeIAVLLeafPrefix(prefix []byte) uint64 {
	// ref: https://github.com/cosmos/iavl/blob/master/proof_ics23.go#L96
	_, n1 := binary.Varint(prefix)
	_, n2 := binary.Varint(prefix[n1:])
	version, _ := binary.Varint(prefix[n1+n2:])
	return uint64(version)
}
