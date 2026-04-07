<template>
    <LayoutDashboard>
        <div class="mb-6 flex items-center justify-between">
            <div>
                <h1 class="text-2xl font-bold text-pink-400">Time Packages</h1>
                <p class="text-sm text-gray-400 mt-1">Manage server time packages and pricing</p>
            </div>
            <button @click="openCreateModal"
                class="bg-linear-to-r from-pink-500 to-violet-500 text-white px-4 py-2 rounded-lg transition-all duration-200 hover:opacity-80 flex items-center">
                <PlusIcon class="w-4 h-4 mr-2" />
                Create Package
            </button>
        </div>

        <div v-if="loading" class="flex justify-center items-center py-10">
            <LoaderCircle class="h-8 w-8 animate-spin text-pink-400" />
        </div>

        <div v-else class="space-y-4">
            <div v-for="pkg in packages" :key="pkg.id"
                class="rounded-xl border border-gray-800 bg-[#0d0d17] p-4">
                <div class="flex items-center justify-between">
                    <div class="flex items-center gap-6">
                        <div>
                            <h3 class="font-semibold text-white">{{ pkg.name }}</h3>
                            <p class="text-xs text-gray-500">ID: {{ pkg.id.slice(0, 8) }}...</p>
                        </div>
                        <div class="text-center">
                            <p class="text-2xl font-bold text-purple-400">{{ pkg.minutes }}</p>
                            <p class="text-xs text-gray-500">minutes</p>
                        </div>
                        <div class="text-center">
                            <p class="text-2xl font-bold text-yellow-400">{{ pkg.coin_cost }}</p>
                            <p class="text-xs text-gray-500">coins</p>
                        </div>
                        <div class="text-center">
                            <p class="text-lg font-semibold text-green-400">{{ (pkg.minutes / pkg.coin_cost).toFixed(2) }}</p>
                            <p class="text-xs text-gray-500">min/coin</p>
                        </div>
                    </div>
                    <div class="flex gap-2">
                        <button @click="openEditModal(pkg)"
                            class="rounded-lg border border-gray-700 px-3 py-1.5 text-sm text-gray-300 transition hover:bg-gray-800">
                            Edit
                        </button>
                        <button @click="deletePackage(pkg.id)"
                            class="rounded-lg bg-red-600/80 px-3 py-1.5 text-sm text-white transition hover:bg-red-700">
                            Delete
                        </button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Create/Edit Modal -->
        <div v-if="showModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
            <div class="fixed inset-0 bg-black/60" @click="showModal = false"></div>
            <div class="relative w-full max-w-md rounded-2xl border border-gray-800 bg-[#0d0d17] p-6 shadow-2xl">
                <h3 class="mb-4 text-lg font-bold text-white">
                    {{ editingId ? 'Edit Package' : 'Create Package' }}
                </h3>
                <div class="space-y-3">
                    <div>
                        <label class="mb-1 block text-sm text-gray-400">Name</label>
                        <input v-model="form.name" type="text" placeholder="e.g. Starter"
                            class="w-full rounded-lg border border-gray-700 bg-gray-900 px-3 py-2 text-white focus:border-purple-500 focus:outline-none" />
                    </div>
                    <div>
                        <label class="mb-1 block text-sm text-gray-400">Minutes</label>
                        <input v-model.number="form.minutes" type="number" min="1"
                            class="w-full rounded-lg border border-gray-700 bg-gray-900 px-3 py-2 text-white focus:border-purple-500 focus:outline-none" />
                    </div>
                    <div>
                        <label class="mb-1 block text-sm text-gray-400">Coin Cost</label>
                        <input v-model.number="form.coin_cost" type="number" min="1"
                            class="w-full rounded-lg border border-gray-700 bg-gray-900 px-3 py-2 text-white focus:border-purple-500 focus:outline-none" />
                    </div>
                    <div v-if="form.minutes && form.coin_cost" class="rounded-lg bg-gray-900 px-3 py-2 text-center">
                        <span class="text-sm text-gray-400">Value: </span>
                        <span class="font-semibold text-green-400">{{ (form.minutes / form.coin_cost).toFixed(2) }} min/coin</span>
                    </div>
                </div>
                <div class="mt-4 flex gap-3">
                    <button @click="showModal = false"
                        class="flex-1 rounded-lg border border-gray-700 bg-gray-900 px-4 py-2 text-sm text-gray-300 hover:bg-gray-800">
                        Cancel
                    </button>
                    <button @click="savePackage"
                        class="flex-1 rounded-lg bg-gradient-to-r from-purple-600 to-purple-500 px-4 py-2 text-sm font-medium text-white hover:from-purple-700 hover:to-purple-600">
                        {{ editingId ? 'Save' : 'Create' }}
                    </button>
                </div>
            </div>
        </div>
    </LayoutDashboard>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import LayoutDashboard from '@/components/admin/LayoutDashboard.vue';
import { PlusIcon, LoaderCircle } from 'lucide-vue-next';

interface TimePackage {
    id: string;
    name: string;
    minutes: number;
    coin_cost: number;
}

const loading = ref(true);
const packages = ref<TimePackage[]>([]);
const showModal = ref(false);
const editingId = ref<string | null>(null);
const form = ref({ name: '', minutes: 0, coin_cost: 0 });

async function fetchPackages() {
    try {
        const res = await fetch('/api/admin/time-packages');
        const data = await res.json();
        if (data.success) {
            packages.value = data.time_packages || [];
        }
    } catch (e) {
    } finally {
        loading.value = false;
    }
}

function openCreateModal() {
    editingId.value = null;
    form.value = { name: '', minutes: 5, coin_cost: 10 };
    showModal.value = true;
}

function openEditModal(pkg: TimePackage) {
    editingId.value = pkg.id;
    form.value = { name: pkg.name, minutes: pkg.minutes, coin_cost: pkg.coin_cost };
    showModal.value = true;
}

async function savePackage() {
    if (!form.value.name || form.value.minutes <= 0 || form.value.coin_cost <= 0) return;

    const formData = new FormData();
    formData.append('name', form.value.name);
    formData.append('minutes', String(form.value.minutes));
    formData.append('coin_cost', String(form.value.coin_cost));

    try {
        const url = editingId.value
            ? `/api/admin/time-packages/${editingId.value}/update`
            : '/api/admin/time-packages/create';
        const res = await fetch(url, { method: 'POST', body: formData });
        const data = await res.json();
        if (data.success) {
            showModal.value = false;
            await fetchPackages();
        }
    } catch (e) {
    }
}

async function deletePackage(id: string) {
    if (!confirm('Delete this time package?')) return;
    try {
        const res = await fetch(`/api/admin/time-packages/${id}/delete`, { method: 'POST' });
        const data = await res.json();
        if (data.success) {
            await fetchPackages();
        }
    } catch (e) {
    }
}

onMounted(fetchPackages);
</script>
