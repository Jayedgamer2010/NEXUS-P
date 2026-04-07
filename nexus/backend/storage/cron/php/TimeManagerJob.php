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

namespace MythicalDash\Cron;

use MythicalDash\Chat\TimedTask;
use MythicalDash\Services\ServerTime\SessionManager;
use MythicalDash\Services\ServerTime\TimeCreditManager;
use MythicalDash\Services\ServerTime\NodeSlotManager;

class TimeManagerJob implements TimeTask
{
    public function run()
    {
        $cron = new Cron('time-manager', '1M');
        try {
            $cron->runIfDue(function () {
                $app = \MythicalDash\App::getInstance(false, true);
                $logger = $app->getLogger();

                $logger->info('=== Time Manager Job Started ===');

                $this->checkExpiredSessions($logger);
                $this->processQueue($logger);
                $this->checkCooldowns($logger);

                $logger->info('=== Time Manager Job Completed ===');
                TimedTask::markRun('time-manager', true, 'Time manager job completed');
            }, true);
        } catch (\Exception $e) {
            $app = \MythicalDash\App::getInstance(false, true);
            $app->getLogger()->error('Failed to run time manager job: ' . $e->getMessage());
            TimedTask::markRun('time-manager', false, 'Time manager job failed: ' . $e->getMessage());
        }
    }

    /**
     * Find and handle expired sessions.
     */
    private function checkExpiredSessions($logger): void
    {
        $expiredSessions = SessionManager::getExpiredSessions();

        if (empty($expiredSessions)) {
            $logger->info('No expired sessions found.');
            return;
        }

        $logger->info('Processing ' . count($expiredSessions) . ' expired session(s).');

        foreach ($expiredSessions as $session) {
            $logger->info(
                "Ending session: {$session['id']} for server {$session['server_uuid']} on node {$session['node_id']}"
            );

            // Stop and suspend the server via Pterodactyl
            $this->stopAndSuspendServer($session, $logger);

            // Set session status to cooldown
            $app = \MythicalDash\App::getInstance(false, true);
            $cooldownMinutes = (int) $app->getConfig()->getDBSetting('cooldown_minutes', 2);
            $cooldownUntil = date('Y-m-d H:i:s', strtotime("+{$cooldownMinutes} minutes"));

            SessionManager::update((int) $session['id'], [
                'status' => 'cooldown',
                'cooldown_until' => $cooldownUntil,
            ]);

            // Decrement node slots
            NodeSlotManager::decrementActive($session['node_id']);

            // Set time credits to 0
            TimeCreditManager::setRemaining($session['server_uuid'], 0);

            $logger->info(
                "Session {$session['id']} moved to cooldown until {$cooldownUntil}"
            );
        }
    }

    /**
     * Process queued sessions when slots become available.
     */
    private function processQueue($logger): void
    {
        $nodeSlots = NodeSlotManager::getAll();

        foreach ($nodeSlots as $slot) {
            $nodeId = $slot['node_id'];
            $currentActive = (int) $slot['current_active'];
            $maxActive = (int) $slot['max_active'];

            if ($currentActive >= $maxActive) {
                continue;
            }

            // Get oldest queued session for this node
            $queuedSession = SessionManager::getOldestQueuedForNode($nodeId);

            if (!$queuedSession) {
                continue;
            }

            $logger->info(
                "Starting queued session: {$queuedSession['id']} for server {$queuedSession['server_uuid']}"
            );

            // Get time credits to calculate ends_at
            $timeCredits = TimeCreditManager::getByUserAndServer(
                $queuedSession['user_id'],
                $queuedSession['server_uuid']
            );

            $minutes = $timeCredits ? (int) $timeCredits['minutes_remaining'] : 10;

            $startedAt = date('Y-m-d H:i:s');
            $endsAt = date('Y-m-d H:i:s', strtotime("+{$minutes} minutes"));

            // Update session to active
            SessionManager::update((int) $queuedSession['id'], [
                'status' => 'active',
                'started_at' => $startedAt,
                'ends_at' => $endsAt,
                'queue_position' => null,
            ]);

            // Increment node slots
            NodeSlotManager::incrementActive($nodeId);

            // Recalculate queue positions
            SessionManager::recalculateQueuePositions($nodeId);

            // Try to unsuspend and start the server
            $this->unsuspendAndStartServer($queuedSession, $logger);

            $logger->info("Session {$queuedSession['id']} started, ends at {$endsAt}");
        }
    }

