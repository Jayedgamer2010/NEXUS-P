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

use MythicalDash\App;
use MythicalDash\Chat\User\Session;
use MythicalDash\Chat\columns\UserColumns;

/**
 * GET /api/user/resource-packages
 * Returns available resource packages
 */
$router->get('/api/user/resource-packages', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $db = $appInstance->db->getPdo();

    $stmt = $db->prepare('SELECT id, name, description, memory, disk, cpu, coin_cost FROM resource_packages ORDER BY coin_cost ASC');
    $stmt->execute();
    $packages = $stmt->fetchAll(\PDO::FETCH_ASSOC);

    App::OK(['packages' => $packages]);
});

/**
 * GET /api/user/resource-packages/my-resources
 * Returns current user's extra resources
 */
$router->get('/api/user/resource-packages/my-resources', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $s = new Session($appInstance);
    $uuid = $s->getInfo(UserColumns::UUID, false);
    $userId = (int) $s->getInfo(UserColumns::ID, false);
    $db = $appInstance->db->getPdo();

    $stmt = $db->prepare('SELECT extra_memory, extra_disk, extra_cpu, extra_servers, updated_at FROM user_resources WHERE user_id = :uid');
    $stmt->execute(['uid' => $userId]);
    $resources = $stmt->fetch(\PDO::FETCH_ASSOC);

    if (!$resources) {
        $resources = [
            'extra_memory' => 0,
            'extra_disk' => 0,
            'extra_cpu' => 0,
            'extra_servers' => 0,
            'updated_at' => date('Y-m-d H:i:s'),
        ];
    }

    App::OK(['resources' => $resources]);
});

/**
 * POST /api/user/resource-packages/purchase
 * Purchases a resource package with coins
 */
$router->post('/api/user/resource-packages/purchase', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $s = new Session($appInstance);
    $uuid = $s->getInfo(UserColumns::UUID, false);
    $userId = (int) $s->getInfo(UserColumns::ID, false);
    $db = $appInstance->db->getPdo();

    // Get request body
    $input = json_decode(file_get_contents('php://input'), true);
    $packageId = $input['packageId'] ?? '';

    if (empty($packageId)) {
        App::Error('MISSING_PACKAGE_ID', 'No package ID provided', 400);
    }

    // Get package details
    $stmt = $db->prepare('SELECT * FROM resource_packages WHERE id = :id');
    $stmt->execute(['id' => $packageId]);
    $package = $stmt->fetch(\PDO::FETCH_ASSOC);

    if (!$package) {
        App::Error('PACKAGE_NOT_FOUND', 'Resource package not found', 404);
    }

    // Check user balance
    $currentCoins = (int) $s->getInfo(UserColumns::CREDITS, false);

    if ($currentCoins < $package['coin_cost']) {
        App::Error('NOT_ENOUGH_COINS', 'Not enough coins to purchase this package', 400);
    }

    // Deduct coins atomically
    if (!$s->addCreditsAtomic(-$package['coin_cost'])) {
        App::Error('TRANSACTION_FAILED', 'Failed to process transaction', 500);
    }

    // Get or create user_resources
    $stmt = $db->prepare('SELECT id FROM user_resources WHERE user_id = :uid');
    $stmt->execute(['uid' => $userId]);
    $existing = $stmt->fetch(\PDO::FETCH_ASSOC);

    if ($existing) {
        $stmt = $db->prepare('UPDATE user_resources SET extra_memory = extra_memory + :memory, extra_disk = extra_disk + :disk, extra_cpu = extra_cpu + :cpu, extra_servers = extra_servers + :servers, updated_at = NOW() WHERE user_id = :uid');
    } else {
        $stmt = $db->prepare('INSERT INTO user_resources (user_id, extra_memory, extra_disk, extra_cpu, extra_servers, updated_at) VALUES (:uid, :memory, :disk, :cpu, :servers, NOW())');
    }
    $stmt->execute([
        'uid' => $userId,
        'memory' => (int) $package['memory'],
        'disk' => (int) $package['disk'],
        'cpu' => (int) $package['cpu'],
        'servers' => 0, // Not implemented in packages yet
    ]);

    // Record transaction
    $stmt = $db->prepare('INSERT INTO coin_transactions (user_id, amount, reason) VALUES (:uid, :amount, :reason)');
    $stmt->execute([
        'uid' => $userId,
        'amount' => -$package['coin_cost'],
        'reason' => 'Purchase: ' . $package['name'],
    ]);

    $newBalance = (int) $s->getInfo(UserColumns::CREDITS, false);
    $s->setInfo(UserColumns::CREDITS, $newBalance, false);

    App::OK([
        'package' => $package['name'],
        'newBalance' => $newBalance,
        'resources_added' => [
            'memory' => (int) $package['memory'],
            'disk' => (int) $package['disk'],
            'cpu' => (int) $package['cpu'],
        ],
    ]);
});
