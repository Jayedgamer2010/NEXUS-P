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
use MythicalDash\Permissions;
use MythicalDash\Chat\User\Session;
use MythicalDash\Middleware\PermissionMiddleware;
use MythicalDash\Services\ServerTime\NodeSlotManager;

// GET /api/admin/node-slots
$router->get('/api/admin/node-slots', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyGET();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    $slots = NodeSlotManager::getSlotsWithQueue();
    $appInstance->OK('Node slots retrieved successfully.', ['node_slots' => $slots]);
});

// PATCH /api/admin/node-slots/:nodeId/slots
$router->patch('/api/admin/node-slots/(.*)/slots', function (string $nodeId): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPATCH();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    if (!isset($_POST['max_active']) || (int) $_POST['max_active'] < 1) {
        $appInstance->BadRequest('max_active is required and must be at least 1', ['error_code' => 'INVALID_MAX_ACTIVE']);
        return;
    }

    $maxActive = (int) $_POST['max_active'];
    $existing = NodeSlotManager::get($nodeId);

    if (!$existing) {
        $appInstance->NotFound('Node slot not found', ['error_code' => 'NODE_SLOT_NOT_FOUND']);
        return;
    }

    $result = NodeSlotManager::updateMaxSlots($nodeId, $maxActive);
    if ($result) {
        $appInstance->OK('Node slots updated successfully.', ['error_code' => 'NODE_SLOTS_UPDATED', 'max_active' => $maxActive]);
    } else {
        $appInstance->BadRequest('Failed to update node slots', ['error_code' => 'FAILED_TO_UPDATE']);
    }
});

// GET /api/admin/time-packages
$router->get('/api/admin/time-packages', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyGET();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    $packages = MythicalDash\Services\ServerTime\TimePackageManager::getAll();
    $appInstance->OK('Time packages retrieved successfully.', ['time_packages' => $packages]);
});
