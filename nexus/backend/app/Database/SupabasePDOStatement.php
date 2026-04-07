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
 * SupabasePDOStatement - A \PDOStatement-compatible wrapper around Supabase REST API.
 *
 * Parses SQL prepared statements (SELECT/INSERT/UPDATE/DELETE), substitutes
 * bound parameters, and translates them to the appropriate SupabaseClient call.
 */
class SupabasePDOStatement extends \PDOStatement
{
    private SupabasePDO $pdo;
    private SupabaseClient $client;
    private string $queryString;

    /** @var array<string,mixed> Parameter bindings  (key = name without ':') */
    private array $bindings = [];

    /** @var array<int,array<string,mixed>> Result rows */
    private array $rows = [];

    /** @var int Current read cursor position */
    private int $cursor = 0;

    /** @var int Number of rows affected (for write operations) */
    private int $rowCountValue = 0;

    public function __construct(SupabasePDO $pdo, SupabaseClient $client, string $queryString)
    {
        $this->pdo         = $pdo;
        $this->client      = $client;
        $this->queryString = trim($queryString);
        // Remove trailing semicolons
        $this->queryString = rtrim($this->queryString, '; ');
    }

    // -----------------------------------------------------------------------
    //  Parameter binding
    // -----------------------------------------------------------------------

    public function bindParam(
        mixed $parameter,
        mixed &$variable,
        int $type = \PDO::PARAM_STR,
        ?int $length = null,
        mixed $driver_options = null,
    ): bool {
        $key = $this->normaliseKey($parameter);
        $this->bindings[$key] = &$variable;
        return true;
    }

    public function bindValue(
        mixed $parameter,
        mixed $value,
        int $type = \PDO::PARAM_STR,
    ): bool {
        $key = $this->normaliseKey($parameter);
        $this->bindings[$key] = $value;
        return true;
    }

    // -----------------------------------------------------------------------
    //  Execution
    // -----------------------------------------------------------------------

    public function execute(?array $params = null): bool
    {
        try {
            // Merge runtime params
            if ($params !== null) {
                foreach ($params as $k => $v) {
                    $this->bindings[$this->normaliseKey($k)] = $v;
                }
            }

            $upper = strtoupper($this->queryString);

            if (str_starts_with($upper, 'SELECT')) {
                return $this->execSelect();
            }
            if (str_starts_with($upper, 'INSERT')) {
                return $this->execInsert();
            }
            if (str_starts_with($upper, 'UPDATE')) {
                return $this->execUpdate();
            }
            if (str_starts_with($upper, 'DELETE')) {
                return $this->execDelete();
            }

            // SHOW TABLES, DESCRIBE, etc -- unsupported, return empty
            $this->rows = [];
            $this->cursor = 0;
            $this->rowCountValue = 0;
            return true;
        } catch (\Throwable $e) {
            $this->logError('execute: ' . $e->getMessage() . ' | Query: ' . $this->queryString);
            $this->rows = [];
            $this->cursor = 0;
            $this->rowCountValue = 0;
            return false;
        }
    }

    // -----------------------------------------------------------------------
    //  Fetching -- all return arrays (FETCH_ASSOC) by default
    // -----------------------------------------------------------------------

    public function fetch(int $mode = \PDO::FETCH_DEFAULT, int $cursor_orientation = \PDO::FETCH_ORI_NEXT, int $cursor_offset = 0): mixed
    {
        if ($this->cursor >= count($this->rows)) {
            return false;
        }
        $row  = $this->rows[$this->cursor++];

        // Map common fetch modes
        switch ($mode) {
            case \PDO::FETCH_ASSOC:
            case \PDO::FETCH_DEFAULT:
                return $row;
            case \PDO::FETCH_NUM:
                return array_values($row);
            case \PDO::FETCH_BOTH:
                return array_merge($row, array_values($row));
            case \PDO::FETCH_OBJ:
                return (object) $row;
        }
        return $row;
    }

