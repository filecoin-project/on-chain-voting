import {
  AdminChanged as AdminChangedEvent,
  BeaconUpgraded as BeaconUpgradedEvent,
  Initialized as InitializedEvent,
  OwnershipTransferred as OwnershipTransferredEvent,
  ProposalCancel as ProposalCancelEvent,
  ProposalCount as ProposalCountEvent,
  ProposalCreate as ProposalCreateEvent,
  Upgraded as UpgradedEvent,
  Vote as VoteEvent
} from "../generated/PowerVoting/PowerVoting"
import {
  AdminChanged,
  BeaconUpgraded,
  Initialized,
  OwnershipTransferred,
  Proposal,
  Upgraded,
  Vote,
  VoteResult
} from "../generated/schema"

export function handleAdminChanged(event: AdminChangedEvent): void {
  let entity = new AdminChanged(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousAdmin = event.params.previousAdmin
  entity.newAdmin = event.params.newAdmin

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleBeaconUpgraded(event: BeaconUpgradedEvent): void {
  let entity = new BeaconUpgraded(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.beacon = event.params.beacon

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleInitialized(event: InitializedEvent): void {
  let entity = new Initialized(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.version = event.params.version

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleOwnershipTransferred(
  event: OwnershipTransferredEvent
): void {
  let entity = new OwnershipTransferred(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.previousOwner = event.params.previousOwner
  entity.newOwner = event.params.newOwner

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleProposalCancel(event: ProposalCancelEvent): void {
  let entity=Proposal.load( event.params.id.toString())
  if (entity) {

    entity.status=event.params.status
    entity.save()
  }

}

export function handleProposalCount(event: ProposalCountEvent): void {
  let entity=Proposal.load( event.params.id.toString())
  if (entity) {
    entity.status = event.params.status
    for (let i = 0; i < event.params.voteResult.length; i++) {
      let voteResult = new VoteResult(`${event.params.id.toString()}/${i}`)
      voteResult.optionId =  event.params.voteResult[i].optionId
      voteResult.votes  =event.params.voteResult[i].votes
      voteResult.proposalId = event.params.id.toI32()
      voteResult.proposal= event.params.id.toString()
      voteResult.save()
    }
    entity.status = event.params.status
    entity.voteListCid=event.params.voteListCid
    entity.save()
  }
}

export function handleProposalCreate(event: ProposalCreateEvent): void {
  let entity = new Proposal(
    event.params.id.toString()
  )
  entity.proposalId =  event.params.id.toI32()
  entity.status = event.params.status
  entity.cid = event.params.proposal.cid
  entity.creator = event.params.proposal.creator
  entity.expTime = event.params.proposal.expTime
  entity.chainId = event.params.proposal.chainId
  entity.proposalType=event.params.proposal.proposalType

  entity.save()
}

export function handleUpgraded(event: UpgradedEvent): void {
  let entity = new Upgraded(
    event.transaction.hash.concatI32(event.logIndex.toI32())
  )
  entity.implementation = event.params.implementation

  entity.blockNumber = event.block.number
  entity.blockTimestamp = event.block.timestamp
  entity.transactionHash = event.transaction.hash

  entity.save()
}

export function handleVote(event: VoteEvent): void {
  let vote=Vote.load(`${event.params.voter.toHexString()}/${event.params.id.toI32()}`)
  if (vote) {
    vote.proposalId = event.params.id.toI32()
    vote.voteInfo = event.params.voteInfo
    vote.transactionHash = event.transaction.hash
    vote.save()
  }else {
    let entity = new Vote(`${event.params.voter.toHexString()}/${event.params.id.toI32()}`)
    entity.proposalId = event.params.id.toI32()
    entity.voteInfo = event.params.voteInfo
    entity.transactionHash = event.transaction.hash
    entity.save()
  }
}
