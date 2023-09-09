import gql from 'graphql-tag';
import {
  ApolloClient,
  createHttpLink,
  InMemoryCache
} from '@apollo/client/core';
import { SUBGRAPH_URL } from "../common/consts";

// Create the apollo client
export const apolloClient = (chainId: number) => new ApolloClient({
  link: createHttpLink({
    uri: `${SUBGRAPH_URL}-${chainId}`
  }),
  cache: new InMemoryCache({
    addTypename: false
  }),
  defaultOptions: {
    query: {
      fetchPolicy: 'no-cache'
    }
  },
  typeDefs: gql`
    enum OrderDirection {
      asc
      desc
    }
  `
});
