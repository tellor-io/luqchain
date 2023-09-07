package bridge

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	fmt "fmt"
	"sort"
	"strings"

	"github.com/cometbft/cometbft/libs/protoio"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	cometbft "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"
)

func (s bridgeServer) BlockValidatorInfo(aContext context.Context, req *QueryBlockValidatorInfoRequest) (*QueryBlockValidatorInfoResponse, error) {
	commit, err := s.getCommit(req.Height)
	if err != nil {
		return nil, err
	}

	cosmosAddrs, ethAddrs, err := getEthAddresses(&commit.SignedHeader)
	if err != nil {
		return nil, err
	}

	votingPowers, err := getValidatorVotingPowers(aContext, &s, req.Height)
	if err != nil {
		return nil, err
	}

	validators := make([]ValidatorInfo, len(cosmosAddrs))
	for i, cosmosAddr := range cosmosAddrs {
		validators[i] = ValidatorInfo{
			CosmosAddress: cosmosAddr,
			EthAddress:    ethAddrs[i],
			VotingPower:   votingPowers[cosmosAddr],
		}
	}

	return &QueryBlockValidatorInfoResponse{
		Validators: validators,
	}, nil
}

func getAddressesSignaturesAndPrefix(info *cometbft.SignedHeader) ([]string, []TmSig, CommonEncodedVotePart, error) {
	addrs := []string{}
	mapAddrs := map[string]struct {
		sig     TmSig
		valAddr string
	}{}

	prefix, err := GetPrefix(tmproto.SignedMsgType(info.Commit.Type()), info.Commit.Height, int64(info.Commit.Round))
	if err != nil {
		return nil, nil, CommonEncodedVotePart{}, err
	}

	prefix = append(prefix, []byte{34, 72, 10, 32}...)

	suffix, err := protoio.MarshalDelimited(
		&tmproto.CanonicalPartSetHeader{
			Total: info.Commit.BlockID.PartSetHeader.Total,
			Hash:  info.Commit.BlockID.PartSetHeader.Hash,
		},
	)
	if err != nil {
		return nil, nil, CommonEncodedVotePart{}, err
	}

	suffix = append([]byte{18}, suffix...)

	commonVote := CommonEncodedVotePart{SignedDataPrefix: prefix, SignedDataSuffix: suffix}

	commonPart := append(commonVote.SignedDataPrefix, info.Commit.BlockID.Hash...)
	commonPart = append(commonPart, commonVote.SignedDataSuffix...)

	chainIDBytes := []byte(info.ChainID)
	encodedChainIDConstant := append([]byte{50, uint8(len(chainIDBytes))}, chainIDBytes...)

	for _, vote := range info.Commit.Signatures {
		if !vote.ForBlock() {
			continue
		}

		encodedTimestamp := encodeTime(vote.Timestamp)

		msg := append(commonPart, []byte{42, uint8(len(encodedTimestamp))}...)
		msg = append(msg, encodedTimestamp...)
		msg = append(msg, encodedChainIDConstant...)
		msg = append([]byte{uint8(len(msg))}, msg...)

		addr, v, err := recoverETHAddress(msg, vote.Signature, vote.ValidatorAddress)

		if err != nil {
			return nil, nil, CommonEncodedVotePart{}, err
		}

		addrs = append(addrs, string(addr))
		mapAddrs[string(addr)] = struct {
			sig     TmSig
			valAddr string
		}{
			sig: TmSig{
				common.BytesToHash(vote.Signature[:32]).Hex(),
				common.BytesToHash(vote.Signature[32:]).Hex(),
				uint32(v),
				common.BytesToHash(encodedTimestamp).Hex(),
			},
			valAddr: string(vote.ValidatorAddress), // Storing validator address in the map
		}
	}

	if len(addrs) == 0 {
		return nil, nil, CommonEncodedVotePart{}, fmt.Errorf("No valid precommit")
	}

	signatures := make([]TmSig, len(addrs))
	cosmosAddresses := make([]string, len(addrs)) // Allocate with exact size for efficiency
	sort.Strings(addrs)
	for i, addr := range addrs {
		signatures[i] = mapAddrs[addr].sig
		cosmosAddresses[i] = mapAddrs[addr].valAddr // Extract validator addresses in the same order
	}

	return cosmosAddresses, signatures, commonVote, nil
}

