# I. How to use UCAN signature

1. First, please install Go Toolchain , you can find instructions here (https://go.dev/doc/install),Go version >=1.20. 

2. Get the code of UCAN signature tool.

```
git clone https://gitlab.com/storswiftlabs/wh/dapp/power-voting/ucan-utils
```

3. Install the dependencies.

```
go mod tidy
```

4. Build the binary file.

```
go build -o signature .
```

5. run.

```
./signature --aud 0x257c072306d848A6fd2f662Aead6855A7738dFEF --act add --privateKey g71qO5OWcneQqP2OL4aXfTkk0abmcGcsswVYVZzP+wo= --keyType secp256k1
```
6. Return a UCAN signature.

```
eyJhbGciOiJzZWNwMjU2azEiLCJ0eXBlIjoiSldUIiwidmVyc2lvbiI6IjAuMC4xIn0.eyJpc3MiOiJ0MXkyNHY2Y3BiNzNwbnVkM2tlcHFoN3Zsb2h1YmNqYTR6emtrZ2MyeSIsImF1ZCI6IjB4MjU3YzA3MjMwNmQ4NDhBNmZkMmY2NjJBZWFkNjg1NUE3NzM4ZEZFRiIsImFjdCI6ImFkZCIsInByZiI6IiJ9.qYl0CQhK_EnqoKMf7Ph6x1gx1LW875y-nL__iH89s6MocYgfEZoETWAuPwwIU21LA4f-2LntzgcxdQv0Eks7bwA
```



# II.  Authorization for F1、F2 Owner、F3 addresses



## 1. Add authorization



### 1.1 Create a UCAN signature authorized by Filecoin account to Eth account

1、[Follow the process below to create a UCAN signature with act as add](#i-how-to-use-ucan-signature).

<span style="color:red;">Attention: Field **act** should be set to **add**</span>.

The parameters need to be changed as follows:

```
	aud           = "0x257c072306d848A6fd2f662Aead6855A7738dFEF"  //Actual Eth address that requires authorization.
	act           = "add"	  																			//For "act", input "add"
	privateKey = "g71qO5OWcneQqP2OL4aXfTkk0abmcGcsswVYVZzP+wo="//Input private key against Filecoin address. 
	keyType    = "secp256k1"																		//The encryption algorithm of Filecoin addresses is as follows: addresses starting with f1 use secp256k1, addresses starting with f3 use bls
```



### 1.2 Create a UCAN signature authorized by Eth account to Filecoin account

1. Go to https://vote.storswift.io.

2. Click UCAN Delegates to  authorize.

![img_1.png](img/img_1.png)

3. Select **Filecoin** for UCAN Type. 

4. Enter  **Filecoin address** that requires authorization against field Aud. The Filecoin address is the one that its private key is entered in [1.1 Create a UCAN signature authorized by Filecoin account to Eth account](#11-create-a-ucan-signature-authorized-by-filecoin-account-to-eth-account).

5. Enter **UCAN signature** created in [1.1 Create a UCAN signature authorized by Filecoin account to Eth account](#11-create-a-ucan-signature-authorized-by-filecoin-account-to-eth-account)   against filed Proof.

![img.png](img/img.png)



### 1.3 Authorization

1. After filling in the parameters, click **Authorize** to sign the message and send it on chain, then authorized successfully.



## 2. Cancel authorization



### 2.1 Create a UCAN signature deauthorized by Filecoin account to Eth account

<span style="color:red;">Attention:field act  should be set to del.</span>

1、[Follow the process below to create a UCAN signature with act as del.](#i-how-to-use-ucan-signature)

The parameters need to be changed as follows:

```
	aud           = "0x257c072306d848A6fd2f662Aead6855A7738dFEF"  //Eth address that requires authorization
	act           = "del"	  																			// Input "del" for field "act"
	privateKey = "g71qO5OWcneQqP2OL4aXfTkk0abmcGcsswVYVZzP+wo="//Input the private key against Filecoin address
	keyType    = "secp256k1"																		//The encryption algorithm of Filecoin addresses is as follows: addresses starting with f1 use secp256k1, addresses starting with f3 use bls
```



### 2.2 Create a UCAN signature deauthorized by Eth account to Filecoin account

**Prerequisite: Eth account has UCAN authorization already**。

1. Go to https://vote.storswift.io.

2. Click UCAN Delegates to  cancel authorization.  The website will monitor whether the Eth account has UCAN authorization or not. The action  will cancel the authorization if it does.

![img_1.png](img/img_1.png)

3. Iss & Aud are auto filled, you only need to enter UCAN created in [2.1 Create a UCAN signature deauthorized by Filecoin account to Eth account](#21-create-a-ucan-signature-deauthorized-by-filecoin-account-to-eth-account) against field Proof.

![img_2.png](img/img_2.png)

### 2.3 Deauthorize

1. After filling in the parameters, click **Deauthorize** to sign the message and send it on chain, then deauthorized successfully.



# III. Authorization for  F4 addresses

Power Voting dApp will automatically query the corresponding F4 address and Actor address using Lotus's RPC by the Eth address.



# IV. Authorization for developers



## 1. Add Authorization



### 1.1 Create a UCAN signature authorized by Eth account to Github handle

1. Go to https://vote.storswift.io.

2. Click UCAN Delegates to  authorize.

![img_1.png](img/img_1.png)

3. Select **Github** for UCAN type.

4. Enter **Github handle** that requires authorization in field Aud.

![img_4.png](img/img_4.png)

5. Click Sign to generate. Signature is the UCAN authorized by Eth to Github. In subsequent operations, the Signature needs to be sent to the Github repo.

![img_5.png](img/img_5.png)



### 1.2 Create an initialized public repository on Github

1. Select **Public** and **Add a README file**. The repository name can be customized. There are no special requirements for that. Here UCAN is used for repo name as demonstration.

![img_3.png](img/img_3.png)



### 1.3 Create a Token used to upload UCAN signature to the repository

1. Select **Developer settings** in  [Github Settings](https://github.com/settings/profile).


![img_7.png](img/img_7.png)

2. Follow 4 steps below to create Token.

![img_8.png](img/img_8.png)

3. Select **write:packages**,the token name can be customized and there are no special requirements. The demonstration here uses **ucan** as the token name.

![img_9.png](img/img_9.png)

4. Remember to save the Token and you will not be able to view the Token after leaving the page.

![img_6.png](img/img_6.png)



### 1.4 Upload the UCAN signature to Github repository (authorized by ETH address to Github handle) 

1. Command.

```
  curl -L \
  -X POST \
  -H "Accept: application/vnd.Github+json" \
  -H "Authorization: Bearer <TOKEN>" \
  -H "X-Github-Api-Version: 2022-11-28" \
  https://api.Github.com/repos/<OWNER>/<REPO>/git/blobs \
  -d '{"content":"<CONTENT>","encoding":"utf-8"}'
  
```

2. Example:

​	2.1 The OWNER here [1.1 Create a UCAN signature authorized by Eth account to Github account](#11-create-a-ucan-signature-authorized-by-eth-account-to-github-handle) is the Github handle entered in field Aud.

​	2.2 The UCAN signature here is the one generated from [1.1 Create a UCAN signature authorized by Eth account to Github account](#11-create-a-ucan-signature-authorized-by-eth-account-to-github-handle).

​	2.3 The REPO here is repo name created from [1.2 Create an initialized public repository on Github](#12-create-an-initialized-public-repository-on-github).

​	2.4 The TOKEN here is one generated from [1.3 Create a Token used to upload the UCAN signature to the repository.](#13-create-a-token-used-to-upload-ucan-signature-to-the-repository)

```
  curl -L \
  -X POST \
  -H "Accept: application/vnd.Github+json" \
  -H "Authorization: Bearer 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" \
  -H "X-Github-Api-Version: 2022-11-28" \
  https://api.Github.com/repos/<OWNER>/ucan/git/blobs \
  -d '{"content":"eyJhbGciOiJlY2RzYSIsInR5cGUiOiJKV1QiLCJ2ZXJzaW9uIjoiMC4wLjEifQ.eyJpc3MiOiIweDI1N2MwNzIzMDZkODQ4QTZmZDJmNjYyQWVhZDY4NTVBNzczOGRGRUYiLCJhdWQiOiJ0ZXN0IiwicHJmIjoiIiwiYWN0IjoiYWRkIn0.MHhmZWE5YTE5NjdjYzQ1ZDJjMmIxNTcyZDAyMzI0OGM1YWY1N2ZiNTE3ZDMxMGY3MmRhNWNiZTEyY2MxY2VjY2FjMGE1NzMwMmRmODk0ZjU1NTE2MWU4MDk3Nzc4YmFkN2M5ZDg4NzFjNmY5ODI1NmRhM2FjY2IxMGRlMzczNWY4NDFj","encoding":"utf-8"}'
```

3. Request returns the results.

```
{
  "sha": "30662d9adc5588d55739c30299ca180e85126f54",
  "url": "https://api.Github.com/repos/<OWNER>/<REPO>/git/blobs/<FILE_SHA>"
}
```



### 1.5 Enter the returned URL on website and proceed to the next step, then wait for the node to get the data

1. Enter the **returned URL**  as required and then click **Authorize**.

![img_11.png](img/img_11.png)



## 2. Deauthorization



### 2.1 Create a UCAN signature deauthorized by Eth account to Github handle

1. Go to https://vote.storswift.io.

2. Click UCAN Delegates to  authorize.

![img_1.png](img/img_1.png)

3. After authorized successfully for Developers, the authorized Github handle will be displayed when entering authorization page. No need to enter parameters, click on **Sign** and you will get UCAN signature for cancelling authorization. 

![img_12.png](img/img_12.png)

![img_13.png](img/img_13.png)



### 2.2  Upload the UCAN signature to Github repository (deauthorized by ETH address to Github handle) 

1. Create a new public repository on Github, refer to [1.2 Create an initialized public repository on Github](#12-create-an-initialized-public-repository-on-github) if necessary.

2. Create a Token used to upload the UCAN signature to repository, refer to [1.3 Create a Token used to upload the UCAN signature to the repository](#13-create-a-token-used-to-upload-ucan-signature-to-the-repository) if necessary.

3. Upload  the UCAN signature to Github repository (deauthorized by ETH account to Github)  , UCAN signature is generated from [2.1 Create a UCAN signature deauthorized by Eth account to Github account](#21-create-a-ucan-signature-deauthorized-by-filecoin-account-to-eth-account),
   [1.4 Upload the UCAN signature to Github repository (authorized by ETH address to Github handle) ](#14-upload-the-ucan-signature-to-github-repository-authorized-by-eth-address-to-github-handle-).



### 2.3 Enter the returned URL on website and proceed to the next step, then wait for the node to get the data

1. Enter the returned URL like below and click **Deauthorize.**

![img_14.png](img/img_14.png)



### 2.4 Delete the Github repository that saves UCAN signature

1.  After deauthorization, the Eth account can still obtain the authorized UCAN signature in the repo through the URL and authorize it again. To avoid the case mentioned before,  there is need to  delete the repository that saves authorized&deauthorized UCAN signature. 

2. Find the **settings** for repository that saves UCAN signature.

![img_15.png](img/img_15.png)

3. Select **Delete this repository** at the bottom of the page.

![img_16.png](img/img_16.png)
