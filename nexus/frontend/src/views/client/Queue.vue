<template>
    <LayoutDashboard>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-2xl font-bold text-white">Server Queue</h1>
                <p class="text-sm text-gray-400 mt-1">Manage your server queue positions and active sessions</p>
            </div>
        </div>

        <div v-if="loading" class="flex justify-center items-center py-16">
            <LoaderCircle class="h-8 w-8 animate-spin text-purple-400" />
        </div>

        <div v-else-if="servers.length === 0"
            class="rounded-xl border border-gray-800 bg-[#0d0d17] py-16 text-center">
            <Clock class="mx-auto mb-4 h-12 w-12 text-gray-600" />
            <h3 class="text-lg font-semibold text-gray-400">No servers in queue</h3>
            <p class="text-sm text-gray-500 mt-1">Create a server to join the queue</p>
            <router-link to="/server/create"
                class="mt-4 inline-block rounded-lg bg-gradient-to-r from-purple-600 to-purple-500 px-6 py-2 text-sm font-medium text-white transition hover:from-purple-700 hover:to-purple-600">
                Create Server
            </router-link>
        </div>

        <div v-else class="space-y-4">
            <div v-for="server in servers" :key="server.uuid"
                class="rounded-xl border border-gray-800 bg-[#0d0d17] p-4">
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-3 min-w-0">
                        <!-- Status indicator -->
                        <span class="relative flex h-3 w-3 shrink-0">
                            <span class="absolute inline-flex h-full w-full animate-ping rounded-full opacity-75"
                                :class="statusColors[server.status]?.dot"></span>
                            <span class="relative inline-flex h-3 w-3 rounded-full"
                                :class="statusColors[server.status]?.bg"></span>
                        </span>
                        <div class="min-w-0">
                            <p class="font-semibold text-white truncate">{{ server.name }}</p>
                            <p class="text-xs text-gray-500 truncate">{{ server.node_id || 'Unknown node' }}</p>
                        </div>
                    </div>

                    <div class="flex items-center gap-4">
                        <!-- Status badge -->
                        <span class="rounded-full px-3 py-1 text-xs font-semibold uppercase"
                            :class="statusColors[server.status]?.badge">
                            {{ server.status }}
                        </span>

                        <!-- Queue position or time remaining -->
                        <span class="text-sm font-mono" :class="statusColors[server.status]?.text">
                            <template v-if="server.status === 'queued'">
                                #{{ server.queue_position }}
                            </template>
                            <template v-else-if="server.status === 'active' && server.ends_at">
                                {{ formatTimeRemaining(server.ends_at) }}
                            </template>
                            <template v-else-if="server.status === 'cooldown' && server.cooldown_until">
                                {{ formatCooldown(server.cooldown_until) }}
                            </template>
                        </span>

                        <!-- Action button -->
                        <button v-if="server.status === 'queued'" @click="leaveQueue(server.uuid)"
                            class="shrink-0 rounded-lg bg-red-600/80 px-3 py-1.5 text-xs font-medium text-white transition hover:bg-red-700">
                            Leave Queue
                        </button>
                        <span v-else-if="server.status === 'active'" class="shrink-0 text-xs text-green-400">
                            Running
                        </span>
                        <span v-else-if="server.status === 'cooldown'" class="shrink-0 text-xs text-orange-400">
                            Cooling down
                        </span>
                        <span v-else class="shrink-0 text-xs text-gray-500">
                            Suspended
                        </span>
                    </div>
                </div>
            </div>
        </div>
    </LayoutDashboard>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue';
import LayoutDashboard from '@/components/client/LayoutDashboard.vue';
import { Clock, LoaderCircle } from 'lucide-vue-next';
import TimeSystem from '@/mythicaldash/client/TimeSystem';

interface QueuedServer {
    uuid: string;
    name: string;
    node_id: string;
    status: 'queued' | 'active' | 'cooldown' | 'suspended';
    queue_position: number | null;
    ends_at: string | null;
    cooldown_until: string | null;
}

const loading = ref(true);
const servers = ref<QueuedServer[]>([]);
let pollingInterval: number | null = null;

const statusColors: Record<string, { dot: string; bg: string; badge: string; text: string }> = {
    active: { dot: 'bg-green-400', bg: 'bg-green-500', badge: 'bg-green-500/20 text-green-400', text: 'text-green-400' },
    queued: { dot: 'bg-yellow-400', bg: 'bg-yellow-500', badge: 'bg-yellow-500/20 text-yellow-400', text: 'text-yellow-400' },
    cooldown: { dot: 'bg-orange-400', bg: 'bg-orange-500', badge: 'bg-orange-500/20 text-orange-400', text: 'text-orange-400' },
    suspended: { dot: '', bg: 'bg-gray-500', badge: 'bg-gray-500/20 text-gray-400', text: 'text-gray-400' },
};

async function fetchServers() {
    try {
        // First get user's servers
        const serverResp = await fetch('/api/user/session');
        const serverData = await serverResp.json();
        if (serverData.success && serverData.servers) {
            const serverList = serverData.servers;

            // Then get time status for each
            const withTime = [];
            for (const sv of serverList) {
                try {
                    const timeResp = await fetch(`/api/user/server/time/${sv.uuid}/status`);
                    const timeData = await timeResp.json();
                    if (timeData.success && timeData.time_status?.status !== 'suspended') {
                        withTime.push({
                            uuid: sv.uuid,
                            name: sv.name || 'Server',
                            node_id: timeData.time_status.node_slots ? `Node` : '',
                            ...timeData.time_status,
                        });
                    }
                } catch {
                    // Skip servers without time status
                }
            }
            servers.value = withTime;
        }
    } catch (e) {
        console.error('Failed to fetch queue:', e);
    } finally {
        loading.value = false;
    }
}

function formatTimeRemaining(endsAt: string): string {
    const diffMs = new Date(endsAt).getTime() - Date.now();
    if (diffMs <= 0) return '00:00';
    const mins = Math.floor(diffMs / 60000);
    return `${mins}m remaining`;
}

function formatCooldown(cooldownUntil: string): string {
    const diffMs = new Date(cooldownUntil).getTime() - Date.now();
    if (diffMs <= 0) return 'Expired';
    const mins = Math.floor(diffMs / 60000);
    return `${mins}s cooldown`;
}

async function leaveQueue(serverUuid: string) {
    if (!confirm('Leave queue? You will receive a full coin refund.')) return;
    try {
        const result = await TimeSystem.leaveQueue(serverUuid);
        if (result.success) {
            await fetchServers();
        }
    } catch (e) {
        console.error('Failed to leave queue:', e);
    }
}

onMounted(() => {
    fetchServers();
    pollingInterval = window.setInterval(fetchServers, 5000);
});

onUnmounted(() => {
    if (pollingInterval) clearInterval(pollingInterval);
});
</script>
