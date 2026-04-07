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
use MythicalDash\Config\ConfigInterface;
use MythicalDash\Chat\columns\UserColumns;

/**
 * GET /api/user/daily-claim/info
 * Returns claim info including streak, today's reward amount, and whether user can claim
 */
$router->get('/api/user/daily-claim/info', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $s = new Session($appInstance);
    $uuid = $s->getInfo(UserColumns::UUID, false);
    $userId = (int) $s->getInfo(UserColumns::ID, false);
    $db = $appInstance->db->getPdo();

    // Get last claim time from coin_transactions
    $stmt = $db->prepare('SELECT created_at FROM coin_transactions WHERE user_id = :uid AND reason = :reason ORDER BY created_at DESC LIMIT 1');
    $stmt->execute(['uid' => $userId, 'reason' => 'Daily Claim']);
    $lastClaim = $stmt->fetch(\PDO::FETCH_ASSOC);

    $canClaim = true;
    $secondsUntilNext = 0;

    if ($lastClaim) {
        $lastClaimTime = strtotime($lastClaim['created_at']);
        $now = time();
        $diff = $now - $lastClaimTime;

        if ($diff < 86400) { // 24 hours
            $canClaim = false;
            $secondsUntilNext = 86400 - $diff;
        }
    }

    // Calculate streak
    $streak = 0;
    if ($lastClaim) {
        $stmt = $db->prepare('SELECT created_at FROM coin_transactions WHERE user_id = :uid AND reason = :reason ORDER BY created_at DESC');
        $stmt->execute(['uid' => $userId, 'reason' => 'Daily Claim']);
        $claims = $stmt->fetchAll(\PDO::FETCH_ASSOC);

        $streak = 1;
        $prevDate = date('Y-m-d', strtotime($claims[0]['created_at']));
        for ($i = 1; $i < count($claims); $i++) {
            $currDate = date('Y-m-d', strtotime($claims[$i]['created_at']));
            $expected = date('Y-m-d', strtotime($prevDate . ' -1 day'));
            if ($currDate === $expected) {
                $streak++;
                $prevDate = $currDate;
            } else {
                break;
            }
        }
    }

    // Today's amount (base + streak bonus)
    $coinsPerDay = (int) $appInstance->getConfig()->getDBSetting(ConfigInterface::DAILY_COINS, 10);
    $streakBonus = min($streak, 10); // Cap bonus at 10
    $todayAmount = $coinsPerDay + ($streakBonus > 1 ? $streakBonus : 0);

    // Total claimed
    $stmt = $db->prepare("SELECT COALESCE(SUM(amount), 0) as total FROM coin_transactions WHERE user_id = :uid AND reason = 'Daily Claim'");
    $stmt->execute(['uid' => $userId]);
    $totalClaimed = (int) $stmt->fetchColumn();

    // Recent history
    $stmt = $db->prepare('SELECT amount, reason, created_at FROM coin_transactions WHERE user_id = :uid AND reason = :reason ORDER BY created_at DESC LIMIT 10');
    $stmt->execute(['uid' => $userId, 'reason' => 'Daily Claim']);
    $history = $stmt->fetchAll(\PDO::FETCH_ASSOC);

    App::OK([
        'claimInfo' => [
            'streak' => $streak,
            'today_amount' => $todayAmount,
            'total_claimed' => $totalClaimed,
        ],
        'canClaim' => $canClaim,
        'secondsUntilNext' => $secondsUntilNext,
        'history' => $history,
    ]);
});

/**
 * POST /api/user/daily-claim/claim
 * Claims the daily coins reward
 */
$router->post('/api/user/daily-claim/claim', function (): void {
    App::init();
    $appInstance = App::getInstance(true);
    $s = new Session($appInstance);
    $uuid = $s->getInfo(UserColumns::UUID, false);
    $userId = (int) $s->getInfo(UserColumns::ID, false);
    $db = $appInstance->db->getPdo();

    // Check last claim
    $stmt = $db->prepare('SELECT created_at FROM coin_transactions WHERE user_id = :uid AND reason = :reason ORDER BY created_at DESC LIMIT 1');
    $stmt->execute(['uid' => $userId, 'reason' => 'Daily Claim']);
    $lastClaim = $stmt->fetch(\PDO::FETCH_ASSOC);

    if ($lastClaim) {
        $lastClaimTime = strtotime($lastClaim['created_at']);
        $now = time();
        $diff = $now - $lastClaimTime;

        if ($diff < 86400) {
            App::Error('DAILY_ALREADY_CLAIMED', 'You have already claimed your daily coins today', $diff);
        }
    }

    // Calculate amount
    $coinsPerDay = (int) $appInstance->getConfig()->getDBSetting(ConfigInterface::DAILY_COINS, 10);

    // Calculate streak
    $streak = 0;
    if ($lastClaim) {
        $stmt = $db->prepare('SELECT created_at FROM coin_transactions WHERE user_id = :uid AND reason = :reason ORDER BY created_at DESC');
        $stmt->execute(['uid' => $userId, 'reason' => 'Daily Claim']);
        $claims = $stmt->fetchAll(\PDO::FETCH_ASSOC);

        $streak = 1;
        $prevDate = date('Y-m-d', strtotime($claims[0]['created_at']));
        for ($i = 1; $i < count($claims); $i++) {
            $currDate = date('Y-m-d', strtotime($claims[$i]['created_at']));
            $expected = date('Y-m-d', strtotime($prevDate . ' -1 day'));
            if ($currDate === $expected) {
                $streak++;
                $prevDate = $currDate;
            } else {
                break;
            }
        }
    } else {
        $streak = 1;
    }

    $streakBonus = min($streak, 10);
    $amount = $coinsPerDay + ($streakBonus > 1 ? $streakBonus : 0);

    // Add coins atomically
    if (!$s->addCreditsAtomic($amount)) {
        return;
    }

    // Record transaction
    $stmt = $db->prepare('INSERT INTO coin_transactions (user_id, amount, reason) VALUES (:uid, :amount, :reason)');
    $stmt->execute([
        'uid' => $userId,
        'amount' => $amount,
        'reason' => 'Daily Claim',
    ]);

    $newBalance = (int) $s->getInfo(UserColumns::CREDITS, false);
    $s->setInfo(UserColumns::CREDITS, $newBalance, false);

    App::OK([
        'amount' => $amount,
        'newBalance' => $newBalance,
        'streak' => $streak,
    ]);
});
