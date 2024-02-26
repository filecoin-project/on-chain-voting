# Compilation of the Oracle

## 1. Obtain the code for the Oracle Contract, with the repository branch set to: filecoin.

```python
git clone https://gitlab.com/storswiftlabs/wh/dapp/power-voting/kyc-oracle.git
```

## 2. Switch branch and enter the contract directory.

```python
git checkout filecoin
cd contract
```

## 3. Copy the code to[Remix](https://remix.ethereum.org/).

![Untitled](img/1.png)

## 4. Open the Oracle.sol file and compile it.

![Untitled](img/2.png)

## 5. Connect to MetaMask and switch to the Filecoin network.

## 6. Check 'Deploy with Proxy' and then click the 'Deploy' button.

![Untitled](img/3.png)

## 7. After deployment, there are two contracts: ORACLE is the logic contract, and ERC1967PROXY is the proxy contract.

![Untitled](img/4.png)

## 8. Open the ERC1967PROXY contract, call the updatePowerVotingContract method, and pass the deployed PowerVoting contract as an argument.

![Untitled](img/5.png)

## 9. Call the updateNodeAllowList method and pass the address used by the Oracle Node and true to update the Oracle Node whitelist.

![Untitled](img/6.png)