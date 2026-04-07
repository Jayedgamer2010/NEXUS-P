<?php

/*
 * This file is part of MythicalDash.
 *
 * MIT License
 *
 * Copyright (c) 2020-2025 MythicalSystems
 * Copyright (c) 2020-2025 Cassian Gherman (NaysKutzu)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

namespace MythicalDash\Database;

/**
 * SupabaseClient - Low-level HTTP client for the Supabase REST API (PostgREST).
 *
 * Provides methods for CRUD operations against Supabase tables with support
 * for filtering, ordering, pagination, and count queries.
 *
 * WHERE clause translation:
 *   ['column', '=', 'value']    -> ?column=eq.value
 *   ['column', 'LIKE', '%val%'] -> ?column=like.*val*
 *   ['column', '>', 5]          -> ?column=gt.5
 *   ['column', '!=', 'val']     -> ?column=neq.val
 */
class SupabaseClient
{
    private string $baseUrl;
    private string $apiKey;

    /** Track rows affected by the last write operation for rowCount() compatibility. */
    public int $lastRowCount = 0;

    /**
     * Create a new SupabaseClient.
     *
     * @param string|null $baseUrl Override URL (reads SUPABASE_URL from env if null)
     * @param string|null $apiKey  Override key  (reads SUPABASE_SERVICE_KEY from env if null)
     */
    public function __construct(?string $baseUrl = null, ?string $apiKey = null)
    {
        $this->baseUrl = rtrim($baseUrl ?? ($_ENV['SUPABASE_URL'] ?? ''), '/');
        $this->apiKey  = $apiKey ?? ($_ENV['SUPABASE_SERVICE_KEY'] ?? ($_ENV['SUPABASE_ANON_KEY'] ?? ''));

        if (empty($this->baseUrl)) {
            throw new \Exception('SUPABASE_URL is not set in environment variables.');
        }
        if (empty($this->apiKey)) {
            throw new \Exception('SUPABASE_SERVICE_KEY (or SUPABASE_ANON_KEY) is not set in environment variables.');
        }
    }

    // ----------------------------------------------------------------
    //  High-level CRUD
    // ----------------------------------------------------------------

    /**
     * Fetch rows from a table.
     *
     * @param string      $table   Table name
     * @param array       $filters Array of [column, operator, value] (AND logic)
     * @param string      $select  Columns to select, default '*'
     * @param string|null $orderBy Column ordering, e.g. 'id.desc' or 'id'
     * @param int|null    $limit   Maximum rows
     * @param int|null    $offset  Offset
     * @return array Rows
     */
    public function getRows(
        string $table,
        array $filters = [],
        string $select = '*',
        ?string $orderBy = null,
        ?int $limit = null,
        ?int $offset = null,
    ): array {
        $url  = $this->buildUrl($table, $select, $filters, $orderBy, $limit, $offset);
        $resp = $this->request('GET', $url, null, true);
        return is_array($resp) ? $resp : [];
    }

    /**
     * Fetch a single row.
     * @return array|null
     */
    public function getRow(
        string $table,
        array $filters = [],
        string $select = '*',
    ): ?array {
        $rows = $this->getRows($table, $filters, $select, null, 1, 0);
        return !empty($rows) ? $rows[0] : null;
    }

    /**
     * Insert one or more records.
     *
     * @param string      $table      Table name
     * @param array       $data       Associative array or array of arrays (bulk)
     * @param string      $select     Columns to return, default '*'
     * @param string|null $onConflict Upsert conflict target column(s)
     * @return array Inserted record(s)
     */
    public function create(string $table, array $data, string $select = '*', ?string $onConflict = null): array
    {
        $url  = $this->buildUrl($table, $select, [], null, null, null, $onConflict);
        $body = json_encode($data);
        $resp = $this->request('POST', $url, $body, true);
        return is_array($resp) ? $resp : [];
    }

    /**
     * Update records matching filters.
     * @return int Rows updated
     */
    public function update(string $table, array $data, array $filters = []): int
    {
        $url  = $this->buildUrl($table, '*', $filters);
        $body = json_encode($data);
        $this->request('PATCH', $url, $body, false);
        return $this->lastRowCount;
    }

