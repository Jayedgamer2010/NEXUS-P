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
use MythicalDash\Hooks\Pterodactyl\Admin\Servers;
use MythicalDash\Services\ServerTime\TimeCreditManager;
use MythicalDash\Services\ServerTime\SessionManager;
use MythicalDash\Services\ServerTime\NodeSlotManager;
use MythicalDash\Services\ServerTime\TimePackageManager;

// GET /api/user/server/time/:uuid/status
$router->get('/api/user/server/time/(.*)/status', function (string $uuid): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyGET();
    $session = new Session($appInstance);
    if (!$session->isLoggedIn()) {
        $appInstance->Unauthorized('You must be logged in', ['error_code' => 'NOT_LOGGED_IN']);
        return;
    }
    $userId = $session->getInfo(UserColumns::UUID, false);

    // Get time credits for this server
    $timeCredits = TimeCreditManager::getByUserAndServer($userId, $uuid);
    // Get active session for this server
    $currentSession = SessionManager::getActiveSession($userId, $uuid);
    // Get node slot info
    $nodeSlots = ['active' => 0, 'max' => 4, 'queue_length' => 0];
    if ($currentSession && isset($currentSession['node_id'])) {
        $nodeSlotData = NodeSlotManager::get($currentSession['node_id']);
        if ($nodeSlotData) {
            $nodeSlots = [
                'active' => (int) $nodeSlotData['current_active'],
                'max' => (int) $nodeSlotData['max_active'],
                'queue_length' => SessionManager::getQueueLength($currentSession['node_id']),
            ];
        }
    }

    // Calculate status and remaining time
    $minutesRemaining = $timeCredits ? (int) $timeCredits['minutes_remaining'] : 0;
    $status = 'suspended';
    $queuePosition = null;
    $cooldownUntil = null;
    $endsAt = null;

    if ($currentSession) {
        $sessionStatus = $currentSession['status'];
        if ($sessionStatus === 'active') {
            $status = 'active';
            $endsAt = $currentSession['ends_at'];
            // Calculate remaining from ends_at
            $endsAtTs = strtotime($endsAt);
            $now = time();
            if ($endsAtTs > $now) {
                $minutesRemaining = max(1, (int) ceil(($endsAtTs - $now) / 60));
            } else {
                $minutesRemaining = 0;
            }
        } elseif ($sessionStatus === 'queued') {
            $status = 'queued';
            $queuePosition = (int) $currentSession['queue_position'];
        } elseif ($sessionStatus === 'cooldown') {
            $status = 'cooldown';
            $cooldownUntil = $currentSession['cooldown_until'];
        } elseif ($sessionStatus === 'suspended') {
            $status = 'suspended';
        }
    }

    $appInstance->OK('Time status retrieved successfully.', [
        'minutes_remaining' => $minutesRemaining,
        'status' => $status,
        'queue_position' => $queuePosition,
        'cooldown_until' => $cooldownUntil,
        'ends_at' => $endsAt,
        'node_slots' => $nodeSlots,
        'first_start_available' => TimeCreditManager::isFirstStartAvailable($uuid),
    ]);
});

// GET /api/user/server/time/:uuid/queue-status
$router->get('/api/user/server/time/(.*)/queue-status', function (string $uuid): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyGET();
    $session = new Session($appInstance);
    if (!$session->isLoggedIn()) {
        $appInstance->Unauthorized('You must be logged in', ['error_code' => 'NOT_LOGGED_IN']);
        return;
    }
    $userId = $session->getInfo(UserColumns::UUID, false);

    $currentSession = SessionManager::getSessionByStatus($userId, $uuid, 'queued');
    if (!$currentSession) {
        $appInstance->BadRequest('No queued session found for this server', ['error_code' => 'NOT_QUEUED']);
        return;
    }

    $position = (int) $currentSession['queue_position'];
    $nodeId = $currentSession['node_id'];

    // Count how many sessions are ahead in queue for same node
    $aheadOfYou = SessionManager::countQueuedAhead($nodeId, (int) $currentSession['id']);

    // Get estimated wait time
    $estimatedWait = SessionManager::estimateWaitMinutes($nodeId, $position);

    // Get node info
    $nodeSlotData = NodeSlotManager::get($nodeId);
    $nodeActive = $nodeSlotData ? (int) $nodeSlotData['current_active'] : 0;
    $nodeMax = $nodeSlotData ? (int) $nodeSlotData['max_active'] : 4;

    $appInstance->OK('Queue status retrieved successfully.', [
        'position' => $position,
        'ahead_of_you' => $aheadOfYou,
        'estimated_wait_minutes' => $estimatedWait,
        'node_active' => $nodeActive,
        'node_max' => $nodeMax,
    ]);
});

