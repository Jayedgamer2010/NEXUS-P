import { useState, useEffect, useCallback } from 'react';

interface UseApiState<T> {
  data: T | null;
  loading: boolean;
  error: string | null;
}

/**
 * Generic data fetching hook with loading/error/data states.
 *
 * Usage:
 *   const { data, loading, error, refetch } = useApi(
 *     () => adminApi.getUsers(page, 20),
 *     [page]
 *   );
 */
function useApi<T>(
  fetchFn: () => Promise<T>,
  deps: unknown[] = [],
  autoFetch = true,
): UseApiState<T> & { refetch: () => void } {
  const [state, setState] = useState<UseApiState<T>>({
    data: null,
    loading: autoFetch,
    error: null,
  });

  const fetch = useCallback(async () => {
    try {
      setState((prev) => ({ ...prev, loading: true, error: null }));
      const result = await fetchFn();
      setState({ data: result, loading: false, error: null });
    } catch (err: unknown) {
      const message =
        err instanceof Error ? err.message : 'An unexpected error occurred';
      setState({ data: null, loading: false, error: message });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps);

  useEffect(() => {
    if (autoFetch) fetch();
  }, [fetch, autoFetch]);

  return { ...state, refetch: fetch };
}

export default useApi;
