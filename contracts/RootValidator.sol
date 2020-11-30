pragma solidity 0.5.2;

import "./MerkleUtils.sol";
import "openzeppelin-solidity/contracts/ownership/Ownable.sol";

contract RootValidator is Ownable {
    bytes32 public root;

    function setRoot(bytes32 merkleRoot) public onlyOwner {
        root = merkleRoot;
    }

    function verifyDataInState(
        bytes memory data,
        bytes32[] memory nodes,
        uint256 leafIndex
    ) public view returns (bool) {
        return MerkleUtils.containedInTree(root, data, nodes, leafIndex);
    }

    function withdraw(
        bytes memory data,
        bytes32[] memory nodes,
        uint256 leafIndex
    ) public {}
}
