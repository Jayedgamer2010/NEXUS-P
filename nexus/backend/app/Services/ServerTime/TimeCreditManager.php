<?php

namespace MythicalDash\Services\ServerTime;

use MythicalDash\App;
use PDO;

class TimeCreditManager
{
    public static function getByUserAndServer(string $userId, string $serverUuid): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM server_time_credits WHERE user_id = ? AND server_uuid = ? LIMIT 1");
            $stmt->execute([$userId, $serverUuid]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            return null;
        }
    }

    public static function addMinutes(string $userId, string $serverUuid, int $minutes): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $existing = self::getByUserAndServer($userId, $serverUuid);
            if ($existing) {
                $newTotal = (int) $existing['minutes_remaining'] + $minutes;
                $newPurchased = (int) $existing['total_minutes_purchased'] + $minutes;
                $stmt = $pdo->prepare("UPDATE server_time_credits SET minutes_remaining = ?, total_minutes_purchased = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?");
                return $stmt->execute([$newTotal, $newPurchased, $existing['id']]);
            }
            $stmt = $pdo->prepare("INSERT INTO server_time_credits (user_id, server_uuid, minutes_remaining, total_minutes_purchased) VALUES (?, ?, ?, ?)");
            return $stmt->execute([$userId, $serverUuid, $minutes, $minutes]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function setRemaining(string $serverUuid, int $minutes): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("UPDATE server_time_credits SET minutes_remaining = ?, updated_at = CURRENT_TIMESTAMP WHERE server_uuid = ?");
            return $stmt->execute([$minutes, $serverUuid]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function isFirstStartAvailable(string $serverUuid): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT server_uuid FROM server_first_starts WHERE server_uuid = ? AND used = false LIMIT 1");
            $stmt->execute([$serverUuid]);
            return $stmt->fetch(PDO::FETCH_ASSOC) === false;
        } catch (\Exception $e) {
            return true;
        }
    }

    public static function hasFirstStartBeenUsed(string $serverUuid): bool
    {
        return !self::isFirstStartAvailable($serverUuid);
    }

    public static function markFirstStartUsed(string $serverUuid, ?string $userId = null): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("INSERT INTO server_first_starts (server_uuid, used) VALUES (?, true) ON CONFLICT (server_uuid) DO UPDATE SET used = true");
            return $stmt->execute([$serverUuid]);
        } catch (\Exception $e) {
            return false;
        }
    }
}