// POST /api/user/server/time/:uuid/queue
$router->post('/api/user/server/time/(.*)/queue', function (string $uuid): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPOST();
    $session = new Session($appInstance);
    if (!$session->isLoggedIn()) {
        $appInstance->Unauthorized('You must be logged in', ['error_code' => 'NOT_LOGGED_IN']);
        return;
    }
    $userId = $session->getInfo(UserColumns::UUID, false);

    // Check if there's already an active session
    $existingSession = SessionManager::getActiveSession($userId, $uuid);
    if ($existingSession) {
        $appInstance->BadRequest('Server already has an active session', [
            'error_code' => 'SESSION_ALREADY_ACTIVE',
            'current_status' => $existingSession['status'],
        ]);
        return;
    }

    $packageId = $_POST['package_id'] ?? null;
    $minutesToAdd = 0;
    $isFirstStart = false;

    // Check if this is the first start
    if (!TimeCreditManager::hasFirstStartBeenUsed($uuid)) {
        // First start - free 10 minutes
        $isFirstStart = true;
        $minutesToAdd = (int) $appInstance->getConfig()->getDBSetting('free_first_start_minutes', 10);
        TimeCreditManager::markFirstStartUsed($uuid, $userId);
        // Update or create time credits
        TimeCreditManager::addMinutes($userId, $uuid, $minutesToAdd);
    } elseif ($packageId) {
        // Need to purchase a time package
        $package = TimePackageManager::getById((string) $packageId);
        if (!$package) {
            $appInstance->BadRequest('Invalid time package', ['error_code' => 'INVALID_PACKAGE']);
            return;
        }

        // Check user has enough coins
        $userCoins = $session->getInfo(UserColumns::COINS, false);
        if ((int) $userCoins < (int) $package['coin_cost']) {
            $appInstance->BadRequest('Insufficient coins', [
                'error_code' => 'INSUFFICIENT_COINS',
                'required' => (int) $package['coin_cost'],
                'available' => (int) $userCoins,
            ]);
            return;
        }

        $minutesToAdd = (int) $package['minutes'];

        // Deduct coins immediately
        $session->removeCreditsAtomic((int) $package['coin_cost']);

        // Record coin transaction (via direct query since we don't have the helper here)
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("INSERT INTO coin_transactions (user_id, amount, reason) VALUES (?, ?, ?)");
            $stmt->execute([
                $userId,
                -(int) $package['coin_cost'],
                "Time package: {$package['name']} - {$package['minutes']} minutes",
            ]);
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('Failed to record coin transaction: ' . $e->getMessage());
        }

        // Add time credits
        TimeCreditManager::addMinutes($userId, $uuid, $minutesToAdd);
    } else {
        $appInstance->BadRequest('First start already used. Must provide a package to purchase.', [
            'error_code' => 'FIRST_START_USED',
        ]);
        return;
    }

    // Now try to start the server - get node ID from Pterodactyl
    // We need to find the Pterodactyl server ID from our DB
    $serverInfoDb = MythicalDash\Chat\Servers\Server::getById($uuid);
    if (!$serverInfoDb) {
        // Try to get from Pterodactyl by UUID
        $pterodactylServers = Servers::getUserServers();
        $foundServer = null;
        foreach ($pterodactylServers['data'] ?? [] as $sv) {
            if (($sv['attributes']['external_id'] ?? '') === $uuid || ($sv['attributes']['description'] ?? '') === $uuid) {
                $foundServer = $sv;
                break;
            }
        }
        if (!$foundServer) {
            $appInstance->BadRequest('Server not found in Pterodactyl', ['error_code' => 'SERVER_NOT_FOUND']);
            return;
        }
        $nodeId = (string) ($foundServer['attributes']['relationships']['node']['attributes']['identifier'] ?? '');
        $pterodactylServerId = $foundServer['attributes']['id'];
    } else {
        $pterodactylServerId = $serverInfoDb['pterodactyl_id'];
        $nodeId = $serverInfoDb['node_id'] ?? '';
    }

    // If no node_id in DB, get it from Pterodactyl
    if (empty($nodeId) && $pterodactylServerId) {
        try {
            $svDetails = MythicalDash\Hooks\Pterodactyl\Admin\Servers::getServerPterodactylDetails((int) $pterodactylServerId);
            if ($svDetails) {
                $nodeId = (string) ($svDetails['attributes']['relationships']['node']['attributes']['identifier'] ?? '');
            }
        } catch (\Exception $e) {
            // Fall back to a default node
            $nodeId = $uuid;
        }
    }

    if (empty($nodeId)) {
        $nodeId = $uuid;
    }

    // Ensure node_slots entry exists
    $defaultMaxSlots = (int) $appInstance->getConfig()->getDBSetting('default_node_max_slots', 4);
    NodeSlotManager::ensureExists($nodeId, $defaultMaxSlots);

    // Check if a slot is available
    $nodeSlotData = NodeSlotManager::get($nodeId);
    $slotsAvailable = $nodeSlotData ? (int) $nodeSlotData['current_active'] < (int) $nodeSlotData['max_active'] : true;

    $sessionStatus = 'queued';
    $queuePosition = 0;
    $startedAt = null;
    $endsAt = null;

    if ($slotsAvailable) {
        // Start immediately
        $sessionStatus = 'active';
        $startedAt = date('Y-m-d H:i:s');
        $endsAt = date('Y-m-d H:i:s', strtotime("+{$minutesToAdd} minutes"));

        // Increment node slot counter
        NodeSlotManager::incrementActive($nodeId);

        // Try to unsuspend and start the server via Pterodactyl
        if ($pterodactylServerId) {
            try {
                MythicalDash\Hooks\Pterodactyl\Admin\Servers::unsuspendPterodactylServer((int) $pterodactylServerId);
                MythicalDash\Hooks\Pterodactyl\Admin\Servers::sendPterodactylPowerAction((int) $pterodactylServerId, 'start');
            } catch (\Exception $e) {
                App::getInstance(false)->getLogger()->error('Failed to start server: ' . $e->getMessage());
            }
        }
    } else {
        // Add to queue
        $queuePosition = SessionManager::getQueueLength($nodeId) + 1;
    }

    // Create the session
    $newSession = SessionManager::create($userId, $uuid, $nodeId, $sessionStatus, $queuePosition, $startedAt, $endsAt);

    $appInstance->OK($isFirstStart ? 'First start - server added to queue/started.' : 'Time purchased - server queued/started.', [
        'minutes_remaining' => $minutesToAdd,
        'status' => $sessionStatus,
        'queue_position' => $queuePosition,
        'ends_at' => $endsAt,
        'first_start' => $isFirstStart,
        'session_id' => $newSession['id'] ?? null,
    ]);
});