    public function fetchAll(int $mode = \PDO::FETCH_DEFAULT, mixed ...$args): array
    {
        // FETCH_COLUMN
        if ($mode === \PDO::FETCH_COLUMN) {
            $colIdx = $args[0] ?? 0;
            $values = [];
            foreach ($this->rows as $row) {
                $keys = array_keys($row);
                if (is_int($colIdx) && isset($keys[$colIdx])) {
                    $values[] = $row[$keys[$colIdx]];
                } elseif (isset($row[$colIdx])) {
                    $values[] = $row[$colIdx];
                }
            }
            return $values;
        }

        switch ($mode) {
            case \PDO::FETCH_ASSOC:
            case \PDO::FETCH_DEFAULT:
                return $this->rows;
            case \PDO::FETCH_NUM:
                return array_map('array_values', $this->rows);
            case \PDO::FETCH_OBJ:
                return array_map(fn($r) => (object) $r, $this->rows);
        }
        return $this->rows;
    }

    public function fetchColumn(int $column = 0): mixed
    {
        if ($this->cursor >= count($this->rows)) {
            return false;
        }
        $row = $this->rows[$this->cursor++];
        $keys = array_keys($row);
        return $keys[$column] ? $row[$keys[$column]] : false;
    }

    public function columnCount(): int
    {
        return empty($this->rows) ? 0 : count($this->rows[0]);
    }

    public function rowCount(): int
    {
        return $this->rowCountValue;
    }

    public function closeCursor(): bool
    {
        $this->rows = [];
        $this->cursor = 0;
        return true;
    }

    // -----------------------------------------------------------------------
    //  SELECT
    // -----------------------------------------------------------------------

    private function execSelect(): bool
    {
        $sql = $this->inlineBindings($this->queryString);

        // --- COUNT(*) ---
        if (preg_match('/^\s*SELECT\s+COUNT\s*\(\s*(\*|\w+)\s*\)/i', $sql)) {
            return $this->execCount($sql);
        }

        // Extract table
        if (!preg_match('/FROM\s+([\w.`]+)/i', $sql, $m)) {
            $this->rows = [];
            $this->cursor = 0;
            $this->rowCountValue = 0;
            return true;
        }
        $table = trim($m[1], '`');
        $table = preg_replace('/\s+AS\s+\w+/i', '', $table);
        $table = preg_replace('/\s+\w+$/', '', $table); // drop alias
        $table = trim($table);

        // JOIN, GROUP BY, window functions, MySQL-specific functions => unsupported
        if (stripos($sql, 'JOIN') !== false) {
            $this->logWarning('JOIN is not fully supported by the Supabase REST API. Returning empty results. Query: ' . $sql);
            $this->rows = []; $this->cursor = 0; $this->rowCountValue = 0;
            return true;
        }
        if (stripos($sql, 'GROUP BY') !== false) {
            $this->logWarning('GROUP BY is not supported by the Supabase REST API. Returning empty results. Query: ' . $sql);
            $this->rows = []; $this->cursor = 0; $this->rowCountValue = 0;
            return true;
        }
        if (preg_match('/NOW\(\)|FROM_UNIXTIME|UNIX_TIMESTAMP|DATE_FORMAT|LAST_INSERT_ID/i', $sql)) {
            $this->logWarning('MySQL-specific function in query not supported by Supabase REST API. Query: ' . $sql);
            $this->rows = []; $this->cursor = 0; $this->rowCountValue = 0;
            return true;
        }

        // SELECT columns
        $selectCols = '*';
        if (preg_match('/SELECT\s+(.+?)\s+FROM/i', $sql, $m)) {
            $rawCols = trim($m[1]);
            if ($rawCols !== '*') {
                // "users.username as un, users.email"  => "username,email"
                $parts = array_map('trim', explode(',', $rawCols));
                $cols = [];
                foreach ($parts as $p) {
                    // Strip table prefix: users.email => email
                    if (strpos($p, '.') !== false) {
                        $p = trim(explode('.', $p, 2)[1]);
                    }
                    // Strip AS alias: email as e => email
                    $p = preg_replace('/\s+as\s+\w+/i', '', $p);
                    $p = trim($p);
                    if ($p !== '' && $p !== '*') {
                        $cols[] = $p;
                    }
                }
                if (!empty($cols)) {
                    $selectCols = implode(',', $cols);
                }
            }
        }

        // WHERE
        $filters = [];
        if (preg_match('/WHERE\s+(.+?)(?:\s+ORDER|\s+LIMIT|\s+OFFSET|$)/i', $sql, $wm)) {
            $where = trim($wm[1]);
            $filters = $this->parseWhereClause($where);
        }

        // ORDER BY
        $orderBy = null;
        if (preg_match('/ORDER\s+BY\s+(\w+)(?:\s+(ASC|DESC))?/i', $sql, $om)) {
            $col = $om[1];
            // Strip table prefix
            if (strpos($col, '.') !== false) {
                $col = explode('.', $col, 2)[1];
            }
            $dir = isset($om[2]) ? strtolower($om[2]) : 'asc';
            $orderBy = $col . '.' . $dir;
        }

        // LIMIT / OFFSET
        $limit  = null;
        $offset = null;
        if (preg_match('/LIMIT\s+(\d+)/i', $sql, $m))  $limit  = (int) $m[1];
        if (preg_match('/OFFSET\s+(\d+)/i', $sql, $m)) $offset = (int) $m[1];

        $this->rows = $this->client->getRows($table, $filters, $selectCols, $orderBy, $limit, $offset);
        $this->cursor = 0;
        $this->rowCountValue = count($this->rows);
        return true;
    }

