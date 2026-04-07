<template>
    <div class="rounded-xl border border-gray-800 bg-[#0d0d17] p-4 text-center">
        <div class="mb-2 text-sm font-medium text-gray-400">Queue Position</div>
        <div class="mb-3 text-5xl font-bold text-yellow-400">#{{ queueInfo.position }}</div>

        <!-- Progress dots -->
        <div class="mb-3 flex items-center justify-center gap-1">
            <span v-for="i in queueInfo.node_max" :key="i"
                class="h-3 w-3 rounded-full"
                :class="i <= queueInfo.node_active ? 'bg-green-500' : 'bg-gray-700'">
            </span>
            <span v-for="i in queueInfo.position" :key="'q' + i"
                class="h-3 w-3 rounded-full animate-pulse bg-yellow-500/50">
            </span>
        </div>

        <p class="mb-1 text-sm text-gray-400">
            {{ queueInfo.ahead_of_you }} server{{ queueInfo.ahead_of_you !== 1 ? 's' : '' }} ahead of you
        </p>
        <p class="mb-3 text-sm text-gray-500">
            Estimated wait: ~{{ queueInfo.estimated_wait_minutes }} minutes
        </p>

        <!-- Waiting spinner -->
        <div class="mb-4 flex items-center justify-center">
            <svg class="h-5 w-5 animate-spin text-yellow-400" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"></path>
            </svg>
        </div>

        <button @click="$emit('leave')"
            class="w-full rounded-lg bg-red-600/80 px-4 py-2 text-xs font-medium text-white transition hover:bg-red-700">
            Leave Queue (full refund)
        </button>
    </div>
</template>

<script setup lang="ts">
interface QueueInfo {
    position: number;
    ahead_of_you: number;
    estimated_wait_minutes: number;
    node_active: number;
    node_max: number;
}

defineProps<{
    queueInfo: QueueInfo;
}>();

defineEmits<{
    leave: [];
}>();
</script>
