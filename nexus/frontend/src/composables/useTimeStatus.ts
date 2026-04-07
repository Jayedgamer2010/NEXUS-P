import { ref, onMounted, onUnmounted } from 'vue';

export interface TimeStatus {
    minutes_remaining: number;
    status: 'queued' | 'active' | 'cooldown' | 'suspended';
    queue_position: number | null;
    cooldown_until: string | null;
    ends_at: string | null;
    node_slots: { active: number; max: number; queue_length: number };
    first_start_available: boolean;
}

export function useTimeStatus(serverId: string) {
    const timeStatus = ref<TimeStatus | null>(null);
    const loading = ref(true);
    const error = ref<string | null>(null);
    let pollingInterval: number | null = null;

    async function fetchStatus() {
        try {
            const response = await fetch(`/api/user/server/time/${serverId}/status`);
            const data = await response.json();
            if (data.success) {
                timeStatus.value = data.time_status;
                error.value = null;
            } else {
                error.value = data.error_code || 'Failed to fetch time status';
            }
        } catch (e) {
            error.value = 'Network error';
        } finally {
            loading.value = false;
        }
    }

    function startPolling(intervalMs = 5000) {
        fetchStatus();
        pollingInterval = window.setInterval(fetchStatus, intervalMs);
    }

    function stopPolling() {
        if (pollingInterval) {
            clearInterval(pollingInterval);
            pollingInterval = null;
        }
    }

    onMounted(() => startPolling());
    onUnmounted(() => stopPolling());

    return { timeStatus, loading, error, fetchStatus, startPolling, stopPolling };
}

export function formatCountdown(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
}
