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

namespace MythicalDash\Services\ServerTime;

use MythicalDash\App;
use PDO;

class SessionManager
{
    /**
     * Get active session for a user and server.
     */
    public static function getActiveSession(string $userId, string $serverUuid): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE user_id = ? AND server_uuid = ? AND status IN ('active', 'queued', 'cooldown', 'suspended') ORDER BY created_at DESC LIMIT 1");
            $stmt->execute([$userId, $serverUuid]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::getActiveSession - ' . $e->getMessage());
            return null;
        }
    }

    /**
     * Get session by status for a user and server.
     */
    public static function getSessionByStatus(string $userId, string $serverUuid, string $status): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE user_id = ? AND server_uuid = ? AND status = ? ORDER BY created_at DESC LIMIT 1");
            $stmt->execute([$userId, $serverUuid, $status]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::getSessionByStatus - ' . $e->getMessage());
            return null;
        }
    }

    /**
     * Create a new session.
     */
    public static function create(string $userId, string $serverUuid, string $nodeId, string $status = 'queued', int $queuePosition = 0, ?string $startedAt = null, ?string $endsAt = null): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("INSERT INTO server_sessions (user_id, server_uuid, node_id, status, queue_position, started_at, ends_at) VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, created_at");
            $stmt->execute([$userId, $serverUuid, $nodeId, $status, $queuePosition, $startedAt, $endsAt]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return array_merge(['user_id' => $userId, 'server_uuid' => $serverUuid, 'node_id' => $nodeId, 'status' => $status, 'queue_position' => $queuePosition, 'started_at' => $startedAt, 'ends_at' => $endsAt], $result ?: []);
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::create - ' . $e->getMessage());
            return [];
        }
    }

    /**
     * Update session.
     */
    public static function update(int $sessionId, array $data): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $setClauses = [];
            $params = [];
            foreach ($data as $key => $value) {
                $setClauses[] = "$key = ?";
                $params[] = $value;
            }
            $setClauses[] = "updated_at = CURRENT_TIMESTAMP";
            $params[] = $sessionId;
            $stmt = $pdo->prepare("UPDATE server_sessions SET " . implode(', ', $setClauses) . " WHERE id = ?");
            return $stmt->execute($params);
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::update - ' . $e->getMessage());
            return false;
        }
    }

    /**
     * Delete a session.
     */
    public static function delete(int $sessionId): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("DELETE FROM server_sessions WHERE id = ?");
            return $stmt->execute([$sessionId]);
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::delete - ' . $e->getMessage());
            return false;
        }
    }

    /**
     * Get count of queued sessions for a node.
     */
    public static function getQueueLength(string $nodeId): int
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT COUNT(*) as cnt FROM server_sessions WHERE node_id = ? AND status = 'queued'");
            $stmt->execute([$nodeId]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return (int) ($result['cnt'] ?? 0);
        } catch (\Exception $e) {
            return 0;
        }
    }

    /**
     * Count queued sessions ahead of a given session.
     */
    public static function countQueuedAhead(string $nodeId, int $sessionId): int
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT COUNT(*) as cnt FROM server_sessions WHERE node_id = ? AND status = 'queued' AND created_at < (SELECT created_at FROM server_sessions WHERE id = ?)");
            $stmt->execute([$nodeId, $sessionId]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return (int) ($result['cnt'] ?? 0);
        } catch (\Exception $e) {
            return 0;
        }
    }

    /**
     * Estimate wait time in minutes based on sessions ahead.
     */
    public static function estimateWaitMinutes(string $nodeId, int $position): int
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            // Get average session time remaining for active sessions on this node
            $stmt = $pdo->prepare("SELECT AVG(EXTRACT(EPOCH FROM (ends_at - NOW())) / 60) as avg_remaining FROM server_sessions WHERE node_id = ? AND status = 'active' AND ends_at > NOW()");
            $stmt->execute([$nodeId]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            $avgRemaining = (float) ($result['avg_remaining'] ?? 10);
            return (int) ceil($avgRemaining * $position);
        } catch (\Exception $e) {
            return $position * 10;
        }
    }

    /**
     * Get all active sessions on a node.
     */
    public static function getActiveSessionsForNode(string $nodeId): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE node_id = ? AND status = 'active'");
            $stmt->execute([$nodeId]);
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return [];
        }
    }

    /**
     * Get oldest queued session for a node.
     */
    public static function getOldestQueuedForNode(string $nodeId): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE node_id = ? AND status = 'queued' ORDER BY created_at ASC LIMIT 1");
            $stmt->execute([$nodeId]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            return null;
        }
    }

    /**
     * Recalculate queue positions for a node.
     */
    public static function recalculateQueuePositions(string $nodeId): void
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT id FROM server_sessions WHERE node_id = ? AND status = 'queued' ORDER BY created_at ASC");
            $stmt->execute([$nodeId]);
            $sessions = $stmt->fetchAll(PDO::FETCH_ASSOC);
            $position = 1;
            foreach ($sessions as $session) {
                $update = $pdo->prepare("UPDATE server_sessions SET queue_position = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?");
                $update->execute([$position, (int) $session['id']]);
                $position++;
            }
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('SessionManager::recalculateQueuePositions - ' . $e->getMessage());
        }
    }

    /**
     * Get all expired active sessions.
     */
    public static function getExpiredSessions(): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE status = 'active' AND ends_at < CURRENT_TIMESTAMP");
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return [];
        }
    }

    /**
     * Get all finished cooldown sessions.
     */
    public static function getFinishedCooldowns(): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_sessions WHERE status = 'cooldown' AND cooldown_until < CURRENT_TIMESTAMP");
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return [];
        }
    }
}