    /**
     * Process finished cooldowns -> suspended state.
     */
    private function checkCooldowns($logger): void
    {
        $finishedCooldowns = SessionManager::getFinishedCooldowns();

        if (empty($finishedCooldowns)) {
            $logger->info('No finished cooldowns found.');
            return;
        }

        $logger->info('Processing ' . count($finishedCooldowns) . ' finished cooldown(s).');

        foreach ($finishedCooldowns as $session) {
            SessionManager::update((int) $session['id'], [
                'status' => 'suspended',
                'cooldown_until' => null,
            ]);

            $logger->info("Session {$session['id']} moved to suspended state");
        }
    }

    /**
     * Stop and suspend a server via Pterodactyl API.
     */
    private function stopAndSuspendServer(array $session, $logger): void
    {
        try {
            $app = \MythicalDash\App::getInstance(false, true);
            $config = $app->getConfig();
            $pterodactylUrl = $config->getDBSetting(\MythicalDash\Config\ConfigInterface::PTERODACTYL_BASE_URL, '');
            $pterodactylKey = $config->getDBSetting(\MythicalDash\Config\ConfigInterface::PTERODACTYL_API_KEY, '');

            if (empty($pterodactylUrl) || empty($pterodactylKey)) {
                $logger->warning('Pterodactyl API not configured, skipping server stop/suspend');
                return;
            }

            $server = MythicalDash\Hooks\Pterodactyl\Admin\Servers::getServerPterodactylDetails((int) $session['pterodactyl_id']);
            if (!$server) {
                // Try to find by UUID
                $server = MythicalDash\Hooks\Pterodactyl\Admin\Servers::getServerByUuid($session['server_uuid']);
                if (!$server) {
                    $logger->warning("Could not find Pterodactyl server for session {$session['id']}");
                    return;
                }
            }

            $pterodactylServerId = $server['attributes']['id'] ?? null;
            if (!$pterodactylServerId) {
                return;
            }

            // Send stop signal
            MythicalDash\Hooks\Pterodactyl\Admin\Servers::sendPterodactylPowerAction((int) $pterodactylServerId, 'stop');
            $logger->info("Stop signal sent to Pterodactyl server {$pterodactylServerId}");

            // Suspend after stopping
            MythicalDash\Hooks\Pterodactyl\Admin\Servers::suspendPterodactylServer((int) $pterodactylServerId);
            $logger->info("Pterodactyl server {$pterodactylServerId} suspended");
        } catch (\Exception $e) {
            $logger->error('Failed to stop/suspend server: ' . $e->getMessage());
        }
    }

    /**
     * Unsuspend and start a server via Pterodactyl API.
     */
    private function unsuspendAndStartServer(array $session, $logger): void
    {
        try {
            // Get server details
            $serverInfoDb = MythicalDash\Chat\Servers\Server::getById($session['server_uuid']);
            if (!$serverInfoDb) {
                $logger->warning("Server {$session['server_uuid']} not found in MythicalDash DB");
                return;
            }

            $pterodactylServerId = $serverInfoDb['pterodactyl_id'] ?? null;
            if (!$pterodactylServerId) {
                $logger->warning("No Pterodactyl server ID for {$session['server_uuid']}");
                return;
            }

            // Unsuspend
            MythicalDash\Hooks\Pterodactyl\Admin\Servers::unsuspendPterodactylServer((int) $pterodactylServerId);
            $logger->info("Pterodactyl server {$pterodactylServerId} unsuspended");

            // Start
            MythicalDash\Hooks\Pterodactyl\Admin\Servers::sendPterodactylPowerAction((int) $pterodactylServerId, 'start');
            $logger->info("Start signal sent to Pterodactyl server {$pterodactylServerId}");
        } catch (\Exception $e) {
            $logger->error('Failed to unsuspend/start server: ' . $e->getMessage());
        }
    }
}
