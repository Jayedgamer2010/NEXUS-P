<?php

use MythicalDash\App;
use MythicalDash\Permissions;
use MythicalDash\Chat\User\Session;
use MythicalDash\Middleware\PermissionMiddleware;
use MythicalDash\Services\ServerTime\TimePackageManager;

// GET /api/admin/time-packages
$router->get('/api/admin/time-packages', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyGET();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    $packages = TimePackageManager::getAll();
    $appInstance->OK('Time packages retrieved successfully.', ['time_packages' => $packages]);
});

// POST /api/admin/time-packages/create
$router->post('/api/admin/time-packages/create', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPOST();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    if (!isset($_POST['name']) || empty($_POST['name']) || !isset($_POST['minutes']) || !isset($_POST['coin_cost'])) {
        $appInstance->BadRequest('Name, minutes, and coin_cost are required', ['error_code' => 'MISSING_FIELDS']);
        return;
    }

    $name = htmlspecialchars(strip_tags($_POST['name']));
    $minutes = (int) $_POST['minutes'];
    $coinCost = (int) $_POST['coin_cost'];

    if ($minutes <= 0 || $coinCost <= 0) {
        $appInstance->BadRequest('Minutes and coin_cost must be positive', ['error_code' => 'INVALID_VALUES']);
        return;
    }

    $result = TimePackageManager::create($name, $minutes, $coinCost);
    if ($result) {
        $appInstance->OK('Time package created successfully.', ['error_code' => 'TIME_PACKAGE_CREATED']);
    } else {
        $appInstance->BadRequest('Failed to create time package', ['error_code' => 'FAILED_TO_CREATE']);
    }
});

// PATCH/POST /api/admin/time-packages/:id/update
$router->post('/api/admin/time-packages/(.*)/update', function (string $id): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPOST();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    $existing = TimePackageManager::getById($id);
    if (!$existing) {
        $appInstance->NotFound('Time package not found', ['error_code' => 'TIME_PACKAGE_NOT_FOUND']);
        return;
    }

    $data = [];
    if (isset($_POST['name']) && !empty($_POST['name'])) {
        $data['name'] = htmlspecialchars(strip_tags($_POST['name']));
    }
    if (isset($_POST['minutes'])) {
        $data['minutes'] = (int) $_POST['minutes'];
    }
    if (isset($_POST['coin_cost'])) {
        $data['coin_cost'] = (int) $_POST['coin_cost'];
    }

    if (empty($data)) {
        $appInstance->BadRequest('No fields to update', ['error_code' => 'NO_FIELDS']);
        return;
    }

    $result = TimePackageManager::update($id, $data);
    if ($result) {
        $appInstance->OK('Time package updated successfully.', ['error_code' => 'TIME_PACKAGE_UPDATED']);
    } else {
        $appInstance->BadRequest('Failed to update time package', ['error_code' => 'FAILED_TO_UPDATE']);
    }
});

// POST /api/admin/time-packages/:id/delete
$router->post('/api/admin/time-packages/(.*)/delete', function (string $id): void {
    App::init();
    $appInstance = App::getInstance(true);
    $appInstance->allowOnlyPOST();
    $session = new Session($appInstance);
    PermissionMiddleware::handle($appInstance, Permissions::ADMIN_SERVERS_LIST, $session);

    $existing = TimePackageManager::getById($id);
    if (!$existing) {
        $appInstance->NotFound('Time package not found', ['error_code' => 'TIME_PACKAGE_NOT_FOUND']);
        return;
    }

    $result = TimePackageManager::delete($id);
    if ($result) {
        $appInstance->OK('Time package deleted successfully.', ['error_code' => 'TIME_PACKAGE_DELETED']);
    } else {
        $appInstance->BadRequest('Failed to delete time package', ['error_code' => 'FAILED_TO_DELETE']);
    }
});
