<template>
    <!-- Announcement Banner for Dashboard -->
    <div v-if="activeAnnouncements.length > 0" class="mb-6 space-y-3">
        <div
            v-for="announcement in activeAnnouncements"
            :key="announcement.id"
            :class="[
                'rounded-xl p-4 border transition-all duration-300',
                typeClasses[announcement.type] || typeClasses.info,
                'animate-fade-in',
            ]"
        >
            <div class="flex items-start gap-3">
                <!-- Icon -->
                <div class="flex-shrink-0 mt-0.5">
                    <component :is="typeIcons[announcement.type] || typeIcons.info" class="w-5 h-5" />
                </div>

                <!-- Content -->
                <div class="flex-1 min-w-0">
                    <h3 class="font-semibold text-sm text-gray-100">{{ announcement.title }}</h3>
                    <p class="text-xs text-gray-300 mt-1 line-clamp-2">{{ announcement.content }}</p>
                </div>

                <!-- Actions -->
                <div class="flex-shrink-0 flex items-center gap-2">
                    <button
                        v-if="announcement.content.length > 60"
                        @click="readMore(announcement)"
                        class="text-xs text-indigo-400 hover:text-indigo-300 whitespace-nowrap"
                    >
                        {{ $t('announcements.read_more') }}
                    </button>
                    <button
                        @click="dismissAnnouncement(announcement.id)"
                        class="w-6 h-6 flex items-center justify-center rounded-full hover:bg-gray-500/20 transition-colors"
                    >
                        <XIcon class="w-4 h-4 text-gray-400" />
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, h } from 'vue';
import { useRouter } from 'vue-router';
import { XIcon, AlertTriangle, Bell, Info, CheckCircle } from 'lucide-vue-next';
import Announcements from '@/mythicaldash/Announcements';

const router = useRouter();
const announcements = ref<Array<{ id: string; title: string; content: string; type: string; assets?: { images?: string[] } }>>([]);
const dismissedIds = ref<Set<string>>(new Set());

const activeAnnouncements = computed(() => {
    return announcements.value.filter(a => !dismissedIds.value.has(a.id));
});

const typeClasses: Record<string, string> = {
    info: 'bg-blue-500/10 border-blue-500/30',
    warning: 'bg-yellow-500/10 border-yellow-500/30',
    error: 'bg-red-500/10 border-red-500/30',
    success: 'bg-green-500/10 border-green-500/30',
    maintenance: 'bg-purple-500/10 border-purple-500/30',
};

const typeIcons: Record<string, any> = {
    info: () => h(Info, { class: 'text-blue-400' }),
    warning: () => h(AlertTriangle, { class: 'text-yellow-400' }),
    error: () => h(AlertTriangle, { class: 'text-red-400' }),
    success: () => h(CheckCircle, { class: 'text-green-400' }),
    maintenance: () => h(Bell, { class: 'text-purple-400' }),
};

const dismissAnnouncement = (id: string) => {
    dismissedIds.value.add(id);
    try {
        const dismissed = JSON.parse(localStorage.getItem('dismissed_announcements') || '[]');
        dismissed.push(id);
        localStorage.setItem('dismissed_announcements', JSON.stringify(dismissed));
    } catch {
        // ignore
    }
};

const readMore = (_announcement: { id: string }) => {
    router.push('/announcements');
};

onMounted(async () => {
    try {
        const data = await Announcements.fetchAnnouncements();
        if (data && Array.isArray(data)) {
            announcements.value = data;
        }
    } catch {
        // ignore
    }

    // Load dismissed IDs from localStorage
    try {
        const dismissed = JSON.parse(localStorage.getItem('dismissed_announcements') || '[]');
        dismissedIds.value = new Set(dismissed);
    } catch {
        // ignore
    }
});
</script>
