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
 * SupabasePDO - A PDO-compatible wrapper around the Supabase REST API client.
 *
 * This class extends \PDO so it satisfies any type-hint that expects a PDO object.
 * Internally it translates prepare()/query()/exec() into Supabase REST API calls.
 */
class SupabasePDO extends \PDO
{
    private SupabaseClient $client;
    private string $_lastInsertId = '0';

    public function __construct(SupabaseClient $client)
    {
        // Initialise the parent with a dummy DSN so instanceof checks pass
        parent::__construct('sqlite::memory:');
        $this->client = $client;
    }

    public function getClient(): SupabaseClient
    {
        return $this->client;
    }

    // -----------------------------------------------------------------------
    //  PDO interface stubs
    // -----------------------------------------------------------------------

    public function prepare(string $query, array $options = []): SupabasePDOStatement
    {
        return new SupabasePDOStatement($this, $this->client, $query);
    }

    public function query(string $query, ?int $fetchMode = null, mixed ...$fetchModeArgs): \PDOStatement|false
    {
        $stmt = $this->prepare($query);
        $stmt->execute();
        return $stmt;
    }

    public function exec(string $query): int
    {
        $stmt = $this->prepare($query);
        return $stmt->execute() ? $stmt->rowCount() : 0;
    }

    public function lastInsertId(?string $name = null): string
    {
        return $this->_lastInsertId;
    }

    public function setLastInsertId(string $id): void
    {
        $this->_lastInsertId = $id;
    }

    public function beginTransaction(): bool
    {
        return true;  // Supabase REST does not support multi-statement transactions
    }

    public function commit(): bool { return true; }
    public function rollBack(): bool { return true; }

    public function setAttribute(int $attribute, mixed $value): bool { return true; }
    public function getAttribute(int $attribute): mixed { return null; }

    public function quote(string $string, int $parameter_type = \PDO::PARAM_STR): string
    {
        return "'" . str_replace("'", "''", $string) . "'";
    }

    public function inTransaction(): bool { return false; }

    public function errorCode(): ?string { return null; }
    public function errorInfo(): array { return ['00000', null, null]; }
}
