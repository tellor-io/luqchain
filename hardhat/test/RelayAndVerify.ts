import { ethers } from "hardhat";
import { expect } from "chai";
import axios from 'axios';
import { validators, encodedChainID } from "../config/contractConfig";
import { Bridge } from "../typechain-types/Bridge";

describe("Test relay and verify", function () {
    let bridge: Bridge;

    beforeEach(async function () {
        const BridgeFactory = await ethers.getContractFactory("Bridge");
        bridge = await BridgeFactory.deploy(validators, encodedChainID);
        return { bridge };
    });

    it("should call relayAndVerify with data from API", async function () {
        const retrieveAllurl = "http://localhost:1317/luqchain/luqchain/retrieve_all";
        const reports = await axios.get(retrieveAllurl);
        const height = reports.headers["grpc-metadata-x-cosmos-block-height"];
        const reportsLength = reports.data.report.length;
        // if reportsLength is 0, then there is no data to relay exit test
        if (reportsLength === 0) {
            console.log("No data to relay, exiting test");
            return;
        }
        const reportsData = reports.data.report[reportsLength - 1];
        const qid = ethers.keccak256(`0x${reportsData.qdata}`).substring(2);
        const timestamp: number = parseInt(reportsData.timestamp);

        const url = `http://localhost:1317/luqchain/bridge/proof?height=${height}&qid=${qid}&timestamp=${timestamp}`;

        // // Fetch data from the API
        const response = await axios.get(url);
        const apiData = response.data;

        let proofBytes = apiData.evmProofBytes;
        let key = apiData.result.reportDataProof.dataKey;

        const tx = await bridge.relayAndVerify(proofBytes, `0x${key}`);
        let receipt = await tx.wait();
        expect(receipt?.status).to.equal(1);
        let details = await bridge.blockDetails(height);
        expect(details.timeSecond).to.equal(parseInt(apiData.result.blockRelayProof.blockHeaderMerkleParts.timeSecond));
        expect(details.timeNanoSecondFraction).to.equal(parseInt(apiData.result.blockRelayProof.blockHeaderMerkleParts.timeNanosecond));
        expect(details.moduleState).to.equal(apiData.result.blockRelayProof.MultistoreProof.luqchain_iavl_state_hash);

        let verifyReport = await bridge.verifyOracleResult(proofBytes, `0x${key}`);
        expect(verifyReport[1]).to.equal(12350);
    });
});
