<?php

namespace MythicalDash\Services\ServerTime;

use MythicalDash\App;
use PDO;

class NodeSlotManager
{
    public static function get(string $nodeId): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM node_slots WHERE node_id = ? LIMIT 1");
            $stmt->execute([$nodeId]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            return null;
        }
    }

    public static function ensureExists(string $nodeId, int $maxActive = 4): bool
    {
        try {
            if (self::get($nodeId)) {
                return true;
            }
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("INSERT INTO node_slots (node_id, max_active, current_active) VALUES (?, ?, 0) ON CONFLICT (node_id) DO NOTHING");
            return $stmt->execute([$nodeId, $maxActive]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function incrementActive(string $nodeId): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("UPDATE node_slots SET current_active = current_active + 1, updated_at = CURRENT_TIMESTAMP WHERE node_id = ?");
            return $stmt->execute([$nodeId]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function decrementActive(string $nodeId): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("UPDATE node_slots SET current_active = GREATEST(current_active - 1, 0), updated_at = CURRENT_TIMESTAMP WHERE node_id = ?");
            return $stmt->execute([$nodeId]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function updateMaxSlots(string $nodeId, int $maxActive): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("UPDATE node_slots SET max_active = ?, updated_at = CURRENT_TIMESTAMP WHERE node_id = ?");
            return $stmt->execute([$maxActive, $nodeId]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function getAll(): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM node_slots ORDER BY node_id");
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return [];
        }
    }

    public static function getSlotsWithQueue(): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("
                SELECT ns.*, COUNT(ss.id) as queue_length
                FROM node_slots ns
                LEFT JOIN server_sessions ss ON ss.node_id = ns.node_id AND ss.status = 'queued'
                GROUP BY ns.id
                ORDER BY ns.node_id
            ");
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return self::getAll();
        }
    }
}
