package bridge

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MultistoreEvm struct {
	LuqchainIAVLStateHash            common.Hash
	MintStoreMerkleHash              common.Hash
	IcacontrollerToIcahostMerkleHash common.Hash
	FeegrantToIbcMerkleHash          common.Hash
	AccToEvidenceMerkleHash          common.Hash
	ParamsToVestingMerkleHash        common.Hash
}

type BlockHeaderEvm struct {
	VersionAndChainIdHash            common.Hash
	Height                           uint64
	TimeSecond                       uint64
	TimeNanoSecondFraction           uint32
	LastBlockIdCommitMerkleHash      common.Hash
	NextValidatorConsensusMerkleHash common.Hash
	LastResultsHash                  common.Hash
	EvidenceProposerMerkleHash       common.Hash
}

type TmSigEvm struct {
	R                common.Hash
	S                common.Hash
	V                uint8
	EncodedTimestamp []byte
}

type IavlEvm struct {
	IsDataOnRight  bool
	SubtreeHeight  uint8
	SubtreeSize    *big.Int
	SubtreeVersion *big.Int
	SiblingHash    common.Hash
}
