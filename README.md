# Power Voting

## ${\color{black}{1 \ Overview}}$

PowerVoting is going to enable people voting in trustless based on FEVM, and it uses the timelock feature based on drand which can make sure the data security during the voting process. Timelock feature gives you time based encryption and decryption capabilities by relying on a drand threshold network. please refer to https://drand.love/ for drand network details, which is verifiable, unpredictable and unbiased random numbers as a service. Per encryption/decryption time based, the secret info can’t be leaked before end time, thus it doesn’t need a centralized public notary to monitor the process and justify the result.

## ${\color{black}{2 \ Problem}}$

In the centralized voting process, since the vote counting power is in the hands of the centralized organization, it will cause problems such as vote fraud and black box operation of vote counting, resulting in the voting results being manipulated by others, which cannot truly reflect the wishes of the community. 

In the community voting process governed by DAO, since the voting data of other community members can be seen before the vote counting time, the community members will be affected by the existing voting data before voting, and some members will even take advantage of a large number of voting rights in their hands to vote at the end of the voting process to make the voting results are reversed, resulting in unfair voting.

## ${\color{black}{3 \ Solution}}$
Power Voting stores voting information on the blockchain, and makes voting rights into NFTs. All voting operations are completed on the chain, which is open and transparent. When community members vote, they use the timelock feature based on drand to encrypt the voting content, and it cannot be decrypted until the counting time arrives, so that no one can know the voting information of other members before the counting time arrives. After the counting time arrives, any voting participant can initiate a vote count without being affected by any centralized organization.

## ${\color{black}{5 \ Smart \ Contract \ Functionalities}}$
* Makes voting rights into NFTs
* Stores voting information and states to the filecoin storage network
* Uses drand timelock feature to to encrypt the voting content and decrypt after time ends

## ${\color{black}{7 \ Roadmap}}$
| Time  | Status |
| ------------- | ------------- |
| 2023-01-30  | System research and requirement analysis |
| 2023-02-15  | Develop roadmap |
| 2023-02-20  |  NFT contract, voting contract and voting based on timelock feature design |
| 2023-03-03  |  NFT contract, voting contract and voting based on timelock feature development |
| 2023-03-25  |  Power Voting frontend Development |
| 2023-04-12   |  Release beta version |
