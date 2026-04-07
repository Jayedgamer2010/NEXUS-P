<?php

namespace MythicalDash\Services\ServerTime;

use MythicalDash\App;
use PDO;

class TimePackageManager
{
    public static function getById(string $id): ?array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM time_packages WHERE id = ? LIMIT 1");
            $stmt->execute([$id]);
            $result = $stmt->fetch(PDO::FETCH_ASSOC);
            return $result ?: null;
        } catch (\Exception $e) {
            return null;
        }
    }

    public static function getAll(): array
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT * FROM time_packages ORDER BY coin_cost ASC");
            $stmt->execute();
            return $stmt->fetchAll(PDO::FETCH_ASSOC);
        } catch (\Exception $e) {
            return [];
        }
    }

    public static function create(string $name, int $minutes, int $coinCost): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("INSERT INTO time_packages (name, minutes, coin_cost) VALUES (?, ?, ?)");
            return $stmt->execute([$name, $minutes, $coinCost]);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function update(string $id, array $data): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $setClauses = [];
            $params = [];
            foreach ($data as $key => $value) {
                $setClauses[] = "$key = ?";
                $params[] = $value;
            }
            $params[] = $id;
            $stmt = $pdo->prepare("UPDATE time_packages SET " . implode(', ', $setClauses) . " WHERE id = ?");
            return $stmt->execute($params);
        } catch (\Exception $e) {
            return false;
        }
    }

    public static function delete(string $id): bool
    {
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("DELETE FROM time_packages WHERE id = ?");
            return $stmt->execute([$id]);
        } catch (\Exception $e) {
            return false;
        }
    }
}
