/**
 * @file apollo.ts
 * @brief Konfigurace Apollo klienta pro GraphQL požadavky a subscription spojení.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import {
  ApolloClient,
  ApolloLink,
  HttpLink,
  InMemoryCache,
  Observable,
  split,
} from "@apollo/client";
import { GraphQLWsLink } from "@apollo/client/link/subscriptions";
import { createClient } from "graphql-ws";
import { getMainDefinition } from "@apollo/client/utilities";
import { apiEndpoints } from "./apiEndpoints";
import {
  isUnauthorizedRuntimeError,
  redirectToLoginForExpiredSession,
} from "../modules/auth/hooks/sessionRedirect";

const httpLink = new HttpLink({
  uri: apiEndpoints.graphqlHttp,
  credentials: "include",
});

const wsLink = new GraphQLWsLink(
  createClient({
    url: apiEndpoints.graphqlWebSocket,
    connectionParams: {},
  })
);

const splitLink = split(
  ({ query }) => {
    const definition = getMainDefinition(query);
    return (
      definition.kind === "OperationDefinition" &&
      definition.operation === "subscription"
    );
  },
  wsLink,
  httpLink
);

const unauthorizedLink = new ApolloLink((operation, forward) => {
  if (!forward) {
    return new Observable((observer) => {
      observer.complete();
    });
  }

  return new Observable((observer) => {
    const subscription = forward(operation).subscribe({
      next: (result) => {
        if (result.errors?.some((entry) => isUnauthorizedRuntimeError(entry))) {
          redirectToLoginForExpiredSession();
        }

        observer.next(result);
      },
      error: (error) => {
        if (isUnauthorizedRuntimeError(error)) {
          redirectToLoginForExpiredSession();
        }

        observer.error(error);
      },
      complete: () => {
        observer.complete();
      },
    });

    return () => subscription.unsubscribe();
  });
});

export const apolloClient = new ApolloClient({
  link: ApolloLink.from([unauthorizedLink, splitLink]),
  cache: new InMemoryCache(),
});