    /**
     * Delete records matching filters.
     * @return int Rows deleted
     */
    public function delete(string $table, array $filters = []): int
    {
        $url = $this->buildUrl($table, '*', $filters);
        $this->request('DELETE', $url, null, false);
        return $this->lastRowCount;
    }

    /**
     * Count rows matching filters.
     */
    public function count(string $table, array $filters = []): int
    {
        $url     = $this->buildUrl($table, '', $filters);  // empty select -> HEAD
        $headers = $this->defaultHeaders();
        $headers[] = 'Prefer: count=exact';

        $respHeaders = [];
        $ch = $this->initCurl('HEAD', $url, null, $headers, $respHeaders);
        curl_exec($ch);
        curl_close($ch);

        // Content-Range header format: */<total>
        if (isset($respHeaders['content-range']) && preg_match('#/(\d+)$#', $respHeaders['content-range'], $m)) {
            return (int) $m[1];
        }
        return 0;
    }

    /**
     * Check if any row matching the filters exists.
     */
    public function exists(string $table, array $filters = []): bool
    {
        return $this->count($table, $filters) > 0;
    }

    // ----------------------------------------------------------------
    //  Internal helpers
    // ----------------------------------------------------------------

    private function buildUrl(
        string $table,
        string $select,
        array $filters,
        ?string $orderBy = null,
        ?int $limit = null,
        ?int $offset = null,
        ?string $onConflict = null,
    ): string {
        $url = $this->baseUrl . '/rest/v1/' . rawurlencode($table);
        $params = [];

        if ($select !== '') {
            $params['select'] = $select;
        }

        foreach ($filters as $filter) {
            if (is_array($filter) && count($filter) >= 3) {
                $col = $filter[0];
                $op  = strtoupper($filter[1]);
                $val = $filter[2];

                $restOp = $this->operatorToRest($op, $val);
                if (array_key_exists($col, $params)) {
                    // Supabase doesn't natively allow duplicate keys for multiple filters
                    // on the same column with different values.  We append with a suffix.
                    $idx = 1;
                    while (array_key_exists($col . '_' . $idx, $params)) {
                        $idx++;
                    }
                    $params[$col . '_' . $idx] = $restOp;
                } else {
                    $params[$col] = $restOp;
                }
            } elseif (is_array($filter) && count($filter) === 2) {
                // Alternative format: ['raw_string' => value]
                $params[$filter[0]] = $filter[1];
            }
        }

        if ($orderBy !== null) {
            $parts = explode('.', $orderBy, 2);
            $col   = $parts[0];
            $dir   = isset($parts[1]) ? strtolower($parts[1]) : 'asc';
            $params['order'] = $col . '.' . $dir;
        }

        if ($limit !== null) {
            $params['limit'] = (string) $limit;
        }
        if ($offset !== null) {
            $params['offset'] = (string) $offset;
        }
        if ($onConflict !== null) {
            $params['on_conflict'] = $onConflict;
        }

        if (!empty($params)) {
            $url .= '?' . $this->buildQueryString($params);
        }

        return $url;
    }

    /**
     * Build query string preserving Supabase operators (=, like, gt, etc.).
     */
    private function buildQueryString(array $params): string
    {
        $pairs = [];
        foreach ($params as $key => $val) {
            $pairs[] = rawurlencode($key) . '=' . rawurlencode((string) $val);
        }
        return implode('&', $pairs);
    }

    /**
     * Translate a SQL operator to a PostgREST filter operator.
     * e.g. '=' -> 'eq', 'LIKE' -> 'like', '>' -> 'gt'
     */
    private function operatorToRest(string $operator, $value): string
    {
        $map = [
            '='        => 'eq',
            '=='       => 'eq',
            '!='       => 'neq',
            '<>'       => 'neq',
            '>'        => 'gt',
            '>='       => 'gte',
            '<'        => 'lt',
            '<='       => 'lte',
            'LIKE'     => 'like',
            'ILIKE'    => 'ilike',
            'NOT LIKE' => 'not.like',
            'IN'       => 'in',
            'IS'       => 'is',
            'IS NOT'   => 'not.is',
        ];

        $op = $map[strtoupper($operator)] ?? 'eq';

        // PostgREST uses * as wildcard instead of SQL %
        if (strtoupper($operator) === 'LIKE' || strtoupper($operator) === 'NOT LIKE') {
            $value = str_replace('%', '*', $value);
        }

        if (strtoupper($operator) === 'IN' && is_array($value)) {
            $value = implode(',', array_map('strval', $value));
        }

        return $op . '.' . (string) $value;
    }

