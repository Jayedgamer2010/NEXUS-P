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
 *
 * Please rather than modifying the dashboard code try to report the thing you wish on our github or write a plugin
 */

namespace MythicalDash\Chat;

use MythicalDash\Database\SupabasePDO;
use MythicalDash\Database\SupabasePDOStatement;
use PDO;
use PDOException;

class Database
{
    private $pdo;
    private $mysqli;
    private $host;
    private $dbName;
    private $username;
    private $password;
    private $port;
    private ?\MythicalDash\Database\SupabaseClient $supabaseClient = null;

    /**
     * Database constructor.
     *
     * @param string $host the hostname or Supabase URL
     * @param string $dbName the database name (unused for Supabase mode)
     * @param string|null $username the username (unused for Supabase mode)
     * @param string|null $password the password (unused for Supabase mode)
     * @param int $port the port (unused for Supabase mode)
     * @param \MythicalDash\Database\SupabaseClient|null $supabaseClient pre-built client (optional)
     *
     * @throws \Exception if connection fails
     */
    public function __construct(
        $host,
        $dbName,
        $username = null,
        $password = null,
        int $port = 5432,
        ?\MythicalDash\Database\SupabaseClient $supabaseClient = null,
    ) {
        if ($supabaseClient !== null) {
            $this->supabaseClient = $supabaseClient;
            $this->pdo = new SupabasePDO($supabaseClient);
        } elseif (!empty($_ENV['SUPABASE_URL']) && !empty($_ENV['SUPABASE_SERVICE_KEY'])) {
            // Supabase mode
            $this->supabaseClient = new \MythicalDash\Database\SupabaseClient($_ENV['SUPABASE_URL'], $_ENV['SUPABASE_SERVICE_KEY']);
            $this->pdo = new SupabasePDO($this->supabaseClient);
        } else {
            // Fallback: traditional MySQL via PDO
            $dsn = "mysql:host=$host;port=$port;dbname=$dbName";
            try {
                $this->pdo = new \PDO($dsn, $username, $password);
                $this->pdo->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
            } catch (\PDOException $e) {
                throw new \Exception('Connection failed: ' . $e->getMessage());
            }
        }

        $this->host = $host;
        $this->dbName = $dbName;
        $this->username = $username;
        $this->password = $password;
        $this->port = $port;
    }

    public function getPdo(): \PDO
    {
        return $this->pdo;
    }

    public function getMysqli(): \mysqli
    {
        if ($this->supabaseClient !== null) {
            throw new \Exception('mysqli is not available in Supabase mode. Use getPdo() for the PDO-compatible wrapper instead.');
        }
        return new \mysqli($this->host, $this->username, $this->password, $this->dbName);
    }

    /**
     * Check whether we are running in Supabase mode.
     */
    public function isSupabaseMode(): bool
    {
        return $this->supabaseClient !== null;
    }

    /**
     * Get the SupabaseClient directly if needed.
     */
    public function getSupabaseClient(): ?\MythicalDash\Database\SupabaseClient
    {
        return $this->supabaseClient;
    }

    /**
     * Create a SupabasePDO instance (PDO-compatible wrapper).
     * Models calling getPdoConnection() receive this and can use
     * prepare(), query(), fetch(), fetchAll(), etc. as before.
     *
     * @return SupabasePDO|\PDO
     */
    public static function getPdoConnection(): \PDO
    {
        \MythicalDash\App::getInstance(true)->loadEnv();

        $supabaseClient = null;
        if (!empty($_ENV['SUPABASE_URL']) && !empty($_ENV['SUPABASE_SERVICE_KEY'])) {
            $supabaseClient = new \MythicalDash\Database\SupabaseClient($_ENV['SUPABASE_URL'], $_ENV['SUPABASE_SERVICE_KEY']);
            return new SupabasePDO($supabaseClient);
        }

        // Fallback: real MySQL PDO
        $host     = $_ENV['DATABASE_HOST'] ?? '127.0.0.1';
        $database = $_ENV['DATABASE_DATABASE'] ?? 'mythicaldash';
        $user     = $_ENV['DATABASE_USER'] ?? 'root';
        $pass     = $_ENV['DATABASE_PASSWORD'] ?? '';
        $port     = (int) ($_ENV['DATABASE_PORT'] ?? 3306);
        $dsn = new \PDO("mysql:host=$host;port=$port;dbname=$database", $user, $pass);
        $dsn->setAttribute(\PDO::ATTR_ERRMODE, \PDO::ERRMODE_EXCEPTION);
        return $dsn;
    }