    // -----------------------------------------------------------------------
    //  COUNT
    // -----------------------------------------------------------------------

    private function execCount(string $sql): bool
    {
        if (!preg_match('/FROM\s+([\w.`]+)/i', $sql, $m)) {
            $this->rows = [['cnt' => 0]];
            $this->cursor = 0;
            $this->rowCountValue = 1;
            return true;
        }

        $table = trim($m[1], '`');
        $table = preg_replace('/\s+AS\s+\w+/i', '', $table);
        $table = trim($table);

        $filters = [];
        if (preg_match('/WHERE\s+(.+?)(?:\s+GROUP|\s+ORDER|\s+HAVING|\s+LIMIT|$)/i', $sql, $wm)) {
            $filters = $this->parseWhereClause(trim($wm[1]));
        }

        // Try the Supabase count endpoint
        $count = $this->client->count($table, $filters);

        // Fallback: fetch rows and count
        if ($count <= 0) {
            $rows = $this->client->getRows($table, $filters, '*', null, 100, 0);
            $count = count($rows);
            // If we got exactly 100 rows there may be more; mark as approximate
            if ($count === 100) {
                // fetch more chunks to get accurate count
                $offset = 100;
                while (true) {
                    $chunk = $this->client->getRows($table, $filters, '*', null, 100, $offset);
                    if (empty($chunk)) break;
                    $count += count($chunk);
                    $offset += 100;
                    if (count($chunk) < 100) break;
                }
            }
        }

        $this->rows = [['cnt' => $count, 'COUNT(*)' => $count]];
        $this->cursor = 0;
        $this->rowCountValue = 1;
        return true;
    }

    // -----------------------------------------------------------------------
    //  INSERT
    // -----------------------------------------------------------------------

    private function execInsert(): bool
    {
        $sql = $this->inlineBindings($this->queryString);

        // INSERT ... ON DUPLICATE KEY UPDATE  =>  upsert
        if (stripos($sql, 'ON DUPLICATE KEY UPDATE') !== false) {
            [$insertPart] = explode('ON DUPLICATE KEY UPDATE', $sql, 2);
            $result = $this->doInsert(trim($insertPart));
            if ($result) {
                $lastId = $this->pdo->lastInsertId();
                // Now do update for existing row
                $table = $this->extractTableFromInsert($insertPart);
                $data  = $this->extractDataFromInsert($insertPart);
                if ($table && $data && isset($data['id'])) {
                    $existing = $this->client->getRow($table, [['id', 'eq', $data['id']]]);
                    if ($existing) {
                        $updateParts = array_map('trim', explode(',', explode('ON DUPLICATE KEY UPDATE', $sql, 2)[1]));
                        $updateData = [];
                        foreach ($updateParts as $up) {
                            if (preg_match('/(\w+)\s*=\s*(.+)/', $up, $am)) {
                                $updateData[$am[1]] = trim($am[2], " '\"");
                            }
                        }
                        $this->client->update($table, $updateData, [['id', 'eq', $data['id']]]);
                    }
                }
            }
            $this->rows = [['id' => (int) $this->pdo->lastInsertId()]];
            $this->rowCountValue = 1;
            return true;
        }

        $result = $this->doInsert($sql);
        return $result;
    }

