# StorSwift Power Voting

## 1. Overview

Power Voting dApp utilizes timelock based on smart contract technology to achieve fair and private voting. Before the voting deadline, no one’s voting results will be seen by others, and the voting process will not be disturbed by other participant’s voting results. After the voting deadline, anyone can count the votes in a decentralized manner, and the results of the counting will executed and stored by smart contract and will not be manipulated by any centralized organization or individual. 

Power Voting dApp aims to become the infrastructure of DAO governance.

## 2. Problem

In the community voting process governed by DAO, since the voting data of other community members can be seen before the vote counting time, the community members will be affected by the existing voting data before voting, and some members will even take advantage of a large number of voting rights in their hands to vote at the end of the voting process to make the voting results are reversed, resulting in unfair voting.

In the centralized voting process, since the vote counting power is in the hands of the centralized organization, it will cause problems such as vote fraud and black box operation of vote counting, resulting in the voting results being manipulated by others, which cannot truly reflect the wishes of the community.

## 3. Solution

Power Voting dApp stores voting information on the blockchain, and all voting operations are executed on the chain, which is open and transparent. 

When community members vote, they use the timelock technology to lock the voting content, and voting content cannot be viewed until the voting expiration time reaches, so that no one can know the voting information of other members before voting expiration time reaches. 

After the counting time arrives, any voting participant can initiate a vote count without being affected by any centralized organization.

## 4. Timelock

When creating a proposal, the creator will enter a voting expiration time, and Power Voting dApp will store the proposal content and voting expiration time together on the blockchain. When user votes on a proposal, Power Voting dApp will call Drand Timelock API to encrypt user's voting data and store the encrypt data into contract, the encrypt data won't be decrypt until the proposal expiration time. When proposal expiration time reached, Power Voting dApp will call Drand Timelock API to decrypt user's voting data to count the proposal. Power Voting dApp will lock all users' voting content and not allow anyone to query voting content until voting expiration time, to make sure no one can know the voting information of other members before voting expiration time reaches.

## 5. Voting Power Snapshot

Power Oracle will request raw data from FileCoin, GitHub and other data sources to calculate role identity and voting power and save them into Power Oracle contracts. SP and Client respectively invoke the `PowerAPI.minerRawPower(filActorld)` and `DataCapAPI.balance(filActorld)` interfaces to retrieve power. Power Oracle contracts will store 60 days history of voting power. When users vote, only the percentage is recorded, not the actual voting power. During the vote counting process, a random weight will be selected from the 60 days history and multiplied by the percentage to calculate the vote.

## 6. Power Voting Flowchart

![](img/flowchart.png)

## 7. Power Voting Sequence Chart

![](img/timing_diagram.png)

## 8. UCAN Design

![](img/ucan1.png)
![](img/ucan2.png)

## 9. Install

#### 1.Environment and Development Tools

1.Node.js 14 or later installed

2.npm 7 or later installed

3.Yarn v1 or v2 installed

4.Git

<img src="./img/git.png" style="zoom:50%;" alt="" />

#### 2.Download Source Code

Download the source code with the following command:

```
git clone https://github.com/black-domain/power-voting.git
```

#### 3.Install Dependencies

Install dependencies with the following command:

```
yarn
```

After yarn, you will get a 'node_modules' folder in the root directory.


<img src="img/node_modules.png" style="zoom:50%;"  alt="" />

#### 4.Update Keys in .env.example

Deploying PowerVoting and Oracle contract on Filecoin main network and replace the following address in ‘/.env.example’

<img src="img/mainnet.png" alt="" />

If you deploy the contract on Filecoin test network Calibration, you should replace the following address in ‘/.env.example’

<img src="img/testnet.png"  alt="" />

If you modify the contract code, you need to update the following abi in ‘/src/common/abi’

<img src="img/abi.png" style="zoom:50%;"  alt="" />

#### 5. Update Wallet Connect Project Id

Create wallet connect project id by https://www.rainbowkit.com/docs/migration-guide#012x-breaking-changes

Set 'WALLET_CONNECT_ID'  in ‘/.env.example’

#### 6.Build And Package

Build  with the following command:

```
yarn build
```

After building, you will get a 'dist' folder in the root directory.

<img src="img/dist.png" style="zoom:50%;"  alt="" />

#### 7.Deployment

To deploy the 'dist ' folder generated after building your front-end project, you can follow these steps:

1. **Upload the dist folder to the server**: Upload the `dist` folder to your server. You can use FTP tools, SSH, or other methods to transfer the files to a specific directory on your server.
2. **Configure the Web Server**: Ensure that your web server (such as Nginx, Apache, etc.) is properly configured, and you know where to serve static files from. Add a new site or virtual host in the configuration file and set the document root to point to the uploaded `dist` folder.
3. **Start the Web Server**: Start or restart your web server to apply the new configuration.
4. **Access the Website**: Open your browser and enter your domain name or server IP address. You should be able to see your deployed front-end application.


