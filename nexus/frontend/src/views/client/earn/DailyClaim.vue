<template>
    <LayoutDashboard>
        <div class="p-6">
            <!-- Header Section -->
            <div class="flex items-center justify-between mb-6">
                <div>
                    <h1 class="text-2xl font-bold text-gray-100">{{ $t('daily_claim.title') }}</h1>
                    <p class="text-gray-400 mt-1">{{ $t('daily_claim.subtitle') }}</p>
                </div>
                <div class="px-4 py-2 bg-[#1A1825] rounded-lg border border-[#2a2a3f]/30">
                    <span class="text-sm text-gray-400">{{ $t('daily_claim.balance') }}:</span>
                    <span class="ml-2 text-lg font-bold text-green-400">{{ balance }}</span>
                    <span class="ml-1 text-sm text-gray-400">{{ $t('daily_claim.coins') }}</span>
                </div>
            </div>

            <!-- Error Message -->
            <div v-if="error" class="mb-6 p-4 bg-red-500/10 border border-red-500/20 rounded-lg">
                <p class="text-red-400">{{ error }}</p>
            </div>

            <!-- Success Message -->
            <div v-if="success" class="mb-6 p-4 bg-green-500/10 border border-green-500/20 rounded-lg">
                <p class="text-green-400">{{ success }}</p>
            </div>

            <!-- Streak Info -->
            <div v-if="claimInfo" class="mb-6 grid gap-6 sm:grid-cols-3">
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-3xl font-bold text-indigo-400">{{ claimInfo.streak }}</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('daily_claim.streak') }}</div>
                </div>
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-3xl font-bold text-green-400">{{ claimInfo.today_amount }}</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('daily_claim.today_reward') }}</div>
                </div>
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-3xl font-bold text-yellow-400">{{ claimInfo.total_claimed }}</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('daily_claim.total_claimed') }}</div>
                </div>
            </div>

            <!-- Claim Button -->
            <div class="flex flex-col items-center gap-6">
                <div v-if="!canClaim" class="bg-[#1A1825] rounded-xl p-8 border border-[#2a2a3f]/30 text-center max-w-md w-full">
                    <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-gray-700/50 flex items-center justify-center">
                        <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                    </div>
                    <h3 class="text-lg font-semibold text-gray-200 mb-2">{{ $t('daily_claim.already_claimed') }}</h3>
                    <p class="text-gray-400 mb-4">{{ $t('daily_claim.next_claim') }}:</p>
                    <div class="text-2xl font-bold text-indigo-400">{{ countdownText }}</div>
                </div>

                <button
                    v-else
                    @click="handleClaim"
                    :disabled="isClaiming"
                    class="group relative px-8 py-4 bg-gradient-to-r from-green-600 to-emerald-600 hover:from-green-500 hover:to-emerald-500 disabled:from-gray-600 disabled:to-gray-700 text-white text-lg font-semibold rounded-xl shadow-lg shadow-green-500/20 hover:shadow-green-500/40 transition-all duration-300 transform hover:scale-105"
                >
                    <span class="flex items-center gap-3">
                        <svg v-if="!isClaiming" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                        <svg v-else class="w-6 h-6 animate-spin" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                        </svg>
                        {{ isClaiming ? $t('daily_claim.claiming') : $t('daily_claim.claim_now') }}
                    </span>
                </button>
            </div>

            <!-- Claim History -->
            <div class="mt-8">
                <h2 class="text-lg font-semibold text-gray-200 mb-4">{{ $t('daily_claim.history') }}</h2>
                <div v-if="history.length === 0" class="bg-[#1A1825] rounded-lg p-6 border border-[#2a2a3f]/30 text-center">
                    <p class="text-gray-400">{{ $t('daily_claim.no_history') }}</p>
                </div>
                <div v-else class="space-y-2">
                    <div v-for="entry in history" :key="entry.created_at"
                        class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 flex items-center justify-between">
                        <div class="flex items-center gap-3">
                            <div class="w-10 h-10 rounded-full bg-green-500/20 flex items-center justify-center">
                                <svg class="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                        d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                </svg>
                            </div>
                            <div>
                                <div class="text-gray-200">{{ entry.reason }}</div>
                                <div class="text-xs text-gray-400">{{ formatDate(entry.created_at) }}</div>
                            </div>
                        </div>
                        <div class="text-lg font-bold text-green-400">+{{ entry.amount }}</div>
                    </div>
                </div>
            </div>
        </div>
    </LayoutDashboard>
</template>

<script setup lang="ts">
import LayoutDashboard from '@/components/client/LayoutDashboard.vue';
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { MythicalDOM } from '@/mythicaldash/MythicalDOM';
import Session from '@/mythicaldash/Session';

MythicalDOM.setPageTitle(MythicalDOM.getTranslation('daily_claim.title'));

const balance = ref(Session.getInfoInt('credits') ?? 0);
const error = ref<string | null>(null);
const success = ref<string | null>(null);
const isClaiming = ref(false);
const canClaim = ref(false);
const claimInfo = ref<{
    streak: number;
    today_amount: number;
    total_claimed: number;
} | null>(null);
const history = ref<Array<{ amount: number; reason: string; created_at: string }>>([]);

// Countdown timer
const timeUntilNextClaim = ref(0);
let countdownInterval: ReturnType<typeof setInterval> | null = null;

const countdownText = computed(() => {
    const seconds = timeUntilNextClaim.value;
    if (seconds <= 0) return '00:00:00';
    const h = Math.floor(seconds / 3600).toString().padStart(2, '0');
    const m = Math.floor((seconds % 3600) / 60).toString().padStart(2, '0');
    const s = (seconds % 60).toString().padStart(2, '0');
    return `${h}:${m}:${s}`;
});

const fetchClaimInfo = async () => {
    try {
        const response = await fetch('/api/user/daily-claim/info');
        if (!response.ok) throw new Error('Failed to fetch claim info');
        const data = await response.json();
        if (data.success) {
            claimInfo.value = data.claimInfo;
            canClaim.value = data.canClaim;
            timeUntilNextClaim.value = data.secondsUntilNext ?? 0;
            history.value = data.history ?? [];
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Failed to fetch claim info';
    }
};

const handleClaim = async () => {
    if (isClaiming.value) return;
    isClaiming.value = true;
    error.value = null;
    success.value = null;

    try {
        const response = await fetch('/api/user/daily-claim/claim', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
        });
        const data = await response.json();

        if (data.success) {
            success.value = `+${data.amount} ${MythicalDOM.getTranslation('daily_claim.coins')}!`;
            balance.value += data.amount;
            canClaim.value = false;
            timeUntilNextClaim.value = 86400; // 24 hours
            await fetchClaimInfo(); // refresh history
        } else {
            error.value = data.message || 'Failed to claim';
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Failed to claim';
    } finally {
        isClaiming.value = false;
    }
};

const formatDate = (dateStr: string) => {
    return new Date(dateStr).toLocaleDateString(undefined, {
        month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit',
    });
};

onMounted(() => {
    fetchClaimInfo();
    countdownInterval = setInterval(() => {
        if (timeUntilNextClaim.value > 0 && !canClaim.value) {
            timeUntilNextClaim.value--;
        }
    }, 1000);
});

onUnmounted(() => {
    if (countdownInterval) clearInterval(countdownInterval);
});
</script>
