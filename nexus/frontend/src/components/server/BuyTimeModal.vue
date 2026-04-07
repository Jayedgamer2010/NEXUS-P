<template>
    <Teleport to="body">
        <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4" @keydown.self.esc="$emit('close')">
    <div class="fixed inset-0 bg-black/60 backdrop-blur-sm"></div>
    <div class="relative w-full max-w-2xl rounded-2xl border border-gray-800 bg-[#0d0d17] p-6 shadow-2xl">
        <!-- Header -->
        <div class="mb-6 flex items-center justify-between">
            <h2 class="text-xl font-bold text-white">Buy Server Time</h2>
            <div class="flex items-center gap-2 rounded-lg bg-gray-900 px-3 py-1.5 text-sm">
                <Coins class="h-4 w-4 text-yellow-400" />
                <span class="text-yellow-400 font-semibold">{{ userCoins }}</span>
                <span class="text-gray-500">coins</span>
            </div>
        </div>

                <!-- Warning -->
                <div class="mb-4 rounded-lg border border-yellow-900/50 bg-yellow-900/20 px-4 py-2 text-sm text-yellow-400">
                    Coins will be charged immediately when joining queue. Full refund if you leave before server starts.
                </div>

                <!-- Packages Grid -->
                <div class="mb-6 grid grid-cols-2 gap-4 md:grid-cols-4">
                    <div v-for="pkg in packages" :key="pkg.id" @click="selectPackage(pkg)"
                        :class="[
                            'relative cursor-pointer rounded-xl border-2 p-4 text-center transition-all',
                            selectedPackage?.id === pkg.id
                                ? (pkg.affordable ? 'border-purple-500 bg-purple-500/10 shadow-lg shadow-purple-500/20' : 'border-gray-600 bg-gray-800/50')
                                : 'border-gray-800 bg-gray-900/50 hover:border-gray-700'
                        ]">
                        <div :class="!pkg.affordable && 'opacity-50'">
                            <h3 class="mb-1 font-bold text-white">{{ pkg.name }}</h3>
                            <div class="mb-2 flex items-center justify-center gap-1">
                                <Clock class="h-4 w-4 text-purple-400" />
                                <span class="text-2xl font-bold text-purple-400">{{ pkg.minutes }}</span>
                                <span class="text-xs text-gray-400">min</span>
                            </div>
                            <div class="flex items-center justify-center gap-1">
                                <Coins class="h-3.5 w-3.5 text-yellow-500" />
                                <span class="text-sm font-semibold" :class="pkg.affordable ? 'text-yellow-400' : 'text-gray-500'">
                                    {{ pkg.coin_cost }} coins
                                </span>
                            </div>
                            <p class="mt-1 text-xs text-gray-500">{{ (pkg.minutes / pkg.coin_cost).toFixed(2) }} min/coin</p>
                        </div>
                        <!-- Not enough coins overlay -->
                        <div v-if="!pkg.affordable"
                            class="absolute inset-0 flex items-center justify-center rounded-xl bg-[#0d0d17]/80">
                            <p class="text-xs font-medium text-gray-400">Not enough coins</p>
                        </div>
                    </div>
                </div>

                <!-- Confirm -->
                <div class="flex gap-3">
                    <button @click="$emit('close')"
                        class="flex-1 rounded-lg border border-gray-700 bg-gray-900 px-4 py-2.5 text-sm font-medium text-gray-300 transition hover:bg-gray-800">
                        Cancel
                    </button>
                    <button @click="confirm" :disabled="!selectedPackage"
                        class="flex-1 rounded-lg bg-gradient-to-r from-purple-600 to-purple-500 px-4 py-2.5 text-sm font-medium text-white transition disabled:cursor-not-allowed disabled:opacity-50 hover:from-purple-700 hover:to-purple-600">
                        Join Queue & Buy
                    </button>
                </div>
            </div>
        </div>
    </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { Clock, Coins } from 'lucide-vue-next';

interface TimePackage {
    id: string;
    name: string;
    minutes: number;
    coin_cost: number;
    affordable: boolean;
}

const props = defineProps<{
    show: boolean;
    userCoins: number;
    packages: TimePackage[];
}>();

const emit = defineEmits<{
    close: [];
    confirm: [packageId: string];
}>();

const selectedPackage = ref<TimePackage | null>(null);

function selectPackage(pkg: TimePackage) {
    if (pkg.affordable) {
        selectedPackage.value = pkg;
    }
}

function confirm() {
    if (selectedPackage.value) {
        emit('confirm', selectedPackage.value.id);
        selectedPackage.value = null;
    }
}
</script>