    private function defaultHeaders(): array
    {
        return [
            'apikey: ' . $this->apiKey,
            'Authorization: Bearer ' . $this->apiKey,
            'Content-Type: application/json',
            'Prefer: return=representation',
        ];
    }

    /**
     * Perform HTTP request to Supabase.
     *
     * @param string      $method     HTTP verb
     * @param string      $url        Full URL
     * @param string|null $body       JSON body for POST/PATCH
     * @param bool        $returnBody Whether to decode & return the response body
     * @return mixed
     */
    private function request(string $method, string $url, ?string $body, bool $returnBody)
    {
        $headers         = $this->defaultHeaders();
        $responseHeaders = [];

        $this->lastRowCount = 0;

        $ch = $this->initCurl($method, $url, $body, $headers, $responseHeaders);
        $raw    = curl_exec($ch);
        $http   = curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $err    = curl_error($ch);
        curl_close($ch);

        if ($err) {
            $this->logError("cURL error ({$err})");
            return null;
        }

        if ($http >= 400) {
            $this->logError("HTTP {$http}: {$raw}");
            return null;
        }

        // Count returned rows from Content-Range header (e.g. "0-14/123" -> 15 rows returned, 123 total)
        if (isset($responseHeaders['content-range'])) {
            $header = $responseHeaders['content-range'];
            // Pattern:  "start-end/total"  or  "*/total" (for count)
            if (preg_match('#(\d+)-(\d+)/(\d+)#', $header, $m)) {
                $this->lastRowCount = max(1, ((int) $m[2]) - ((int) $m[1]) + 1);
            } elseif (preg_match('#\*/(\d+)#', $header, $m)) {
                $this->lastRowCount = max(0, ((int) $m[1]) - 0);
            }
        }

        if ($returnBody && $raw !== '' && $raw !== false) {
            $decoded = json_decode($raw, true);
            if (json_last_error() === JSON_ERROR_NONE) {
                return $decoded;
            }
            return $raw;
        }

        return null;
    }

    /**
     * Initialize a cURL handle.
     */
    private function initCurl(
        string $method,
        string $url,
        ?string $body,
        array $headers,
        array &$responseHeaders,
    ) {
        $ch = curl_init();

        curl_setopt_array($ch, [
            CURLOPT_URL            => $url,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_HEADER         => false,
            CURLOPT_TIMEOUT        => 30,
            CURLOPT_FOLLOWLOCATION => true,
            CURLOPT_HTTPHEADER     => $headers,
            CURLOPT_CUSTOMREQUEST  => strtoupper($method),
        ]);

        if ($body !== null) {
            curl_setopt($ch, CURLOPT_POSTFIELDS, $body);
        }

        // Capture response headers
        curl_setopt($ch, CURLOPT_HEADERFUNCTION, function ($ch, $line) use (&$responseHeaders): int {
            $len = strlen($line);
            if (($pos = strpos($line, ':')) !== false) {
                $key = strtolower(trim(substr($line, 0, $pos)));
                $val = trim(substr($line, $pos + 1));
                if (!empty($key)) {
                    $responseHeaders[$key] = $val;
                }
            }
            return $len;
        });

        return $ch;
    }

    private function logError(string $msg): void
    {
        try {
            $logger = \MythicalDash\App::getInstance(true)->getLogger();
            $logger->error('[SupabaseClient] ' . $msg);
        } catch (\Throwable $e) {
            error_log('[SupabaseClient] ' . $msg);
        }
    }

    public function logWarning(string $msg): void
    {
        try {
            $logger = \MythicalDash\App::getInstance(true)->getLogger();
            $logger->warning('[SupabaseClient] ' . $msg);
        } catch (\Throwable $e) {
            error_log('[SupabaseClient] ' . $msg);
        }
    }
}
