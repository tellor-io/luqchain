// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.17;

import "openzeppelin/utils/structs/EnumerableMap.sol";
import "./ReportCodec.sol";
import "forge-std/console.sol";

//                                              ___________________________[AppHash]______________
//                                             /                                                  \
//                         _________________[I19]_________________                           ____[I20*]____
//                        /                                        \	                    /              \
//             _______[I15]______                          _______[I16*]______          [GHIJ]           [KLMN]
//            /                  \                        /                  \
//       __[I8*]__             __[I9]__              __[I10]__           __[I11]__
//      /         \          /         \            /         \         /         \
//    [I0]       [I1]     [I2*]        [I3]       [I4]      [I5]      [I6]       [I7]
//   /   \      /   \     /   \      /    \      /    \    /    \    /    \     /    \
// [0]   [1]  [2]   [3] [4]   [5]  [6*]    [7*] [8]  [9]  [A]   [B] [C]   [D]  [E]   [F]
// Right[7], Left[I2], Left[I8], Left[I16], Right[I20]
// [0] - acc (auth) [1] - authz    [2] - bank   [3] - capability [4] - consensus  [5] - crisis        [6] - distr
// [7] - evidence   [8] - feegrant [9] - gov    [A] - group      [B] - ibc        [C] - icacontroller [D] - icahost
// [E] - luqchain   [F] - mint     [G] - params [H] - slashing   [I] - staking    [J] - transfer      [K] - upgrade
// [L] - vesting