func getEthAddresses(info *cometbft.SignedHeader) ([]string, []string, error) {
	// Get signatures and common encoded vote parts from the signed header
	cosmosAddresses, signatures, commonEncodedPart, err := getAddressesSignaturesAndPrefix(info)
	if err != nil {
		return nil, nil, err
	}

	chainIDBytes := []byte(info.ChainID)
	encodedChainIDConstant := append([]byte{50, uint8(len(chainIDBytes))}, chainIDBytes...)

	var addresses []string
	for _, sig := range signatures {
		address, err := CheckTimeAndRecoverSigner(sig, commonEncodedPart, encodedChainIDConstant)
		if err != nil {
			return nil, nil, err
		}
		addresses = append(addresses, address)
	}

	return cosmosAddresses, addresses, nil
}

func CheckTimeAndRecoverSigner(sig TmSig, commonEncodedPart CommonEncodedVotePart, encodedChainID []byte) (string, error) {
	// Ensure valid timestamp size
	encodedTimestampLen := len(sig.EncodedTimestamp)
	if encodedTimestampLen < 6 || encodedTimestampLen > 12 {
		return "", errors.New("Invalid timestamp's size")
	}

	// Construct the encoded canonical vote
	encodedCanonicalVote := append(commonEncodedPart.SignedDataPrefix, 42)
	encodedCanonicalVote = append(encodedCanonicalVote, byte(encodedTimestampLen))
	encodedCanonicalVote = append(encodedCanonicalVote, sig.EncodedTimestamp...)
	encodedCanonicalVote = append(encodedCanonicalVote, commonEncodedPart.SignedDataSuffix...)
	encodedCanonicalVote = append(encodedCanonicalVote, encodedChainID...)

	// Construct the data to hash
	dataToHash := append([]byte{byte(len(encodedCanonicalVote))}, encodedCanonicalVote...)
	hashedData := sha256.Sum256(dataToHash)

	// Perform ecrecover
	publicKeyECDSA, err := ecrecover(hashedData[:], sig.V, sig.R, sig.S)
	if err != nil {
		return "", err
	}

	address := publicKeyBytesToAddress(publicKeyECDSA)
	return strings.ToLower(address), nil
}

func ecrecover(hash []byte, v uint32, rStr, sStr string) (*ecdsa.PublicKey, error) {
	signature := make([]byte, 65)
	copy(signature[32-len(rStr):32], rStr)
	copy(signature[64-len(sStr):64], sStr)
	signature[64] = byte(v - 27)
	return crypto.SigToPub(hash, signature)
}

func publicKeyBytesToAddress(pub *ecdsa.PublicKey) string {
	pubBytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(pubBytes[1:]) // omit the prefix byte 0x04 of uncompressed public key
	if err != nil {
		panic(err)
	}
	return "0x" + hex.EncodeToString(hash.Sum(nil)[12:])
}

func getValidatorVotingPowers(goContext context.Context, b *bridgeServer, height int64) (map[string]int64, error) {
	ctx := sdk.UnwrapSDKContext(goContext)

	// Call the Validators method on the TendermintRPC client
	result, err := b.clientCtx.Client.Validators(ctx, &height, nil, nil)
	if err != nil {
		return nil, err
	}

	// Create a map to store voting powers with the validator address as the key
	votingPowers := make(map[string]int64)

	for _, validator := range result.Validators {
		votingPowers[string(validator.Address)] = validator.VotingPower
	}

	return votingPowers, nil
}