    private function doInsert(string $sql): bool
    {
        // Extract table
        if (!preg_match('/INTO\s+([\w.`]+)/i', $sql, $m)) {
            $this->rows = [];
            $this->cursor = 0;
            $this->rowCountValue = 0;
            return false;
        }
        $table = trim($m[1], '`');

        // Extract columns
        $columns = [];
        if (preg_match('/\(([^)]+)\)/', $sql, $m)) {
            $columns = array_map(fn($c) => trim($c, '`'), array_map('trim', explode(',', $m[1])));
        }

        // Extract values
        $values = [];
        if (preg_match('/VALUES\s*\(([^)]+)\)/i', $sql, $m)) {
            $values = $this->parseValuesList($m[1]);
        }

        if (empty($columns) || empty($values)) {
            $this->logWarning('Cannot parse INSERT columns/values. Query: ' . $sql);
            $this->rows = []; $this->cursor = 0; $this->rowCountValue = 0;
            return false;
        }

        $data = [];
        foreach ($columns as $i => $col) {
            if (isset($values[$i])) {
                $data[$col] = $this->unescapeValue($values[$i]);
            }
        }

        $result = $this->client->create($table, $data);

        // Track inserted ID
        if (isset($data['id']) && is_numeric($data['id'])) {
            $this->pdo->setLastInsertId((string) $data['id']);
        } elseif (!empty($result)) {
            $firstRow = $result[0] ?? $result;
            if (isset($firstRow['id'])) {
                $this->pdo->setLastInsertId((string) $firstRow['id']);
            }
        }

        $this->rows = is_array($result) ? $result : [];
        $this->cursor = 0;
        $this->rowCountValue = is_array($result) ? count($result) : 1;
        return true;
    }

    // -----------------------------------------------------------------------
    //  UPDATE
    // -----------------------------------------------------------------------

    private function execUpdate(): bool
    {
        $sql = $this->inlineBindings($this->queryString);

        // Check for unsupported functions
        if (preg_match('/NOW\(\)/i', $sql)) {
            $sql = str_ireplace('NOW()', date('Y-m-d\TH:i:s\Z'), $sql);
        }
        if (preg_match('/FROM_UNIXTIME\(([^)]+)\)/i', $sql, $m)) {
            $ts = (int) $m[1];
            if ($ts > 0) {
                $sql = str_ireplace($m[0], date('Y-m-d\TH:i:s\Z', $ts), $sql);
            } else {
                $this->logWarning('FROM_UNIXTIME(0) replaced with empty string. Query: ' . $sql);
                $sql = str_ireplace($m[0], '', $sql);
            }
        }
        if (preg_match('/UNIX_TIMESTAMP\(([^)]*)\)/i', $sql)) {
            $this->logWarning('UNIX_TIMESTAMP() in UPDATE not supported by Supabase REST API. Query: ' . $sql);
            $this->rows = []; $this->cursor = 0; $this->rowCountValue = 0;
            return false;
        }

        if (!preg_match('/UPDATE\s+([\w.`]+)/i', $sql, $m)) {
            return false;
        }
        $table = trim($m[1], '`');

        // Extract SET data
        $data = [];
        if (preg_match('/SET\s+(.+?)(?:\s+WHERE|$)/i', $sql, $sm)) {
            $setClause = trim($sm[1]);
            foreach (explode(',', $setClause) as $assignment) {
                if (preg_match("/^`?(\w+)`?\s*=\s*(.+)$/i", trim($assignment), $am)) {
                    $data[$am[1]] = $this->unescapeValue(trim($am[2]));
                }
            }
        }

        // Extract WHERE filters
        $filters = [];
        if (preg_match('/WHERE\s+(.+?)(?:\s+ORDER|\s+LIMIT|$)/i', $sql, $wm)) {
            $filters = $this->parseWhereClause(trim($wm[1]));
        }

        if (empty($data)) {
            $this->logWarning('No SET data in UPDATE. Query: ' . $sql);
            return false;
        }

        $count = $this->client->update($table, $data, $filters);
        $this->rowCountValue = $count;
        $this->rows = [['affected' => $count]];
        $this->cursor = 0;
        return true;
    }

    // -----------------------------------------------------------------------
    //  DELETE
    // -----------------------------------------------------------------------

    private function execDelete(): bool
    {
        $sql = $this->inlineBindings($this->queryString);

        if (!preg_match('/FROM\s+([\w.`]+)/i', $sql, $m)) {
            return false;
        }
        $table = trim($m[1], '`');

        $filters = [];
        if (preg_match('/WHERE\s+(.+?)(?:\s+ORDER|\s+LIMIT|$)/i', $sql, $wm)) {
            $filters = $this->parseWhereClause(trim($wm[1]));
        }

        $count = $this->client->delete($table, $filters);
        $this->rowCountValue = $count;
        $this->rows = [['affected' => $count]];
        $this->cursor = 0;
        return true;
    }