contract Bridge {
    using EnumerableMap for EnumerableMap.AddressToUintMap;

    struct ValidatorWithPower {
        address addr;
        uint256 power;
    }

    struct BlockDetail {
        bytes32 moduleState;
        uint64 timeSecond;
        uint32 timeNanoSecondFraction;
    }

    struct MultistoreData {
        bytes32 luqchainIAVLStateHash;
        bytes32 mintStoreMerkleHash;
        bytes32 icacontrollerToIcahostMerkleHash;
        bytes32 feegrantToIbcMerkleHash;
        bytes32 accToEvidenceMerkleHash;
        bytes32 paramsToVestingMerkleHash;
    }

    struct BlockHeaderMerklePartsData {
        bytes32 versionAndChainIdHash;
        uint64 height;
        uint64 timeSecond;
        uint32 timeNanoSecondFraction;
        bytes32 lastBlockIdCommitMerkleHash;
        bytes32 nextValidatorConsensusMerkleHash;
        bytes32 lastResultsHash;
        bytes32 evidenceProposerMerkleHash;
    }

    struct TMSignatureData {
        bytes32 r;
        bytes32 s;
        uint8 v;
        bytes encodedTimestamp;
    }

    struct IAVLData {
        bool isDataOnRight;
        uint8 subtreeHeight;
        uint256 subtreeSize;
        uint256 subtreeVersion;
        bytes32 siblingHash;
    }

    struct CommonPartsData {
        bytes signedDataPrefix;
        bytes signedDataSuffix;
    }

    /// Mapping from block height to the struct that contains block time and hash of "oracle" iAVL Merkle tree.
    mapping(uint256 => BlockDetail) public blockDetails;
    /// Mapping from an address to its voting power.
    EnumerableMap.AddressToUintMap private validatorPowers;
    /// The total voting power of active validators currently on duty.
    uint256 public totalValidatorPower;
    /// The encoded chain's ID of Band.
    bytes public encodedChainID;

    function concatenate(
        string memory module
    ) public pure returns (bytes memory) {
        bytes memory b_middle = bytes(module);
        bytes memory b_prefix = abi.encodePacked(uint8(b_middle.length)); // Convert uint8 to bytes
        bytes memory b_suffix = abi.encodePacked(uint8(32)); // Convert uint8 to bytes

        return abi.encodePacked(b_prefix, b_middle, b_suffix);
    }

    constructor(
        ValidatorWithPower[] memory validators,
        bytes memory _encodedChainID
    ) {
        for (uint256 idx = 0; idx < validators.length; ++idx) {
            ValidatorWithPower memory validator = validators[idx];
            require(
                validatorPowers.set(validator.addr, validator.power),
                "DUPLICATION_IN_INITIAL_VALIDATOR_SET"
            );
            totalValidatorPower += validator.power;
        }
        encodedChainID = _encodedChainID;
    }

    function getAppHash(
        MultistoreData memory self
    ) public pure returns (bytes32) {
        bytes32 _moduleHash = sha256(
            abi.encodePacked(self.luqchainIAVLStateHash)
        );
        bytes memory _moduleHashWithPrefix = abi.encodePacked(
            concatenate("luqchain"),
            _moduleHash
        );
        bytes32 _leafHash = merkleLeafHash(_moduleHashWithPrefix);
        bytes32 _innerHash = merkleInnerHash(
            _leafHash,
            self.mintStoreMerkleHash
        );
        _innerHash = merkleInnerHash(
            self.icacontrollerToIcahostMerkleHash,
            _innerHash
        );
        _innerHash = merkleInnerHash(self.feegrantToIbcMerkleHash, _innerHash);
        _innerHash = merkleInnerHash(self.accToEvidenceMerkleHash, _innerHash);
        _innerHash = merkleInnerHash(
            _innerHash,
            self.paramsToVestingMerkleHash
        );
        return _innerHash;
    }

    function getBlockHeader(
        BlockHeaderMerklePartsData memory self,
        bytes32 appHash
    ) public pure returns (bytes32) {
        bytes memory _encodedBlockHeight = abi.encodePacked(
            uint8(8),
            encodeVarintUnsigned(self.height)
        );
        bytes memory _encodedBlockTime = encodeTime(
            self.timeSecond,
            self.timeNanoSecondFraction
        );
        bytes memory _encodedAppHash = abi.encodePacked(
            uint8(10),
            uint8(32),
            appHash
        );
        bytes32 _blockHeightLeafHash = merkleLeafHash(_encodedBlockHeight);
        bytes32 _blockTimeLeafHash = merkleLeafHash(_encodedBlockTime);
        bytes32 _appHashLeafHash = merkleLeafHash(_encodedAppHash);
        bytes32 _leftInnerHash = merkleInnerHash(
            _blockHeightLeafHash,
            _blockTimeLeafHash
        );
        _leftInnerHash = merkleInnerHash(
            self.versionAndChainIdHash,
            _leftInnerHash
        );
        _leftInnerHash = merkleInnerHash(
            _leftInnerHash,
            self.lastBlockIdCommitMerkleHash
        );
        bytes32 _rightInnerHash = merkleInnerHash(
            _appHashLeafHash,
            self.lastResultsHash
        );
        _rightInnerHash = merkleInnerHash(
            self.nextValidatorConsensusMerkleHash,
            _rightInnerHash
        );
        _rightInnerHash = merkleInnerHash(
            _rightInnerHash,
            self.evidenceProposerMerkleHash
        );
        return merkleInnerHash(_leftInnerHash, _rightInnerHash);
    }

    /// @dev Returns the address that signed on the given encoded canonical vote message on Cosmos.
    /// @param commonEncodedPart The first common part of the encoded canonical vote.
    /// @param encodedchainID The last part of the encoded canonical vote.
    function checkTimeAndRecoverSigner(
        TMSignatureData memory self,
        bytes memory commonEncodedPart,
        bytes memory encodedchainID
    ) public pure returns (address) {
        // We need to limit the possible size of the encodedCanonicalVote to ensure only one possible block hash.
        // The size of the encodedTimestamp will be between 6 and 12 according to the following two constraints.
        // 1. The size of an encoded Unix's second is 6 bytes until over a thousand years in the future.
        // 2. The NanoSecond size can vary from 0 to 6 bytes.
        // Therefore, 6 + 0 <= the size <= 6 + 6.
        require(
            6 <= self.encodedTimestamp.length &&
                self.encodedTimestamp.length <= 12,
            "TMSignature: Invalid timestamp's size"
        );
        bytes memory encodedCanonicalVote = abi.encodePacked(
            commonEncodedPart,
            uint8(42),
            uint8(self.encodedTimestamp.length),
            self.encodedTimestamp,
            encodedchainID
        );
        return
            ecrecover(
                sha256(
                    abi.encodePacked(
                        uint8(encodedCanonicalVote.length),
                        encodedCanonicalVote
                    )
                ),
                self.v,
                self.r,
                self.s
            );
    }

    function getParentHash(
        IAVLData memory self,
        bytes32 dataSubtreeHash
    ) public pure returns (bytes32) {
        (bytes32 leftSubtree, bytes32 rightSubtree) = self.isDataOnRight
            ? (self.siblingHash, dataSubtreeHash)
            : (dataSubtreeHash, self.siblingHash);
        return
            sha256(
                abi.encodePacked(
                    self.subtreeHeight << 1, // Tendermint signed-int8 encoding requires multiplying by 2
                    encodeVarintSigned(self.subtreeSize),
                    encodeVarintSigned(self.subtreeVersion),
                    uint8(32), // Size of left subtree hash
                    leftSubtree,
                    uint8(32), // Size of right subtree hash
                    rightSubtree
                )
            );
    }

    function checkPartsAndEncodedCommonParts(
        CommonPartsData memory self,
        bytes32 blockHash
    ) public pure returns (bytes memory) {
        require(
            self.signedDataPrefix.length == 15 ||
                self.signedDataPrefix.length == 24,
            "CommonEncodedVotePart: Invalid prefix's size"
        );
        require(
            self.signedDataSuffix.length == 38,
            "CommonEncodedVotePart: Invalid suffix's size"
        );

        return
            abi.encodePacked(
                self.signedDataPrefix,
                blockHash,
                self.signedDataSuffix
            );
    }

    /// Perform checking of the block's validity by verify signatures and accumulate the voting power on Band.
    /// @param multiStore Extra multi store to compute app hash. See MultiStore lib.
    /// @param merkleParts Extra merkle parts to compute block hash. See BlockHeaderMerkleParts lib.
    /// @param commonEncodedVotePart The common part of a block that all validators agree upon.
    /// @param signatures The signatures signed on this block, sorted alphabetically by address.
    function verifyBlockHeader(
        MultistoreData memory multiStore,
        BlockHeaderMerklePartsData memory merkleParts,
        CommonPartsData memory commonEncodedVotePart,
        TMSignatureData[] memory signatures
    ) internal view returns (bytes32) {
        // Computes Tendermint's block header hash at this given block.
        bytes32 blockHeader = getBlockHeader(
            merkleParts,
            getAppHash(multiStore)
        );
        // Verify the prefix, suffix and then compute the common encoded part.
        bytes memory commonEncodedPart = checkPartsAndEncodedCommonParts(
            commonEncodedVotePart,
            blockHeader
        );
        // Create a local variable to prevent reading that state repeatedly.
        bytes memory _encodedChainID = encodedChainID;

        // Counts the total number of valid signatures signed by active validators.
        address lastSigner = address(0);
        uint256 sumVotingPower = 0;
        for (uint256 idx = 0; idx < signatures.length; ++idx) {
            address signer = checkTimeAndRecoverSigner(
                signatures[idx],
                commonEncodedPart,
                _encodedChainID
            );
            require(signer > lastSigner, "INVALID_SIGNATURE_SIGNER_ORDER");
            (bool success, uint256 power) = validatorPowers.tryGet(signer);
            if (success) {
                sumVotingPower += power;
            }
            lastSigner = signer;
        }
        // Verifies that sufficient validators signed the block and saves the oracle state.
        require(
            sumVotingPower * 3 > totalValidatorPower * 2,
            "INSUFFICIENT_VALIDATOR_SIGNATURES"
        );

        return multiStore.luqchainIAVLStateHash;
    }

    /// @param multiStore Extra multi store to compute app hash. See MultiStore lib.
    /// @param merkleParts Extra merkle parts to compute block hash. See BlockHeaderMerkleParts lib.
    /// @param commonEncodedVotePart The common part of a block that all validators agree upon.
    /// @param signatures The signatures signed on this block, sorted alphabetically by address.
    function relayBlock(
        MultistoreData memory multiStore,
        BlockHeaderMerklePartsData memory merkleParts,
        CommonPartsData memory commonEncodedVotePart,
        TMSignatureData[] memory signatures
    ) public {
        if (
            blockDetails[merkleParts.height].moduleState ==
            multiStore.luqchainIAVLStateHash &&
            blockDetails[merkleParts.height].timeSecond ==
            merkleParts.timeSecond &&
            blockDetails[merkleParts.height].timeNanoSecondFraction ==
            merkleParts.timeNanoSecondFraction
        ) return;

        blockDetails[merkleParts.height] = BlockDetail({
            moduleState: verifyBlockHeader(
                multiStore,
                merkleParts,
                commonEncodedVotePart,
                signatures
            ),
            timeSecond: merkleParts.timeSecond,
            timeNanoSecondFraction: merkleParts.timeNanoSecondFraction
        });
    }

    /// Verifies that the given data is a valid data for the given oracleStateRoot.
    /// @param moduleStateRoot The root hash of the module store.
    /// @param version Lastest block height that the data node was updated.
    /// @param report The report of this request.
    /// @param merklePaths Merkle proof that shows how the data leave is part of the module iAVL.
    function verifyResultWithRoot(
        bytes32 moduleStateRoot,
        uint256 version,
        ReportCodec.Report memory report,
        IAVLData[] memory merklePaths
    ) internal pure returns (ReportCodec.Report memory) {
        // Computes the hash of leaf node for iAVL oracle tree.
        bytes32 dataHash = sha256(ReportCodec.encode(report));

        // Verify proof
        require(
            verifyProof(
                moduleStateRoot,
                version,
                abi.encodePacked(uint8(255), report.value),
                dataHash,
                merklePaths
            ),
            "INVALID_ORACLE_DATA_PROOF"
        );

        return report;
    }

    // /// Performs oracle state relay and oracle data verification in one go. The caller submits
    // /// the encoded proof and receives back the decoded data, ready to be validated and used.
    // /// @param data The encoded data for oracle state relay and data verification.
    // function relayAndVerify(
    //     bytes calldata data
    // ) external returns (ReportCodec.Report memory) {
    //     (bytes memory relayData, bytes memory verifyData) = abi.decode(
    //         data,
    //         (bytes, bytes)
    //     );
    //     (bool relayOk, ) = address(this).call(
    //         abi.encodePacked(this.relayBlock.selector, relayData)
    //     );
    //     require(relayOk, "RELAY_BLOCK_FAILED");
    //     (bool verifyOk, bytes memory verifyResult) = address(this).staticcall(
    //         abi.encodePacked(this.verifyOracleData.selector, verifyData)
    //     );
    //     require(verifyOk, "VERIFY_ORACLE_DATA_FAILED");
    //     return abi.decode(verifyResult, (ReportCodec.Report));
    // }

    // /// Performs oracle state relay and many times of oracle data verification in one go. The caller submits
    // /// the encoded proof and receives back the decoded data, ready to be validated and used.
    // /// @param data The encoded data for oracle state relay and an array of data verification.
    // function relayAndMultiVerify(
    //     bytes calldata data
    // ) external returns (ReportCodec.Report[] memory) {
    //     (bytes memory relayData, bytes[] memory manyVerifyData) = abi.decode(
    //         data,
    //         (bytes, bytes[])
    //     );
    //     (bool relayOk, ) = address(this).call(
    //         abi.encodePacked(this.relayBlock.selector, relayData)
    //     );
    //     require(relayOk, "RELAY_BLOCK_FAILED");

    //     ReportCodec.Report[] memory results = new ReportCodec.Report[](
    //         manyVerifyData.length
    //     );
    //     for (uint256 i = 0; i < manyVerifyData.length; i++) {
    //         (bool verifyOk, bytes memory verifyResult) = address(this)
    //             .staticcall(
    //                 abi.encodePacked(
    //                     this.verifyOracleData.selector,
    //                     manyVerifyData[i]
    //                 )
    //             );
    //         require(verifyOk, "VERIFY_ORACLE_DATA_FAILED");
    //         results[i] = abi.decode(verifyResult, (ReportCodec.Report));
    //     }

    //     return results;
    // }

    // /// Performs oracle state relay and requests count verification in one go. The caller submits
    // /// the encoded proof and receives back the decoded data, ready to be validated and used.
    // /// @param data The encoded data
    // function relayAndVerifyCount(
    //     bytes calldata data
    // ) external returns (uint64, uint64) {
    //     (bytes memory relayData, bytes memory verifyData) = abi.decode(
    //         data,
    //         (bytes, bytes)
    //     );
    //     (bool relayOk, ) = address(this).call(
    //         abi.encodePacked(this.relayBlock.selector, relayData)
    //     );
    //     require(relayOk, "RELAY_BLOCK_FAILED");

    //     (bool verifyOk, bytes memory verifyResult) = address(this).staticcall(
    //         abi.encodePacked(this.verifyRequestsCount.selector, verifyData)
    //     );
    //     require(verifyOk, "VERIFY_REQUESTS_COUNT_FAILED");

    //     return abi.decode(verifyResult, (uint64, uint64));
    // }

    /// Verifies validity of the given data in the Oracle store. This function is used for both
    /// querying an oracle request and request count.
    /// @param rootHash The expected rootHash of the oracle store.
    /// @param version Lastest block height that the data node was updated.
    /// @param key The encoded key of an oracle request or request count.
    /// @param dataHash Hashed data corresponding to the provided key.
    /// @param merklePaths Merkle proof that shows how the data leave is part of the oracle iAVL.
    function verifyProof(
        bytes32 rootHash,
        uint256 version,
        bytes memory key,
        bytes32 dataHash,
        IAVLData[] memory merklePaths
    ) private pure returns (bool) {
        bytes memory encodedVersion = encodeVarintSigned(version);

        bytes32 currentMerkleHash = sha256(
            abi.encodePacked(
                uint8(0), // Height of tree (only leaf node) is 0 (signed-varint encode)
                uint8(2), // Size of subtree is 1 (signed-varint encode)
                encodedVersion,
                uint8(key.length), // Size of data key
                key,
                uint8(32), // Size of data hash
                dataHash
            )
        );

        // Goes step-by-step computing hash of parent nodes until reaching root node.
        for (uint256 idx = 0; idx < merklePaths.length; ++idx) {
            currentMerkleHash = getParentHash(
                merklePaths[idx],
                currentMerkleHash
            );
        }

        // Verifies that the computed Merkle root matches what currently exists.
        return currentMerkleHash == rootHash;
    }

    /// Get number of validators.
    function getNumberOfValidators() public view returns (uint256) {
        return validatorPowers.length();
    }

    /// Get validators by specifying an offset index and a chunk's size.
    /// @param offset An offset index of validators mapping.
    /// @param size The size of the validators chunk.
    function getValidators(
        uint256 offset,
        uint256 size
    ) public view returns (ValidatorWithPower[] memory) {
        ValidatorWithPower[] memory validatorWithPowerList;
        uint256 numberOfValidators = getNumberOfValidators();

        if (offset >= numberOfValidators) {
            // return an empty list
            return validatorWithPowerList;
        } else if (offset + size > numberOfValidators) {
            // reduce size of the entire list
            size = numberOfValidators - offset;
        }

        validatorWithPowerList = new ValidatorWithPower[](size);
        for (uint256 idx = 0; idx < size; ++idx) {
            (address addr, uint256 power) = validatorPowers.at(idx + offset);
            validatorWithPowerList[idx] = ValidatorWithPower({
                addr: addr,
                power: power
            });
        }
        return validatorWithPowerList;
    }

    /// Get all validators with power.
    function getAllValidatorPowers()
        external
        view
        returns (ValidatorWithPower[] memory)
    {
        return getValidators(0, getNumberOfValidators());
    }

    /// Get validator by address
    /// @param addr is an address of the specific validator.
    function getValidatorPower(
        address addr
    ) public view returns (uint256 power) {
        (, power) = validatorPowers.tryGet(addr);
    }

    /// Performs oracle state extraction and verification without saving root hash to storage in one go.
    /// The caller submits the encoded proof and receives back the decoded data, ready to be validated and used.
    /// @param data The encoded data for oracle state relay and data verification.
    function verifyOracleResult(
        bytes calldata data
    ) external view returns (ReportCodec.Report memory) {
        (bytes memory relayData, bytes memory verifyData) = abi.decode(
            data,
            (bytes, bytes)
        );

        (
            MultistoreData memory multiStore,
            BlockHeaderMerklePartsData memory merkleParts,
            CommonPartsData memory commonEncodedVotePart,
            TMSignatureData[] memory signatures
        ) = abi.decode(
                relayData,
                (
                    MultistoreData,
                    BlockHeaderMerklePartsData,
                    CommonPartsData,
                    TMSignatureData[]
                )
            );

        (
            ,
            ReportCodec.Report memory report,
            uint256 version,
            IAVLData[] memory merklePaths
        ) = abi.decode(
                verifyData,
                (uint256, ReportCodec.Report, uint256, IAVLData[])
            );

        return
            verifyResultWithRoot(
                verifyBlockHeader(
                    multiStore,
                    merkleParts,
                    commonEncodedVotePart,
                    signatures
                ),
                version,
                report,
                merklePaths
            );
    }

    // Utility functions
    /// @dev Returns the hash of a Merkle leaf node.
    function merkleLeafHash(
        bytes memory value
    ) internal pure returns (bytes32) {
        return sha256(abi.encodePacked(uint8(0), value));
    }

    /// @dev Returns the hash of internal node, calculated from child nodes.
    function merkleInnerHash(
        bytes32 left,
        bytes32 right
    ) internal pure returns (bytes32) {
        return sha256(abi.encodePacked(uint8(1), left, right));
    }

    function encodeVarintUnsigned(
        uint256 value
    ) internal pure returns (bytes memory) {
        // Computes the size of the encoded value.
        uint256 tempValue = value;
        uint256 size = 0;
        while (tempValue > 0) {
            ++size;
            tempValue >>= 7;
        }
        // Allocates the memory buffer and fills in the encoded value.
        bytes memory result = new bytes(size);
        tempValue = value;
        for (uint256 idx = 0; idx < size; ++idx) {
            result[idx] = bytes1(uint8(128) | uint8(tempValue & 127));
            tempValue >>= 7;
        }
        result[size - 1] &= bytes1(uint8(127)); // Drop the first bit of the last byte.
        return result;
    }

    function encodeVarintSigned(
        uint256 value
    ) internal pure returns (bytes memory) {
        return encodeVarintUnsigned(value * 2);
    }

    function encodeTime(
        uint64 second,
        uint32 nanoSecond
    ) internal pure returns (bytes memory) {
        bytes memory result = abi.encodePacked(
            hex"08",
            encodeVarintUnsigned(uint256(second))
        );
        if (nanoSecond > 0) {
            result = abi.encodePacked(
                result,
                hex"10",
                encodeVarintUnsigned(uint256(nanoSecond))
            );
        }
        return result;
    }
}
