import gql from 'graphql-tag';

export const PROPOSAL_QUERY_ALL = gql`
  query Proposals(
    $first: Int
    $skip: Int
    $orderBy: String
    $orderDirection: OrderDirection
    $id: String
  ){
    proposals(
      first: $first
      skip: $skip
      orderBy:$orderBy
      orderDirection:$orderDirection
      where: { 
        status_lt: 2
      }
    ) {
      id
      cid
      proposalId
      creator
      status
      expTime
      voteResults {
        id,
        optionId,
        votes,
      },
      voteListCid
    }
  }
`;
export const PROPOSAL_QUERY = gql`
  query Proposals(
    $first: Int
    $skip: Int
    $orderBy: String
    $orderDirection: OrderDirection
    $id: String
    $status: Int
  ){
    proposals(
      first: $first
      skip: $skip
      orderBy:$orderBy
      orderDirection:$orderDirection
      where: { 
        status_lt: 2
        status: $status
      }
    ) {
      id
      cid
      proposalId
      creator
      status
      expTime
      voteResults {
        id,
        optionId,
        votes,
      },
      voteListCid
    }
  }
`;
export const VOTE_QUERY = gql`
  query Proposal(
    $id: String!
  ){
    proposal(id: $id) {
      id
      cid
      proposalId
      creator
      status
      voteResults {
        id,
        optionId,
        votes,
      },
      voteListCid
    }
  }
`;