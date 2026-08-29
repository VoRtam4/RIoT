/**
 * @file usePageState.ts
 * @brief Pomocný hook pro ukládání a obnovu stavu stránek při navigaci.
 *
 * @author Vojtěch Hubáček
 *
 * @par Autorský podíl
 * - Vojtěch Hubáček: návrh a implementace celé funkcionality souboru.
 *
 * @ingroup riot_frontend
 */
import { useCallback, useMemo } from "react";
import { useLocation, useNavigate, useSearchParams } from "react-router-dom";

export type PageStateCodec<Q, E> = {
  decodeQuery: (searchParams: URLSearchParams) => Q;
  encodeQuery: (query: Q) => URLSearchParams;
  decodeEntry: (state: unknown) => E;
  encodeEntry: (entry: E) => unknown;
};

type SetPageStateOptions = {
  replace?: boolean;
};

export function usePageState<Q, E>(codec: PageStateCodec<Q, E>) {
  const [searchParams] = useSearchParams();
  const location = useLocation();
  const navigate = useNavigate();

  const query = useMemo(
    () => codec.decodeQuery(searchParams),
    [codec, searchParams],
  );

  const entry = useMemo(
    () => codec.decodeEntry(location.state),
    [codec, location.state],
  );

  const setPageState = useCallback(
    (
      next: {
        query: Q;
        entry: E;
      },
      options?: SetPageStateOptions,
    ) => {
      const nextSearchParams = codec.encodeQuery(next.query);
      const nextSearch = nextSearchParams.toString();

      navigate(
        {
          pathname: location.pathname,
          search: nextSearch ? `?${nextSearch}` : "",
        },
        {
          replace: options?.replace ?? true,
          state: codec.encodeEntry(next.entry),
        },
      );
    },
    [codec, location.pathname, navigate],
  );

  return {
    query,
    entry,
    setPageState,
  };
}
