<template>
    <LayoutDashboard>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-2xl font-bold text-pink-400">Node Slots</h1>
                <p class="text-sm text-gray-400 mt-1">Manage maximum active server slots per node</p>
            </div>
        </div>

        <div v-if="loading" class="flex justify-center items-center py-10">
            <LoaderCircle class="h-8 w-8 animate-spin text-pink-400" />
        </div>

        <div v-else class="space-y-4">
            <div v-for="slot in slots" :key="slot.node_id"
                class="rounded-xl border border-gray-800 bg-[#0d0d17] p-4">
                <div class="flex items-center justify-between mb-3">
                    <div>
                        <h3 class="font-semibold text-white">{{ slot.node_id }}</h3>
                        <p class="text-xs text-gray-500">Max: {{ slot.max_active }} slots</p>
                    </div>
                    <span class="text-sm text-gray-400">{{ slot.current_active }} / {{ slot.max_active }} active</span>
                </div>

                <!-- Progress bar -->
                <div class="mb-3 h-2 w-full overflow-hidden rounded-full bg-gray-800">
                    <div class="h-full rounded-full bg-gradient-to-r from-purple-600 to-purple-400 transition-all"
                        :style="{ width: ((slot.current_active / slot.max_active) * 100) + '%' }"></div>
                </div>

                <div class="flex items-center justify-between">
                    <p class="text-xs text-gray-500">Queue length: {{ slot.queue_length || 0 }}</p>
                    <button @click="openEditModal(slot)"
                        class="rounded-lg bg-gradient-to-r from-pink-500 to-violet-500 px-4 py-1.5 text-sm font-medium text-white transition hover:opacity-80">
                        Edit Slots
                    </button>
                </div>
            </div>
        </div>

        <!-- Edit Modal -->
        <div v-if="showEditModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="fixed inset-0 bg-black/60" @click="showEditModal = false"></div>
            <div class="relative w-full max-w-md rounded-2xl border border-gray-800 bg-[#0d0d17] p-6 shadow-2xl">
                <h3 class="mb-4 text-lg font-bold text-white">Edit Node Slots</h3>
                <p class="mb-2 text-sm text-gray-400">Node: {{ editingSlot?.node_id }}</p>
                <input v-model.number="editMax" type="number" min="1" max="20"
                    class="w-full rounded-lg border border-gray-700 bg-gray-900 px-3 py-2 text-white focus:border-purple-500 focus:outline-none mb-4" />
                <div class="flex gap-3">
                    <button @click="showEditModal = false"
                        class="flex-1 rounded-lg border border-gray-700 bg-gray-900 px-4 py-2 text-sm text-gray-300 hover:bg-gray-800">
                        Cancel
                    </button>
                    <button @click="saveSlots"
                        class="flex-1 rounded-lg bg-gradient-to-r from-purple-600 to-purple-500 px-4 py-2 text-sm font-medium text-white hover:from-purple-700 hover:to-purple-600">
                        Save
                    </button>
                </div>
            </div>
        </div>
    </LayoutDashboard>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import LayoutDashboard from '@/components/admin/LayoutDashboard.vue';
import { LoaderCircle } from 'lucide-vue-next';

interface NodeSlot {
    node_id: string;
    max_active: number;
    current_active: number;
    queue_length?: number;
}

const loading = ref(true);
const slots = ref<NodeSlot[]>([]);
const showEditModal = ref(false);
const editingSlot = ref<NodeSlot | null>(null);
const editMax = ref(4);

async function fetchSlots() {
    try {
        const res = await fetch('/api/admin/node-slots');
        const data = await res.json();
        if (data.success) {
            slots.value = data.node_slots || [];
        }
    } catch (e) {
        slots.value = [];
    } finally {
        loading.value = false;
    }
}

function openEditModal(slot: NodeSlot) {
    editingSlot.value = slot;
    editMax.value = slot.max_active;
    showEditModal.value = true;
}

async function saveSlots() {
    if (!editingSlot.value) return;
    const formData = new FormData();
    formData.append('max_active', String(editMax.value));

    try {
        const res = await fetch(`/api/admin/node-slots/${editingSlot.value.node_id}/slots`, {
            method: 'PATCH',
            body: formData,
        });
        const data = await res.json();
        if (data.success) {
            showEditModal.value = false;
            await fetchSlots();
        }
    } catch (e) {
    console.error('Failed to update slots:', e);
    }
}

onMounted(fetchSlots);
</script>