// POST /api/user/server/time/:uuid/leave-queue
$router->post('/api/user/server/time/(.*)/leave-queue', function (string $uuid): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPOST();
    $session = new Session($appInstance);
    if (!$session->isLoggedIn()) {
        $appInstance->Unauthorized('You must be logged in', ['error_code' => 'NOT_LOGGED_IN']);
        return;
    }
    $userId = $session->getInfo(UserColumns::UUID, false);

    // Check if this is the first start (can't leave queue during a free first start)
    $currentSession = SessionManager::getActiveSession($userId, $uuid);
    if (!$currentSession) {
        $appInstance->BadRequest('No active session found for this server', ['error_code' => 'NO_ACTIVE_SESSION']);
        return;
    }

    if ($currentSession['status'] === 'active') {
        $appInstance->BadRequest('Cannot leave queue while server is active', ['error_code' => 'SESSION_ACTIVE']);
        return;
    }

    if ($currentSession['status'] === 'cooldown') {
        $appInstance->BadRequest('Cannot leave queue while in cooldown', ['error_code' => 'SESSION_COOLDOWN']);
        return;
    }

    // Calculate refund amount based on package that was purchased
    $refundAmount = 0;
    if (isset($currentSession['first_start']) && $currentSession['first_start']) {
        // First start was free - no refund needed
        $refundAmount = 0;
    } else {
        // Try to find the coin transaction for this session to refund
        try {
            $pdo = App::getInstance(false)->getDatabase()->getPdo();
            $stmt = $pdo->prepare("SELECT amount FROM coin_transactions WHERE user_id = ? AND reason LIKE 'Time package:%' ORDER BY created_at DESC LIMIT 1");
            $stmt->execute([$userId]);
            $tx = $stmt->fetch(PDO::FETCH_ASSOC);
            if ($tx) {
                $refundAmount = abs((int) $tx['amount']);
            }

            // Record refund transaction
            if ($refundAmount > 0) {
                $stmt2 = $pdo->prepare("INSERT INTO coin_transactions (user_id, amount, reason) VALUES (?, ?, ?)");
                $stmt2->execute([
                    $userId,
                    $refundAmount,
                    "Queue leave refund",
                ]);
            }
        } catch (\Exception $e) {
            App::getInstance(false)->getLogger()->error('Failed to process refund: ' . $e->getMessage());
        }
    }

    // Refund coins
    if ($refundAmount > 0) {
        $session->addCreditsAtomic($refundAmount);
    }

    // Get node ID before deleting session
    $nodeId = $currentSession['node_id'] ?? '';

    // Delete the session
    SessionManager::delete((int) $currentSession['id']);

    // Recalculate queue positions for remaining queued servers on this node
    if (!empty($nodeId)) {
        SessionManager::recalculateQueuePositions($nodeId);
    }

    $appInstance->OK('Left queue successfully' . ($refundAmount > 0 ? " with $refundAmount coins refunded." : '.'), [
        'refunded_coins' => $refundAmount,
    ]);
});
