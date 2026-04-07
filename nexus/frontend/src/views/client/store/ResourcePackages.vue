<template>
    <LayoutDashboard>
        <div class="p-6">
            <!-- Header Section -->
            <div class="flex items-center justify-between mb-6">
                <div>
                    <h1 class="text-2xl font-bold text-gray-100">{{ $t('resource_packages.title') }}</h1>
                    <p class="text-gray-400 mt-1">{{ $t('resource_packages.subtitle') }}</p>
                </div>
                <div class="px-4 py-2 bg-[#1A1825] rounded-lg border border-[#2a2a3f]/30">
                    <span class="text-sm text-gray-400">{{ $t('resource_packages.balance') }}:</span>
                    <span class="ml-2 text-lg font-bold text-green-400">{{ balance }}</span>
                    <span class="ml-1 text-sm text-gray-400">{{ $t('resource_packages.coins') }}</span>
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

            <!-- Resource Stats -->
            <div class="grid gap-4 sm:grid-cols-4 mb-6">
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-2xl font-bold text-blue-400">{{ resources.extra_memory }} MB</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('resource_packages.extra_memory') }}</div>
                </div>
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-2xl font-bold text-purple-400">{{ resources.extra_disk }} MB</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('resource_packages.extra_disk') }}</div>
                </div>
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-2xl font-bold text-orange-400">{{ resources.extra_cpu }}%</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('resource_packages.extra_cpu') }}</div>
                </div>
                <div class="bg-[#1A1825] rounded-lg p-4 border border-[#2a2a3f]/30 text-center">
                    <div class="text-2xl font-bold text-green-400">{{ resources.extra_servers }}</div>
                    <div class="text-sm text-gray-400 mt-1">{{ $t('resource_packages.extra_servers') }}</div>
                </div>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                <div v-for="n in 3" :key="n" class="bg-[#1A1825] rounded-xl p-6 border border-[#2a2a3f]/30 animate-pulse">
                    <div class="h-6 bg-gray-700 rounded w-2/3 mb-4"></div>
                    <div class="h-4 bg-gray-700 rounded w-full mb-2"></div>
                    <div class="h-4 bg-gray-700 rounded w-1/2 mb-4"></div>
                    <div class="h-10 bg-gray-700 rounded"></div>
                </div>
            </div>

            <!-- Packages Grid -->
            <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
                <div
                    v-for="pkg in packages"
                    :key="pkg.id"
                    class="bg-[#1A1825] rounded-xl p-6 border border-[#2a2a3f]/30 hover:border-indigo-500/50 transition-all duration-300 hover:shadow-lg hover:shadow-indigo-500/10"
                >
                    <h3 class="text-lg font-bold text-gray-100 mb-2">{{ pkg.name }}</h3>
                    <p class="text-gray-400 text-sm mb-4">{{ pkg.description }}</p>

                    <!-- Resource Details -->
                    <div class="space-y-2 mb-6">
                        <div v-if="pkg.memory > 0" class="flex items-center justify-between text-sm">
                            <span class="text-gray-400">{{ $t('resource_packages.memory') }}</span>
                            <span class="text-blue-400 font-medium">+{{ pkg.memory }} MB</span>
                        </div>
                        <div v-if="pkg.disk > 0" class="flex items-center justify-between text-sm">
                            <span class="text-gray-400">{{ $t('resource_packages.disk') }}</span>
                            <span class="text-purple-400 font-medium">+{{ pkg.disk }} MB</span>
                        </div>
                        <div v-if="pkg.cpu > 0" class="flex items-center justify-between text-sm">
                            <span class="text-gray-400">{{ $t('resource_packages.cpu') }}</span>
                            <span class="text-orange-400 font-medium">+{{ pkg.cpu }}%</span>
                        </div>
                    </div>

                    <!-- Price & Buy Button -->
                    <div class="flex items-center justify-between">
                        <div class="flex items-center gap-2">
                            <span class="text-2xl font-bold text-green-400">{{ pkg.coin_cost }}</span>
                            <span class="text-sm text-gray-400">{{ $t('resource_packages.coins') }}</span>
                        </div>
                        <button
                            @click="handlePurchase(pkg)"
                            :disabled="isPurchasing || balance < pkg.coin_cost"
                            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-500 disabled:bg-gray-700 disabled:text-gray-500 text-white text-sm font-medium rounded-lg transition-all duration-200"
                        >
                            {{ isPurchasing ? $t('resource_packages.processing') : $t('resource_packages.purchase') }}
                        </button>
                    </div>

                    <!-- Can't afford message -->
                    <p v-if="balance < pkg.coin_cost" class="text-xs text-red-400 mt-2">
                        {{ $t('resource_packages.not_enough_coins') }}
                    </p>
                </div>
            </div>

            <!-- Empty State -->
            <div v-if="!loading && packages.length === 0" class="text-center py-12">
                <div class="w-24 h-24 mx-auto mb-4 text-gray-400">
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                            d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                </div>
                <h3 class="text-lg font-semibold text-gray-200 mb-2">{{ $t('resource_packages.no_packages') }}</h3>
                <p class="text-gray-400">{{ $t('resource_packages.no_packages_desc') }}</p>
            </div>
        </div>
    </LayoutDashboard>
</template>

<script setup lang="ts">
import LayoutDashboard from '@/components/client/LayoutDashboard.vue';
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import Session from '@/mythicaldash/Session';
import { MythicalDOM } from '@/mythicaldash/MythicalDOM';

MythicalDOM.setPageTitle(MythicalDOM.getTranslation('resource_packages.title'));
const router = useRouter();

interface ResourcePackage {
    id: string;
    name: string;
    description: string;
    memory: number;
    disk: number;
    cpu: number;
    coin_cost: number;
}

interface UserResources {
    extra_memory: number;
    extra_disk: number;
    extra_cpu: number;
    extra_servers: number;
}

const packages = ref<ResourcePackage[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const success = ref<string | null>(null);
const isPurchasing = ref(false);
const balance = ref(Session.getInfoInt('credits') ?? 0);
const resources = ref<UserResources>({
    extra_memory: 0,
    extra_disk: 0,
    extra_cpu: 0,
    extra_servers: 0,
});

const fetchPackages = async () => {
    try {
        const response = await fetch('/api/user/resource-packages');
        if (!response.ok) throw new Error('Failed to fetch packages');
        const data = await response.json();
        if (data.success) {
            packages.value = data.packages;
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Failed to fetch packages';
    }
};

const fetchResources = async () => {
    try {
        const response = await fetch('/api/user/resource-packages/my-resources');
        if (!response.ok) return;
        const data = await response.json();
        if (data.success) {
            resources.value = data.resources;
        }
    } catch {
        // Ignore
    }
};

const handlePurchase = async (pkg: ResourcePackage) => {
    if (isPurchasing.value || balance.value < pkg.coin_cost) return;
    isPurchasing.value = true;
    error.value = null;
    success.value = null;

    try {
        const response = await fetch('/api/user/resource-packages/purchase', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ packageId: pkg.id }),
        });
        const data = await response.json();

        if (data.success) {
            success.value = `${pkg.name} purchased! Balance: ${data.newBalance}`;
            balance.value = data.newBalance;
            await fetchResources();
        } else {
            error.value = data.message || 'Purchase failed';
        }
    } catch (err) {
        error.value = err instanceof Error ? err.message : 'Purchase failed';
    } finally {
        isPurchasing.value = false;
    }
};

onMounted(() => {
    Promise.all([fetchPackages(), fetchResources()]).finally(() => {
        loading.value = false;
    });
});
</script>