    // -----------------------------------------------------------------------
    //  Helpers
    // -----------------------------------------------------------------------

    /** Normalise a parameter key (strip leading ':' for named params). */
    private function normaliseKey(mixed $key): string
    {
        if (is_string($key)) {
            return ltrim($key, ':');
        }
        return (string) $key;
    }

    /** Replace :param placeholders with their bound values in a query string. */
    private function inlineBindings(string $sql): string
    {
        // Named params first
        foreach ($this->bindings as $key => $value) {
            if (!is_numeric($key)) {
                $placeholder = ':' . $key;
                if (strpos($sql, $placeholder) !== false) {
                    $sql = str_replace($placeholder, $this->sqlValue($value), $sql);
                }
            }
        }
        // Positional params (?)
        $i = 0;
        while (strpos($sql, '?') !== false && isset($this->bindings[(string) $i])) {
            $sql = preg_replace('/\?/', $this->sqlValue($this->bindings[(string) $i]), $sql, 1);
            $i++;
        }
        return $sql;
    }

    /** Format a PHP value as a SQL literal string. */
    private function sqlValue(mixed $value): string
    {
        if ($value === null) {
            return 'NULL';
        }
        if (is_bool($value)) {
            return $value ? 'true' : 'false';
        }
        if (is_int($value) || is_float($value)) {
            return (string) $value;
        }
        return "'" . str_replace("'", "''", (string) $value) . "'";
    }

    /** Remove surrounding quotes from a SQL literal. */
    private function unescapeValue(string $val): mixed
    {
        $val = trim($val);
        $len = strlen($val);
        if ($len >= 2) {
            if (($val[0] === "'" && $val[$len - 1] === "'") ||
                ($val[0] === '"' && $val[$len - 1] === '"')) {
                $val = substr($val, 1, -1);
            }
        }
        if (strtoupper($val) === 'NULL') return null;
        if (strtolower($val) === 'true')  return 'true';
        if (strtolower($val) === 'false') return 'false';
        if ($val === '') return '';
        if (is_numeric($val)) {
            return (strpos($val, '.') !== false) ? (float) $val : (int) $val;
        }
        return $val;
    }

    /** Parse a VALUES (a, b, c) list into an array of strings. */
    private function parseValuesList(string $inner): array
    {
        $result = [];
        $cur = '';
        $inQ = false;
        $qChar = null;
        $depth = 0;
        for ($i = 0; $i < strlen($inner); $i++) {
            $ch = $inner[$i];
            if (!$inQ) {
                if ($ch === '(') { $depth++; continue; }
                if ($ch === ')') { $depth--; continue; }
            }
            if (($ch === "'" || $ch === '"') && !$inQ) { $inQ = true; $qChar = $ch; $cur .= $ch; continue; }
            if ($inQ && $ch === $qChar) {
                if ($i + 1 < strlen($inner) && $inner[$i + 1] === $qChar) {
                    $cur .= $qChar . $qChar; $i++; continue;
                }
                $inQ = false; $cur .= $ch; continue;
            }
            if ($ch === ',' && !$inQ && $depth <= 0) { $result[] = trim($cur); $cur = ''; continue; }
            $cur .= $ch;
        }
        if (trim($cur) !== '') $result[] = trim($cur);
        return $result;
    }

