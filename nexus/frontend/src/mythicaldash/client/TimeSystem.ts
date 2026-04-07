export default {
    async getTimeStatus(serverUuid: string) {
        const response = await fetch(`/api/user/server/time/${serverUuid}/status`);
        return response.json();
    },

    async getQueueStatus(serverUuid: string) {
        const response = await fetch(`/api/user/server/time/${serverUuid}/queue-status`);
        return response.json();
    },

    async joinQueue(serverUuid: string, packageId?: string) {
        const formData = new FormData();
        if (packageId) formData.append('package_id', packageId);
        const response = await fetch(`/api/user/server/time/${serverUuid}/queue`, {
            method: 'POST',
            body: formData,
        });
        return response.json();
    },

    async leaveQueue(serverUuid: string) {
        const response = await fetch(`/api/user/server/time/${serverUuid}/leave-queue`, {
            method: 'POST',
        });
        return response.json();
    },

    async getTimePackages() {
        const response = await fetch('/api/admin/time-packages');
        return response.json();
    },

    async getNodeSlots() {
        const response = await fetch('/api/admin/node-slots');
        return response.json();
    },
};