    // -----------------------------------------------------------------------
    //  Static convenience methods
    // -----------------------------------------------------------------------

    /**
     * Get the table row count.
     */
    public static function getTableRowCount(string $table, bool $adminSide = false): int
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare('SELECT COUNT(*) FROM ' . $table . ($adminSide ? ' WHERE deleted = "false"' : ''));
            $stmt->execute();
            return (int) $stmt->fetchColumn();
        } catch (\Exception $e) {
            self::db_Error('Failed to get table row count: ' . $e->getMessage());
            return 0;
        }
    }

    /**
     * Get the count of rows in a table with optional WHERE conditions.
     *
     * @param string $table  the table name
     * @param array  $where  where conditions  (key=>value for equals, or [column, operator, value])
     * @param bool   $includeDeleted  whether to include soft-deleted rows
     */
    public static function getTableColumnCount(string $table, array $where = [], bool $includeDeleted = false): int
    {
        try {
            $db = self::getPdoConnection();

            $conditions = [];
            $params = [];

            if (!$includeDeleted) {
                $conditions[] = "deleted = 'false'";
            }

            foreach ($where as $key => $value) {
                if (is_array($value)) {
                    $conditions[] = "{$value[0]} {$value[1]} ?";
                    $params[] = $value[2];
                } else {
                    $conditions[] = "$key = ?";
                    $params[] = $value;
                }
            }

            $whereClause = !empty($conditions) ? 'WHERE ' . implode(' AND ', $conditions) : '';

            $stmt = $db->prepare('SELECT COUNT(*) FROM ' . $table . ' ' . $whereClause);
            $stmt->execute($params);

            return (int) $stmt->fetchColumn();
        } catch (\Exception $e) {
            self::db_Error('Failed to get table column count: ' . $e->getMessage());
            return 0;
        }
    }

    /**
     * Check if a table exists (always returns true in Supabase mode).
     */
    public static function tableExists(string $table): bool
    {
        $db = self::getPdoConnection();
        if ($db instanceof SupabasePDO) {
            // In Supabase mode we cannot introspect tables; assume true
            return true;
        }
        try {
            $query = $db->query("SHOW TABLES LIKE '$table'");
            return $query->rowCount() > 0;
        } catch (\Exception $e) {
            self::db_Error('Failed to check if table exists: ' . $e->getMessage());
            return false;
        }
    }

    /**
     * Get all tables (not available in Supabase mode).
     */
    public static function getTables(): array
    {
        $db = self::getPdoConnection();
        if ($db instanceof SupabasePDO) {
            self::db_Error('getTables() is not available in Supabase mode.');
            return [];
        }
        try {
            return $db->query('SHOW TABLES')->fetchAll(\PDO::FETCH_COLUMN);
        } catch (\Exception $e) {
            self::db_Error('Failed to get tables: ' . $e->getMessage());
            return [];
        }
    }

    // -----------------------------------------------------------------------
    //  CRUD helpers  (use the SupabaseClient directly for efficiency)
    // -----------------------------------------------------------------------

    /**
     * Create a record in a table.
     *
     * @param string $table the table name
     * @param array $data associative array of column=>value pairs
     * @return int|false the id of the inserted row, or false on failure
     */
    public static function insertRecord(string $table, array $data)
    {
        try {
            $db = self::getPdoConnection();
            $fields = array_keys($data);
            $placeholders = array_map(fn($f) => ':' . $f, $fields);
            $sql = 'INSERT INTO ' . $table . ' (' . implode(',', $fields) . ') VALUES (' . implode(',', $placeholders) . ')';
            $stmt = $db->prepare($sql);
            if ($stmt->execute($data)) {
                $id = $db->lastInsertId();
                return $id !== '0' ? (int) $id : 0;
            }
            return false;
        } catch (\Exception $e) {
            self::db_Error('Failed to insert record into ' . $table . ': ' . $e->getMessage());
            return false;
        }
    }

    /**
     * Update records in a table.
     *
     * @param string $table the table name
     * @param array $data  column=>value pairs to update
     * @param array $where where conditions  (key=>value or [column, operator, value])
     */
    public static function updateRecord(string $table, array $data, array $where = []): int
    {
        try {
            $db = self::getPdoConnection();
            $setParts = [];
            foreach ($data as $col => $val) {
                $setParts[] = "$col = :u_$col";
            }
            $sql = "UPDATE $table SET " . implode(', ', $setParts);

            $conditions = [];
            $params = [];
            foreach ($where as $col => $val) {
                $conditions[] = "$col = :w_$col";
                $params["w_$col"] = $val;
            }
            if (!empty($conditions)) {
                $sql .= ' WHERE ' . implode(' AND ', $conditions);
            }

            // Merge set params
            foreach ($data as $col => $val) {
                $params["u_$col"] = $val;
            }

            $stmt = $db->prepare($sql);
            $stmt->execute($params);
            return $stmt->rowCount();
        } catch (\Exception $e) {
            self::db_Error('Failed to update record in ' . $table . ': ' . $e->getMessage());
            return 0;
        }
    }

    /**
     * Delete records from a table.
     *
     * @param string $table the table name
     * @param array $where where conditions
     */
    public static function deleteRecord(string $table, array $where = []): int
    {
        try {
            $db = self::getPdoConnection();
            $sql = "DELETE FROM $table";

            $conditions = [];
            $params = [];
            foreach ($where as $col => $val) {
                if (is_array($val)) {
                    $conditions[] = "{$val[0]} {$val[1]} :w_$col";
                } else {
                    $conditions[] = "$col = :w_$col";
                }
                $params["w_$col"] = $val;
            }
            if (!empty($conditions)) {
                $sql .= ' WHERE ' . implode(' AND ', $conditions);
            }

            $stmt = $db->prepare($sql);
            $stmt->execute($params);
            return $stmt->rowCount();
        } catch (\Exception $e) {
            self::db_Error('Failed to delete record from ' . $table . ': ' . $e->getMessage());
            return 0;
        }
    }

    /**
     * Fetch rows from a table.
     *
     * @param string $table  the table name
     * @param array  $where  where conditions
     * @param string $select columns (default '*')
     * @param string $orderBy  e.g. 'id DESC'
     * @param int    $limit
     * @param int    $offset
     * @return array
     */
    public static function getRows(
        string $table,
        array $where = [],
        string $select = '*',
        string $orderBy = '',
        int $limit = 0,
        int $offset = 0,
    ): array {
        try {
            $db = self::getPdoConnection();
            $sql = "SELECT $select FROM $table";

            $conditions = [];
            $params = [];
            foreach ($where as $col => $val) {
                if (is_array($val)) {
                    $conditions[] = "{$val[0]} {$val[1]} :w_$col";
                } else {
                    $conditions[] = "$col = :w_$col";
                }
                $params["w_$col"] = $val;
            }
            if (!empty($conditions)) {
                $sql .= ' WHERE ' . implode(' AND ', $conditions);
            }
            if (!empty($orderBy)) {
                $sql .= " ORDER BY $orderBy";
            }
            if ($limit > 0) {
                $sql .= " LIMIT $limit";
            }
            if ($offset > 0) {
                $sql .= " OFFSET $offset";
            }

            $stmt = $db->prepare($sql);
            $stmt->execute($params);
            return $stmt->fetchAll(\PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            self::db_Error('Failed to get rows from ' . $table . ': ' . $e->getMessage());
            return [];
        }
    }

    /**
     * Fetch a single row.
     */
    public static function getRow(
        string $table,
        array $where = [],
        string $select = '*',
    ): ?array {
        $rows = self::getRows($table, $where, $select, '', 1);
        return !empty($rows) ? $rows[0] : null;
    }

    /**
     * Check if any row matching the conditions exists.
     */
    public static function checkIfExist(string $table, array $where = []): bool
    {
        try {
            return self::getTableRowCount($table, false) > 0;
        } catch (\Exception $e) {
            return false;
        }
    }

    /**
     * Run a raw SQL query and return results.
     * NOTE: Supabase REST API cannot execute arbitrary SQL.
     *       Only basic SELECT/INSERT/UPDATE/DELETE will work automatically.
     *       SHOW TABLES, JOINs, GROUP BY etc. will return empty results.
     *
     * @param string $sql raw SQL query
     */
    public static function runSQL(string $sql): array
    {
        try {
            $db = self::getPdoConnection();
            if ($db instanceof SupabasePDO) {
                self::db_Error('WARNING: runSQL() called with Supabase backend. Complex SQL may not work. Query: ' . $sql);
            }
            $stmt = $db->prepare($sql);
            $stmt->execute();
            return $stmt->fetchAll(\PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            self::db_Error('Failed to run SQL: ' . $e->getMessage());
            return [];
        }
    }

    /**
     * Alias for runSQL.
     */
    public static function rawQuery(string $query): array
    {
        return self::runSQL($query);
    }

    /**
     * Get the last insert id (approximate for Supabase UUID tables).
     */
    public static function getLastInsertId(string $table): int
    {
        try {
            $db = self::getPdoConnection();
            return (int) $db->lastInsertId();
        } catch (\Exception $e) {
            self::db_Error('Failed to get last insert ID: ' . $e->getMessage());
            return 0;
        }
    }

    // -----------------------------------------------------------------------
    //  Soft-delete / lock helpers
    // -----------------------------------------------------------------------

    public static function markRecordAsDeleted(string $table, int $row): void
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("UPDATE $table SET deleted = 'true' WHERE id = :id");
            $stmt->execute(['id' => $row]);
        } catch (\Exception $e) {
            self::db_Error('Failed to mark record as deleted: ' . $e->getMessage());
        }
    }

    public static function getDeletedRecords(string $table): array
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("SELECT * FROM $table WHERE deleted = 'true'");
            $stmt->execute();
            return $stmt->fetchAll(\PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            self::db_Error('Failed to get deleted records: ' . $e->getMessage());
            return [];
        }
    }

    public static function restoreRecord(string $table, int $row): void
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("UPDATE $table SET deleted = 'false' WHERE id = :id");
            $stmt->execute(['id' => $row]);
        } catch (\Exception $e) {
            self::db_Error('Failed to restore record: ' . $e->getMessage());
        }
    }

    public static function lockRecord(string $table, int $row): void
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("UPDATE $table SET locked = 'true' WHERE id = :id");
            $stmt->execute(['id' => $row]);
        } catch (\Exception $e) {
            self::db_Error('Failed to lock record: ' . $e->getMessage());
        }
    }

    public static function unlockRecord(string $table, int $row): void
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("UPDATE $table SET locked = 'false' WHERE id = :id");
            $stmt->execute(['id' => $row]);
        } catch (\Exception $e) {
            self::db_Error('Failed to unlock record: ' . $e->getMessage());
        }
    }

    public static function isLocked(string $table, int $row): bool
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("SELECT locked FROM $table WHERE id = :id LIMIT 1");
            $stmt->execute(['id' => $row]);
            $result = $stmt->fetch(\PDO::FETCH_ASSOC);
            return $result && ($result['locked'] == 'true');
        } catch (\Exception $e) {
            self::db_Error('Failed to check for lock: ' . $e->getMessage());
            return false;
        }
    }

    public static function requestSaveAndUnlock(string $table, int $row): void
    {
        try {
            $db = self::getPdoConnection();
            $stmt = $db->prepare("UPDATE $table SET locked = 'false' WHERE id = :id");
            $stmt->execute(['id' => $row]);
        } catch (\Exception $e) {
            self::db_Error('Failed to request save and unlock: ' . $e->getMessage());
        }
    }

    /**
     * Count query: returns the count of rows matching a condition.
     *
     * @param string $table        the table name
     * @param array  $conditions   where conditions (key=>value pairs)
     */
    public static function countQuery(string $table, array $conditions = []): int
    {
        try {
            $db = self::getPdoConnection();
            $sql = "SELECT COUNT(*) as cnt FROM $table";

            $params = [];
            $whereParts = [];
            foreach ($conditions as $col => $val) {
                $whereParts[] = "$col = :w_$col";
                $params["w_$col"] = $val;
            }
            if (!empty($whereParts)) {
                $sql .= ' WHERE ' . implode(' AND ', $whereParts);
            }

            $stmt = $db->prepare($sql);
            $stmt->execute($params);
            return (int) $stmt->fetchColumn();
        } catch (\Exception $e) {
            self::db_Error('Failed to count rows in ' . $table . ': ' . $e->getMessage());
            return 0;
        }
    }

    /**
     * Log a database error.
     */
    public static function db_Error(string $message): void
    {
        $app = \MythicalDash\App::getInstance(true);
        $app->getLogger()->error($message, true);
    }
}
