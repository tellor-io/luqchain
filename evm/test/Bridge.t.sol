// SPDX-License-Identifier: MIT
pragma solidity 0.8.17;

import "forge-std/Test.sol";
import "forge-std/console.sol";
import {Bridge} from "../src/Bridge.sol";

contract BridgeTest is Test {
    Bridge public bridge;
    bytes32 public APP_HASH =
        0xFC24423DC9259DF1A9ACF9F9775A592B32223534940B202AAD0A30EB0EF3D8A8;
    bytes32 public BLOCK_HASH =
        0x5D7FDF779203179A38EDCDB393205708BF5E504747182CAD71B954D249D6DD1C;
    bytes public COMMON =
        hex"080211640000000000000022480a205d7fdf779203179a38edcdb393205708bf5e504747182cad71b954d249d6dd1c1224080112205800736fadf214cd2e7fd0b0072c9623a9e439b593a8595ded03906300f940cf";
    Bridge.MultistoreData public multistoreData =
        // block - 1
        Bridge.MultistoreData({
            luqchainIAVLStateHash: 0xC8E79FA92F8CE06FE87FA6A6AD337DB11DDC7C258F7C04D423A87A6C75FD0AA8,
            mintStoreMerkleHash: 0xC988456482A3FD2BA21408D74FDD19921B1DC9B02A916FDE6B9F923CFCBCECE8,
            icacontrollerToIcahostMerkleHash: 0xAA650406EA0D76E39DD43D2EA6A91E3FDAA1C908FC21A7CA68E5E62CC8115639,
            feegrantToIbcMerkleHash: 0xC2E31F915AC8DB1049C0DDF8C2EF086CC4DB0EDEFE1DE5F9480C83E8A79B2AAC,
            accToEvidenceMerkleHash: 0x2787024CF1CF32B92DAEBAC63BC1846DAA99A33B5C954711EDD3E21597983941,
            paramsToVestingMerkleHash: 0x7FB3D2D5427A3D5FD89274305A4062C22AD4A6F4E112DCAA1128C4DA4DF98E3A
        });
    Bridge.BlockHeaderMerklePartsData public header =
        Bridge.BlockHeaderMerklePartsData({
            versionAndChainIdHash: 0x30a56e8396af0bf8abfe2b36fbc3b9cce5d179f94bd63ec6c906ed9eed360676,
            height: 100,
            timeSecond: 1693919981,
            timeNanoSecondFraction: 970037339,
            lastBlockIdCommitMerkleHash: 0x23f7ba9e16b08522fe473535269ff1d849a130fa310863f9ece80906f7a2258a,
            nextValidatorConsensusMerkleHash: 0x82163e6be22a258ccfa0da69d2cd2f7a3d9dd13c1dea71ea2497dca7dcf6d256,
            lastResultsHash: 0x9fb9c7533caf1d218da3af6d277f6b101c42e3c3b75d784242da663604dd53c2,
            evidenceProposerMerkleHash: 0xb8d13167c10b18857f690a0168462a558bfb22b0eb311a494f656e06222259e6
        });
    Bridge.CommonPartsData public common =
        Bridge.CommonPartsData({
            signedDataPrefix: hex"080211640000000000000022480A20",
            signedDataSuffix: hex"1224080112205800736FADF214CD2E7FD0B0072C9623A9E439B593A8595DED03906300F940CF"
        });

    function setUp() public {
        Bridge.ValidatorWithPower[]
            memory vps = new Bridge.ValidatorWithPower[](2);
        vps[0] = Bridge.ValidatorWithPower(
            0x008c1B0e9CdD79Cf896a6e54cC60353BC104A313,
            uint256(1000000)
        );
        vps[1] = Bridge.ValidatorWithPower(
            0x2461D0b9B9808F56Af2E71B4aEcC90C60Ed9AB90,
            uint256(500000)
        );

        bridge = new Bridge(vps, hex"32086C7571636861696E");
        assertEq(bridge.encodedChainID(), hex"32086C7571636861696E");
        assertEq(bridge.getNumberOfValidators(), 2);
        assertEq(
            bridge.getValidatorPower(
                0x008c1B0e9CdD79Cf896a6e54cC60353BC104A313
            ),
            1000000
        );
        assertEq(bridge.getAllValidatorPowers()[0].power, vps[0].power);
    }

    function testgetBlockHeader() public {
        bytes32 _header = bridge.getBlockHeader(header, APP_HASH);
        bytes32 actualBlockHash = BLOCK_HASH;
        assertEq(_header, actualBlockHash);
    }

    function testgetAppHash() public {
        assertEq(bridge.concatenate("luqchain"), hex"086c7571636861696e20");
        bytes32 _appHash = bridge.getAppHash(multistoreData);
        assertEq(_appHash, APP_HASH);
    }

    function testGetParentHash() public {
        bytes32 _parentHash = bridge.getParentHash(
            Bridge.IAVLData({
                isDataOnRight: true,
                subtreeHeight: 2,
                subtreeSize: 3,
                subtreeVersion: 4261,
                siblingHash: 0x4A396DA9264C7542DEAA198F318E29AA8B65F1BCE56D316C0A64D2B04D9C9795
            }),
            0x9CC4828F4C1F4542B5931969CB07CF9151E0231A88369C95F18ADEA0489539B0
        );
        assertEq(
            _parentHash,
            0x79DFEB19EAF9E40A1C365606979EBE624626DACD97C2FAA552AFF8C80DFAAF6C
        );
    }

    function testCommonencodedPart() public {
        bytes memory _parts = bridge.checkPartsAndEncodedCommonParts(
            Bridge.CommonPartsData({
                signedDataPrefix: hex"080211301100000000000022480A20",
                signedDataSuffix: hex"1224080112207E861272DEA1F9F9A4A9CA0B4C819B2C93C6B7214F5B0289F18B0BF74EAD3DE6"
            }),
            BLOCK_HASH
        );

        assertEq(
            _parts,
            // signedDataPrefix + BLOCK_HASH + signedDataSuffix
            hex"080211301100000000000022480a205d7fdf779203179a38edcdb393205708bf5e504747182cad71b954d249d6dd1c1224080112207e861272dea1f9f9a4a9ca0b4c819b2c93c6b7214f5b0289f18b0bf74ead3de6"
        );
    }

    function testcheckTimeAndRecoverSigner() public {
        address _validator = bridge.checkTimeAndRecoverSigner(
            Bridge.TMSignatureData({
                r: 0x8a58f67a8c667bce7f75de6f577353ef54c2f55d99c447c98687d87aa451d5ba,
                s: 0x14f743dd57ab32df86b22345c3acc1675383175fd82387694720ca2008c38b89,
                v: 28,
                encodedTimestamp: hex"08efd5dca70610d2baed76"
            }),
            COMMON,
            bridge.encodedChainID()
        );
        assertEq(
            abi.encodePacked(
                uint8(50),
                uint8(bytes("luqchain").length),
                bytes("luqchain")
            ),
            hex"32086c7571636861696e"
        );
        assertEq(_validator, 0x008c1B0e9CdD79Cf896a6e54cC60353BC104A313);
    }

    function testrelayBlock() public {
        Bridge.TMSignatureData[] memory sigs = new Bridge.TMSignatureData[](2);
        sigs[0] = Bridge.TMSignatureData({
            r: 0x8a58f67a8c667bce7f75de6f577353ef54c2f55d99c447c98687d87aa451d5ba,
            s: 0x14f743dd57ab32df86b22345c3acc1675383175fd82387694720ca2008c38b89,
            v: 28,
            encodedTimestamp: hex"08efd5dca70610d2baed76"
        });
        sigs[1] = Bridge.TMSignatureData({
            r: 0xdbd8114f64fd0d2f9046a6eb1e633ea0775177acbb1ddff96709f27ba03b627e,
            s: 0x51e76f99b3f57470f904af8f09d83112f3e7f022d72f3afeb30c38d01a9224af,
            v: 27,
            encodedTimestamp: hex"08efd5dca70610c1a6caa801"
        });
        bridge.relayBlock(multistoreData, header, common, sigs);
    }

    function testInclusionProof() public {
        Bridge.ValidatorWithPower[]
            memory vps = new Bridge.ValidatorWithPower[](2);
        vps[0] = Bridge.ValidatorWithPower(
            0x008c1B0e9CdD79Cf896a6e54cC60353BC104A313,
            uint256(1000000)
        );
        vps[1] = Bridge.ValidatorWithPower(
            0x2461D0b9B9808F56Af2E71B4aEcC90C60Ed9AB90,
            uint256(500000)
        );

        Bridge bridge2 = new Bridge(
            vps, // validators
            hex"32086C7571636861696E" // encoded chain ID
        );
        assertEq(bridge2.encodedChainID(), hex"32086C7571636861696E");
        assertEq(bridge2.getNumberOfValidators(), 2);
        assertEq(
            bridge2.getValidatorPower(
                0x008c1B0e9CdD79Cf896a6e54cC60353BC104A313
            ),
            1000000
        );
        assertEq(bridge2.getAllValidatorPowers()[0].power, vps[0].power);

        bytes32 rootHash = 0x75605B4972D8E1E729ED6404AE382E9A04AF62419C07F6B8684F97EF93C56B90;
        uint256 version = 4629;
        bytes memory key = hex"F03E9A3A8125B3030D3DA809A5065FB5F4FB91AE04B45C455218F4844614FC48";
        bytes32 dataHash = 0x2A2F7B2BF7A9028DE05502F0A169A72E64A6F3CD3C26F62D8C2A93618928C4A1;
        Bridge.IAVLData[] memory merklePaths = new Bridge.IAVLData[](3);
        merklePaths[0] = Bridge.IAVLData({
            isDataOnRight: true,
            subtreeHeight: 1,
            subtreeSize: 2,
            subtreeVersion: 4629,
            siblingHash: 0x289F5F11077B09013D392305B9DF9EF17C64BA4DED15F3F445A38E18F84F3F82
        });
        merklePaths[1] = Bridge.IAVLData({
            isDataOnRight: true,
            subtreeHeight: 2,
            subtreeSize: 4,
            subtreeVersion: 4629,
            siblingHash: 0xF6F2BBAC53F49488FC9C3CD8D94707F7906201C1DC7F95D15F0D88F3DF48CB13
        });
        merklePaths[2] = Bridge.IAVLData({
            isDataOnRight: true,
            subtreeHeight: 3,
            subtreeSize: 6,
            subtreeVersion: 4629,
            siblingHash: 0xE3C98274F4F977F312E2CFEB50186142B0ADBA198FFDD2D9EFC2B06A06103B14
        });
        assertEq(bridge2.verifyProof(rootHash, version, key, dataHash, merklePaths), true);
    }
}