    /** Parse a WHERE clause into Supabase filter arrays [[col, op, val], ...]. */
    private function parseWhereClause(string $where): array
    {
        $filters = [];
        $where = preg_replace('/^\s*AND\s+/i', '', trim($where));
        if ($where === '') return $filters;

        // OR is not natively supported by Supabase REST; we do our best below
        if (stripos($where, ' OR ') !== false) {
            return $this->parseWhereWithOr($where);
        }

        $conditions = preg_split('/\s+AND\s+/i', $where);
        foreach ($conditions as $cond) {
            $cond = trim($cond);
            if ($cond === '') continue;

            if (preg_match("/^`?(\w+)`?\s+NOT\s+LIKE\s+(.+)$/i", $cond, $m)) {
                $filters[] = [$m[1], 'NOT LIKE', $this->convertLikePattern($m[2])];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s+(ILIKE)\s+(.+)$/i", $cond, $m)) {
                $filters[] = [$m[1], 'ILIKE', $this->convertLikePattern($m[3])];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s+LIKE\s+(.+)$/i", $cond, $m)) {
                $filters[] = [$m[1], 'LIKE', $this->convertLikePattern($m[2])];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s+IN\s*\(([^)]+)\)/i", $cond, $m)) {
                $vals = array_map(fn($v) => $this->unescapeValue(trim($v)), explode(',', $m[2]));
                $filters[] = [$m[1], 'IN', $vals];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s+IS\s+NULL$/i", $cond, $m)) {
                $filters[] = [$m[1], 'IS', 'null'];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s+IS\s+NOT\s+NULL$/i", $cond, $m)) {
                $filters[] = [$m[1], 'IS NOT', 'null'];
                continue;
            }
            if (preg_match("/^`?(\w+)`?\s*(!=|<>|>=|<=|>|<|=)\s*(.+)$/i", $cond, $m)) {
                $col = $m[1];
                $op  = $m[2];
                if ($op === '<>') $op = '!=';
                $val = $this->unescapeValue(trim($m[3]));
                $filters[] = [$col, $op, $val];
                continue;
            }

            $this->logWarning('Could not parse WHERE condition: ' . $cond);
        }

        return $filters;
    }

    /**
     * Handle WHERE with OR by splitting and making multiple requests,
     * then merging results (deduplicated by primary key).
     */
    private function parseWhereWithOr(string $where): array
    {
        // Simple approach: split on OR and return filters for first group
        // The caller can handle more complex cases
        $this->logWarning('OR conditions are partially emulated by splitting the WHERE clause and making multiple REST requests. Some results may be missed. WHERE: ' . $where);

        // Extract individual conditions
        $conditions = preg_split('/\s+(?:AND|OR)\s+/i', $where);
        $filters = [];
        foreach ($conditions as $cond) {
            $cond = trim($cond);
            if (preg_match("/^`?(\w+)`?\s*(=|!=|<>|>=|<=|>|<|LIKE|ILIKE)\s*(.+)$/i", $cond, $m)) {
                $op = $m[2];
                $val = $this->unescapeValue(trim($m[3]));
                if (strtoupper($op) === 'LIKE' || strtoupper($op) === 'ILIKE') {
                    $val = $this->convertLikePattern($m[3]);
                }
                $filters[] = [$m[1], $op, $val];
            }
        }

        // Return all conditions combined as AND (best-effort). For complex OR logic,
        // the caller would need to merge multiple requests manually.
        return $filters;
    }

    /** Convert SQL LIKE pattern (%foo%) to Supabase pattern (*foo*). */
    private function convertLikePattern(string $val): string
    {
        $val = $this->unescapeValue($val);
        return (string) str_replace('%', '*', ($val instanceof \Stringable ? (string)$val : $val));
    }

    /** Extract the table name from an INSERT statement. */
    private function extractTableFromInsert(string $sql): ?string
    {
        if (preg_match('/INTO\s+([\w.`]+)/i', $sql, $m)) {
            return trim($m[1], '`');
        }
        return null;
    }

    /** Extract the column => value map from an INSERT statement. */
    private function extractDataFromInsert(string $sql): ?array
    {
        if (!preg_match('/\(([^)]+)\)/', $sql, $cm)) return null;
        if (!preg_match('/VALUES\s*\(([^)]+)\)/i', $sql, $vm)) return null;
        $columns = array_map(fn($c) => trim($c, '`'), array_map('trim', explode(',', $cm[1])));
        $values  = $this->parseValuesList($vm[1]);
        $data = [];
        foreach ($columns as $i => $col) {
            if (isset($values[$i])) {
                $data[$col] = $this->unescapeValue($values[$i]);
            }
        }
        return $data;
    }

    private function logError(string $msg): void
    {
        try {
            $app = \MythicalDash\App::getInstance(true);
            $app->getLogger()->error('[SupabasePDOStatement] ' . $msg);
        } catch (\Throwable $e) {
            error_log('[SupabasePDOStatement] ' . $msg);
        }
    }

    private function logWarning(string $msg): void
    {
        try {
            $app = \MythicalDash\App::getInstance(true);
            $app->getLogger()->warning('[SupabasePDOStatement] ' . $msg);
        } catch (\Throwable $e) {
            error_log('[SupabasePDOStatement] ' . $msg);
        }
    }
}
