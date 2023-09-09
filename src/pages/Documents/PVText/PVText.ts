export class PVText {
    pvText(): string {
        const text = this._overview()
            + this._problem()
            + this._solution()
            + this._timelock()
            + this._pvSnapshot()
            + this._timeingDiagram()
            + this._proposalStateFlow()
            + this._proposalType()
        return text;
    }
    private _overview(): string {
        return `## 1. Overview 
&nbsp;&nbsp;&nbsp;&nbsp; Power Voting dApp utilizes timelock based on smart contract technology to achieve fair and private voting. Before the voting deadline, no one’s voting results will be seen by others, and the voting process will not be disturbed by other participant’s voting results. After the voting deadline, anyone can count the votes in a decentralized manner, and the results of the counting will executed and stored by smart contract and will not be manipulated by any centralized organization or individual. 

&nbsp;&nbsp;&nbsp;&nbsp; Power Voting dApp aims to become the infrastructure of DAO governance.
`
    }
    private _problem(): string {
        return `## 2. Problem
&nbsp;&nbsp;&nbsp;&nbsp; In the community voting process governed by DAO, since the voting data of other community members can be seen before the vote counting time, the community members will be affected by the existing voting data before voting, and some members will even take advantage of a large number of voting rights in their hands to vote at the end of the voting process to make the voting results are reversed, resulting in unfair voting.

&nbsp;&nbsp;&nbsp;&nbsp; In the centralized voting process, since the vote counting power is in the hands of the centralized organization, it will cause problems such as vote fraud and black box operation of vote counting, resulting in the voting results being manipulated by others, which cannot truly reflect the wishes of the community.
`
    }
    private _solution(): string {
        return `## 3. Solution
&nbsp;&nbsp;&nbsp;&nbsp; Power Voting dApp stores voting information on the blockchain, and all voting operations are executed on the chain, which is open and transparent. 

&nbsp;&nbsp;&nbsp;&nbsp; When community members vote, they use the timelock technology to lock the voting content, and voting content cannot be viewed until the voting expiration time reaches, so that no one can know the voting information of other members before voting expiration time reaches. 

&nbsp;&nbsp;&nbsp;&nbsp; After the counting time arrives, any voting participant can initiate a vote count without being affected by any centralized organization.       
`
    }
    private _timelock(): string {
        return `## 4. Timelock
&nbsp;&nbsp;&nbsp;&nbsp; When creating a proposal, the creator will enter a voting expiration time, and Power Voting dApp will store the proposal content and voting expiration time together on the blockchain. When user queries voting content, Power Voting dApp will check \`block.timestamp\` to see if it reaches voting expiration time. Power Voting dApp will lock all users' voting content and not allow anyone to query voting content until voting expiration time, to make sure no one can know the voting information of other members before voting expiration time reaches.
`
    }
    private _pvSnapshot(): string {
        return `## 5. Voting Power Snapshot
&nbsp;&nbsp;&nbsp;&nbsp; When creating a proposal, Power Voting dApp will get the current \`block.height\` and store it together with proposal content on the blockchain. When a user votes, Power Voting dApp will obtain the $FIL asset of the user's \`address\` corresponding to the \`block.height\` when the proposal was created at, and then use the asset amount as the voting power to vote.
`
    }
    private _timeingDiagram(): string {
        return `## 6. Timing Diagram
<div  align="center">
<img src=/images/timing_graph.png width=750 height=1600/>
</div>

`
    }
    // ![](/images/timing_graph.png)
//     <div style="width: 60% height:40%" align="center">
// <img src=/images/state_flow.png width=600 height=800/>
// </div>
    private _proposalStateFlow(): string {
        return `<h2>7. Proposal State Flow</h2>

<div align="center">
<img src=/images/state_flow.png width=600 height=800/>
</div>
`
    }
    private _proposalType(): string {
        return `<h2>8. Proposal Type</h2>

<div align="center">
<img src=/images/proposal_type.png width=600 height=800/>
</div>
`
    }

}