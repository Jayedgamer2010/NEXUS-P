<template>
    <div class="rounded-xl border border-gray-800 bg-[#0d0d17] p-4">
        <!-- ACTIVE -->
        <div v-if="status?.status === 'active'" class="flex flex-col gap-3">
            <div class="flex items-center gap-2">
                <span class="relative flex h-3 w-3">
                    <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-400 opacity-75"></span>
                    <span class="relative inline-flex h-3 w-3 rounded-full bg-green-500"></span>
                </span>
                <span class="text-sm font-semibold text-green-400">ACTIVE</span>
                <span :class="timeCritical ? 'text-red-400 font-bold animate-pulse' : 'text-gray-300'" class="ml-auto font-mono text-lg">
                    {{ countdownDisplay }}
                </span>
            </div>
            <div class="h-2 w-full overflow-hidden rounded-full bg-gray-800">
                <div class="h-full rounded-full bg-gradient-to-r from-purple-600 to-purple-400 transition-all duration-1000"
                    :style="{ width: progressPercent + '%' }"></div>
            </div>
            <div class="flex justify-between text-xs text-gray-500">
                <span>Started: {{ formatTime(sessionStart) }}</span>
                <span>Ends: {{ formatTime(endsAt) }}</span>
            </div>
        </div>

        <!-- QUEUED -->
        <div v-else-if="status?.status === 'queued'" class="flex flex-col gap-3">
            <div class="flex items-center gap-2">
                <span class="relative flex h-3 w-3">
                    <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-yellow-400 opacity-75"></span>
                    <span class="relative inline-flex h-3 w-3 rounded-full bg-yellow-500"></span>
                </span>
                <span class="text-sm font-semibold text-yellow-400">QUEUED</span>
                <span class="ml-auto text-lg font-bold text-yellow-400">#{{ status.queue_position }} in queue</span>
            </div>
            <div class="flex items-center justify-center py-2">
                <svg class="h-6 w-6 animate-spin text-yellow-400" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
                </svg>
            </div>
            <p class="text-center text-sm text-gray-400">Est. wait: ~{{ estimatedWait }} minutes</p>
            <button @click="$emit('leave-queue')"
                class="w-full rounded-lg bg-red-600 px-4 py-2 text-sm font-medium text-white transition hover:bg-red-700">
                Leave Queue (full refund)
            </button>
        </div>

        <!-- COOLDOWN -->
        <div v-else-if="status?.status === 'cooldown'" class="flex flex-col gap-3">
            <div class="flex items-center gap-2">
                <span class="relative flex h-3 w-3">
                    <span class="absolute inline-flex h-full w-full animate-ping rounded-full bg-orange-400 opacity-75"></span>
                    <span class="relative inline-flex h-3 w-3 rounded-full bg-orange-500"></span>
                </span>
                <span class="text-sm font-semibold text-orange-400">COOLDOWN</span>
                <span class="ml-auto font-mono text-lg text-orange-400">{{ cooldownCountdown }}</span>
            </div>
            <p class="text-sm text-gray-400">Server suspended - cooling down</p>
        </div>

        <!-- SUSPENDED -->
        <div v-else class="flex flex-col gap-3">
            <div class="flex items-center gap-2">
                <span class="inline-flex h-3 w-3 rounded-full bg-gray-500"></span>
                <span class="text-sm font-semibold text-gray-400">SUSPENDED</span>
            </div>
            <p class="text-sm text-gray-400">Server suspended - buy time to restart</p>
            <button @click="$emit('buy-time')"
                class="w-full rounded-lg bg-gradient-to-r from-purple-600 to-purple-500 px-4 py-2 text-sm font-medium text-white transition hover:from-purple-700 hover:to-purple-600">
                Buy Time
            </button>
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue';

export interface ServerTimeStatus {
    minutes_remaining: number;
    status: 'queued' | 'active' | 'cooldown' | 'suspended';
    queue_position: number | null;
    cooldown_until: string | null;
    ends_at: string | null;
    node_slots: { active: number; max: number; queue_length: number };
    first_start_available: boolean;
}

const props = defineProps<{
    status: ServerTimeStatus | null;
    estimatedWait?: number;
}>();

const emit = defineEmits<{
    'buy-time': [];
    'leave-queue': [];
}>();

const now = ref(Date.now());
let timerInterval: number | null = null;

const countdownDisplay = computed(() => {
    if (!props.status?.ends_at) return '00:00';
    const ends = new Date(props.status.ends_at).getTime();
    const remaining = Math.max(0, (ends - now.value) / 1000);
    const mins = Math.floor(remaining / 60);
    const secs = Math.floor(remaining % 60);
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
});

const timeCritical = computed(() => {
    if (!props.status?.ends_at) return false;
    const ends = new Date(props.status.ends_at).getTime();
    return (ends - Date.now()) / 1000 < 120;
});

const progressPercent = computed(() => {
    if (!props.status?.ends_at || !props.status.minutes_remaining) return 100;
    const totalMs = props.status.minutes_remaining * 60000;
    const remainingMs = Math.max(0, new Date(props.status.ends_at).getTime() - now.value);
    return Math.max(0, Math.min(100, (remainingMs / totalMs) * 100));
});

const cooldownCountdown = computed(() => {
    if (!props.status?.cooldown_until) return '00:00';
    const remaining = Math.max(0, (new Date(props.status.cooldown_until).getTime() - now.value) / 1000);
    const mins = Math.floor(remaining / 60);
    const secs = Math.floor(remaining % 60);
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
});

const sessionStart = computed(() => {
    if (!props.status?.ends_at || !props.status.minutes_remaining) return null;
    return new Date(new Date(props.status.ends_at).getTime() - props.status.minutes_remaining * 60000).toISOString();
});

const endsAt = computed(() => props.status?.ends_at ?? null);

function formatTime(iso: string | null): string {
    if (!iso) return '--:--';
    return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

onMounted(() => {
    timerInterval = window.setInterval(() => { now.value = Date.now(); }, 1000);
});

onUnmounted(() => {
    if (timerInterval) clearInterval(timerInterval);
});
</script>
