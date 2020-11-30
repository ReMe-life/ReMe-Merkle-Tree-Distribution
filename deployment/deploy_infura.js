const etherlime = require('etherlime-lib');
const MerkleUtils = require('../build/MerkleUtils.json');
const RootValidator = require('../build/RootValidator.json');
const ethers = require('ethers');
const utils = ethers.utils;


const deploy = async (network, secret) => {

	const deployer = new etherlime.InfuraPrivateKeyDeployer(secret, network)
	const MerkleUtilsLib = await deployer.deploy(MerkleUtils);
	await deployer.deploy(RootValidator, {
		MerkleUtils: MerkleUtilsLib.contractAddress
	});
};

module.exports = {
	deploy
};