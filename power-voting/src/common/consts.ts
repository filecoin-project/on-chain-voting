// Copyright (C) 2023-2024 StorSwift Inc.
// This file is part of the PowerVoting library.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import { filecoin, filecoinCalibration } from 'wagmi/chains';
export const powerVotingMainNetContractAddress = process.env.POWER_VOTING_MAINNET_CONTRACT_ADDRESS || '';
export const oracleMainNetContractAddress = process.env.ORACLE_MAINNET_CONTRACT_ADDRESS || '';
export const oraclePowerMainNetContractAddress = process.env.ORACLE_POWER_MAINNET_CONTRACT_ADDRESS || '';
export const powerVotingCalibrationContractAddress = process.env.POWER_VOTING_CALIBRATION_CONTRACT_ADDRESS || '';
export const oracleCalibrationContractAddress = process.env.ORACLE_CALIBRATION_CONTRACT_ADDRESS || '';
export const oraclePowerCalibrationContractAddress = process.env.ORACLE_POWER_CALIBRATION_CONTRACT_ADDRESS || '';
export const walletConnectProjectId = process.env.WALLET_CONNECT_ID || '';
export const web3StorageEmail: any = process.env.WEB3_STORAGE_EMAIL || '';

export const contractAddressList = [
  {
    id: filecoin.id,
    address: powerVotingMainNetContractAddress
  },
  {
    id: filecoinCalibration.id,
    address: powerVotingCalibrationContractAddress
  },
];
export const githubApi = 'https://api.github.com/users';
export const proposalResultApi = '/api/proposal/result';
export const proposalHistoryApi = '/api/proposal/history';
export const worldTimeApi = 'https://worldtimeapi.org/api/timezone/Etc/UTC';
export const IN_PROGRESS_STATUS = 0;
export const COMPLETED_STATUS = 1;
export const PENDING_STATUS = 2;
export const VOTE_COUNTING_STATUS = 3;
export const VOTE_ALL_STATUS = 4;
export const WRONG_NET_STATUS = 5;
export const VOTE_OPTIONS = ['Approve', 'Reject'];
export const VOTE_LIST = [
  {
    value: PENDING_STATUS,
    color: 'bg-cyan-700',
    label: 'Pending'
  },
  {
    value: IN_PROGRESS_STATUS,
    color: 'bg-green-700',
    label: 'In Progress'
  },
  {
    value: VOTE_COUNTING_STATUS,
    color: 'bg-yellow-700',
    label: 'Vote Counting'
  },
  {
    value: COMPLETED_STATUS,
    color: 'bg-[#6D28D9]',
    label: 'Completed'
  },
]
export const VOTE_FILTER_LIST = [
  {
    label: "All",
    value: VOTE_ALL_STATUS
  },
  {
    label: "Pending",
    value: PENDING_STATUS
  },
  {
    label: "In Progress",
    value: IN_PROGRESS_STATUS
  },
  {
    label: "Vote Counting",
    value: VOTE_COUNTING_STATUS
  },
  {
    label: "Completed",
    value: COMPLETED_STATUS
  }
];
export const UCAN_TYPE_FILECOIN = 1;
export const UCAN_TYPE_GITHUB = 2;

export const UCAN_TYPE_FILECOIN_OPTIONS = [
  {
    label: 'Filecoin',
    value: UCAN_TYPE_FILECOIN
  },
];

export const UCAN_TYPE_GITHUB_OPTIONS = [
  {
    label: 'Github',
    value: UCAN_TYPE_GITHUB
  }
];

