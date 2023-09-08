import { ethers } from "hardhat";
import { Bridge } from "../typechain-types/Bridge";
import { readFileSync } from 'fs';
import { join } from 'path';
import { Wallet } from 'ethers';

const pathToFile = join(process.env.HOME!, '.luqchain', 'config', 'priv_validator_key.json');

const data: any = JSON.parse(readFileSync(pathToFile, 'utf8'));
const privKeyBase64: string = data.priv_key.value;
const privKeyBytes: Buffer = Buffer.from(privKeyBase64, 'base64');
const privKeyHex: string = privKeyBytes.toString('hex');
const wallet = new Wallet(privKeyHex);

export const validators: Bridge.ValidatorWithPowerStruct[] = [
    { addr: wallet.address, power: ethers.toBigInt("100") }
];

export const encodedChainID = "0x32086C7571636861696E";