export const UCAN_GITHUB_STEP_1 = 1;
export const UCAN_GITHUB_STEP_2 = 2;
export const FILECOIN_AUTHORIZE_DOC = `
# I. How to use UCAN signature

1. First, please install Go Toolchain, you can find instructions here (https://go.dev/doc/install),Go version >=1.20. 

2. Get the code of UCAN signature tool.

\`\`\`
git clone https://gitlab.com/storswiftlabs/wh/dapp/power-voting/ucan-utils
\`\`\`

3. Go into the utils directory and install the dependencies.

\`\`\`
go mod tidy
\`\`\`

4. Build the binary file.

\`\`\`
go build -o signature
\`\`\`

5. run.

\`\`\`
./signature --aud 0x257c072306d848A6fd2f662Aead6855A7738dFEF --act add --privateKey <your_private_key> --keyType secp256k1
\`\`\`

6. Return a UCAN signature.

\`\`\`
eyJhbGciOiJzZWNwMjU2azEiLCJ0eXBlIjoiSldUIiwidmVyc2lvbiI6IjAuMC4xIn0.eyJpc3MiOiJ0MXkyNHY2Y3BiNzNwbnVkM2tlcHFoN3Zsb2h1YmNqYTR6emtrZ2MyeSIsImF1ZCI6IjB4MjU3YzA3MjMwNmQ4NDhBNmZkMmY2NjJBZWFkNjg1NUE3NzM4ZEZFRiIsImFjdCI6ImFkZCIsInByZiI6IiJ9.qYl0CQhK_EnqoKMf7Ph6x1gx1LW875y-nL__iH89s6MocYgfEZoETWAuPwwIU21LA4f-2LntzgcxdQv0Eks7bwA
\`\`\`


# II. Authorization for F1、F2 Owner、F3 addresses

## 1. Add authorization

### 1.1 Create a UCAN signature authorized by Filecoin account to Eth account

Follow the process below to create a UCAN signature with act as add.

<span style="color:red;">Attention: Field **act** should be set to **add**</span>

The parameters need to be changed as follows:

\`\`\`
var (
\taud = "0x257c072306d848A6fd2f662Aead6855A7738dFEF"  //Actual Eth address that requires authorization.
\tact = "add"  //For "act", input "add"
\tprivateKeyStr = "<your_private_key>"  //Input private key against Filecoin address. 
\tkeyTypeStr = "secp256k1"  //The encryption algorithm of Filecoin addresses is as follows: addresses starting with f1 use secp256k1, addresses starting with f3 use bls
)
\`\`\`

### 1.2 Create a UCAN signature authorized by Eth account to Filecoin account

1. Go to https://vote.storswift.io.

2. Click UCAN Delegates to  authorize.

<p>
    <img src="/images/img_1.png" />
</p>

3. Select **Filecoin** for UCAN Type. 

4. Enter  **Filecoin address** that requires authorization against field Aud. The Filecoin address is the one that its private key is entered in [1.1 Create a UCAN signature authorized by Filecoin account to Eth account](#11创建Filecoin账户对Eth账户授权的UCAN签名)

5. Enter **UCAN signature** created in 1.1 Create a UCAN signature authorized by Filecoin account to Eth account against filed Proof.

<p>
    <img src="/images/img.png" />
</p>

### 1.3 Authorization

After filling in the parameters, click **Authorize** to sign the message and send it on chain, then authorized successfully.
`;
export const FILECOIN_DEAUTHORIZE_DOC = `
# I. How to use UCAN signature

1. First, please install Go Toolchain, you can find instructions here (https://go.dev/doc/install),Go version >=1.20. 

2. Get the code of UCAN signature tool.

\`\`\`
git clone https://gitlab.com/storswiftlabs/wh/dapp/power-voting/ucan-utils
\`\`\`

3. Go into the utils directory and install the dependencies.

\`\`\`
go mod tidy
\`\`\`

4. Build the binary file.

\`\`\`
go build -o signature
\`\`\`

5. run.

\`\`\`
./signature --aud 0x257c072306d848A6fd2f662Aead6855A7738dFEF --act add --privateKey <your_private_key> --keyType secp256k1
\`\`\`

6. Return a UCAN signature.

\`\`\`
eyJhbGciOiJzZWNwMjU2azEiLCJ0eXBlIjoiSldUIiwidmVyc2lvbiI6IjAuMC4xIn0.eyJpc3MiOiJ0MXkyNHY2Y3BiNzNwbnVkM2tlcHFoN3Zsb2h1YmNqYTR6emtrZ2MyeSIsImF1ZCI6IjB4MjU3YzA3MjMwNmQ4NDhBNmZkMmY2NjJBZWFkNjg1NUE3NzM4ZEZFRiIsImFjdCI6ImFkZCIsInByZiI6IiJ9.qYl0CQhK_EnqoKMf7Ph6x1gx1LW875y-nL__iH89s6MocYgfEZoETWAuPwwIU21LA4f-2LntzgcxdQv0Eks7bwA
\`\`\`

# II. Cancel authorization

## 1 Create a UCAN signature deauthorized by Filecoin account to Eth account

<span style="color:red;">Attention：field act  should be set to del.</span>

Follow the process below to create a UCAN signature with act as del.

The parameters need to be changed as follows:

\`\`\`
var (
\taud = "0x257c072306d848A6fd2f662Aead6855A7738dFEF"  //Eth address that requires authorization
\tact = "del"  // Input "del" for field "act"
\tprivateKeyStr = "<your_private_key>"  // Input the private key against Filecoin address
\tkeyTypeStr = "secp256k1"  // The encryption algorithm of Filecoin addresses is as follows: addresses starting with f1 use secp256k1, addresses starting with f3 use bls\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t
)
\`\`\`

## 2. Create a UCAN signature deauthorized by Eth account to Filecoin account

**Prerequisite: Eth account has UCAN authorization already**。

2.1 Go to https://vote.storswift.io.

2.2 Click UCAN Delegates to  cancel authorization.  The website will monitor whether the Eth account has UCAN authorization or not. The action  will cancel the authorization if it does.

<p>
    <img src="/images/img_1.png" />
</p>

## 3. Iss & Aud are auto filled, you only need to enter UCAN created in 1.1 Create a UCAN signature deauthorized by Filecoin account to Eth account against field Proof.

<p>
    <img src="/images/img_2.png" />
</p>
`;
export const GITHUB_AUTHORIZE_DOC = `
# I. Authorization for developers

## 1. Add Authorization

### 1.1 Create a UCAN signature authorized by Eth account to Github handle

#### 1. Go to https://vote.storswift.io.

#### 2. Click UCAN Delegates to authorize.

<p>
    <img src="/images/img_2.png" />
</p>

#### 3. Select **Github** for UCAN type

#### 4. Enter **Github handle** that requires authorization in field Aud.

<p>
    <img src="/images/img_4.png" />
</p>

#### 5. Click Sign to generate. Signature is the UCAN authorized by Eth to Github. In subsequent operations, the Signature needs to be sent to the Github repo.

<p>
    <img src="/images/img_5.png" />
</p>

### 1.2 Create an initialized public repository on Github

#### Select **Public** and **Add a README file**. The repository name can be customized. There are no special requirements for that. Here UCAN is used for repo name as demonstration.

<p>
    <img src="/images/img_3.png" />
</p>

### 1.3 Create a Token used to upload UCAN signature to the repository

#### 1. Select **Developer settings** in  [Github Settings ](https://github.com/settings/profile)


<p>
    <img src="/images/img_7.png" />
</p>

#### 2. Follow 4 steps below to create Token.

<p>
    <img src="/images/img_8.png" />
</p>

#### 3. Select **write:packages**, the token name can be customized and there are no special requirements. The demonstration here uses **ucan** as the token name.

<p>
    <img src="/images/img_9.png" />
</p>

#### 4. Remember to save the Token and you will not be able to view the Token after leaving the page.

<p>
    <img src="/images/img_6.png" />
</p>

### 1.4 Upload the UCAN signature to Github repository (authorized by ETH address to Github handle) 

#### 1. Command

\`\`\`
  curl -L \\
  -X POST \\
  -H "Accept: application/vnd.Github+json" \\
  -H "Authorization: Bearer <TOKEN>" \\
  -H "X-Github-Api-Version: 2022-11-28" \\
  https://api.Github.com/repos/<OWNER>/<REPO>/git/blobs \\
  -d '{"content":"<CONTENT>","encoding":"utf-8"}'
  
\`\`\`

#### 2. Example：

​\tThe OWNER here 1.1 Create a UCAN signature authorized by Eth account to Github account is the Github handle entered in field Aud.

​\t The UCAN signature here is the one generated from 1.1 Create a UCAN signature authorized by Eth account to Github account.

​\t The REPO here is repo name created from 1.2 Create an initialized public repository on Github.

​\t The TOKEN here is one generated from 1.3 Create a Token used to upload the UCAN signature to the repository.

\`\`\`
  curl -L \\
  -X POST \\
  -H "Accept: application/vnd.Github+json" \\
  -H "Authorization: Bearer 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" \\
  -H "X-Github-Api-Version: 2022-11-28" \\
  https://api.Github.com/repos/Hzexiang/UCAN/git/blobs \\
  -d '{"content":"eyJhbGciOiJlY2RzYSIsInR5cGUiOiJKV1QiLCJ2ZXJzaW9uIjoiMC4wLjEifQ.eyJpc3MiOiIweDI1N2MwNzIzMDZkODQ4QTZmZDJmNjYyQWVhZDY4NTVBNzczOGRGRUYiLCJhdWQiOiJ0ZXN0IiwicHJmIjoiIiwiYWN0IjoiYWRkIn0.MHhmZWE5YTE5NjdjYzQ1ZDJjMmIxNTcyZDAyMzI0OGM1YWY1N2ZiNTE3ZDMxMGY3MmRhNWNiZTEyY2MxY2VjY2FjMGE1NzMwMmRmODk0ZjU1NTE2MWU4MDk3Nzc4YmFkN2M5ZDg4NzFjNmY5ODI1NmRhM2FjY2IxMGRlMzczNWY4NDFj","encoding":"utf-8"}'
\`\`\`

#### 3. Request returns the results.

\`\`\`
{
  "sha": "30662d9adc5588d55739c30299ca180e85126f54",
  "url": "https://api.Github.com/repos/<OWNER>/<REPO>/git/blobs/<FILE_SHA>"
}
\`\`\`

### 1.5 Enter the returned URL on website and proceed to the next step, then wait for the node to get the data

#### Enter the **returned URL**  as required and then click **Authorize**.

<p>
    <img src="/images/img_11.png" />
</p>
`;
export const GITHUB_DEAUTHORIZE_DOC = `
# I. Deauthorization for developers

## 1. Create a UCAN signature deauthorized by Eth account to Github handle

### 1.1 Go to https://vote.storswift.io.

### 1.2 Click UCAN Delegates to  authorize.

<p>
    <img src="/images/img_1.png" />
</p>

### 1.3 After authorized successfully for Developers, the authorized Github handle will be displayed when entering authorization page. No need to enter parameters, click on **Sign** and you will get UCAN signature for cancelling authorization. 

<p>
    <img src="/images/img_12.png" />
</p>

<p>
    <img src="/images/img_13.png" />
</p>

## 2. Upload the UCAN signature to Github repository (deauthorized by ETH address to Github handle) 

### 2.1 Create a new public repository on Github, refer to 1.2 Create an initialized public repository on Github if necessary.

### 2.2 Create a Token used to upload the UCAN signature to repository, refer to 1.3 Create a Token used to upload the UCAN signature to the repository if necessary

### 2.3 Upload  the UCAN signature to Github repository (deauthorized by ETH account to Github)  , UCAN signature is generated from 2.1 Create a UCAN signature deauthorized by Eth account to Github account 1.4 Upload the UCAN signature to Github repository (authorized by ETH address to Github handle).



## 3. Enter the returned URL on website and proceed to the next step, then wait for the node to get the data

### Enter the returned URL like below and click **Deauthorize**

<p>
    <img src="/images/img_14.png" />
</p>

## 4. Delete the Github repository that saves UCAN signature

### 4.1  After deauthorization, the Eth account can still obtain the authorized UCAN signature in the repo through the URL and authorize it again. To avoid the case mentioned before,  there is need to  delete the repository that saves authorized&deauthorized UCAN signature. 

### 4.2  Find the **settings** for repository that saves UCAN signature.

<p>
    <img src="/images/img_15.png" />
</p>

### 4.3 Select **Delete this repository** at the bottom of the page.

<p>
    <img src="/images/img_16.png" />
</p>
`;
export const DEFAULT_TIMEZONE = Intl.DateTimeFormat().resolvedOptions().timeZone;
export const web3AvatarUrl = 'https://cdn.stamp.fyi/avatar/eth';

export const UCAN_JWT_HEADER = {
  alg: 'ecdsa',
  type: 'JWT',
  version: '0.0.1'
};
export const SUCCESS_INFO= 'success';
export const ERROR_INFO= 'error';
export const OPERATION_CANCELED_MSG= 'Operation Canceled';
export const OPERATION_FAILED_MSG= 'Operation Failed';
export const STORING_DATA_MSG= 'Storing data on chain!';
export const STORING_DATA_SUCCESS_MSG= 'Data stored on chain successfully!';
export const VOTE_SUCCESS_MSG= 'Vote successful!';
export const CHOOSE_VOTE_MSG= 'Please choose a option to vote!';
export const WRONG_START_TIME_MSG= 'Start time can\'t be less than current time!';
export const WRONG_EXPIRATION_TIME_MSG= 'Expiration time can\'t be less than current time!';
export const NO_MINER_ID_MSG= 'Please input your miner ID!';
export const WRONG_MINER_ID_MSG= 'Please check your miner ID!';
export const DUPLICATED_MINER_ID_MSG= 'Your miner ID is duplicated!';
export const NOT_FIP_EDITOR_MSG= 'Please select a FIP Editor to create proposals!';