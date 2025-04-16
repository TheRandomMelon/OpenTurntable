<template>
    <div class="flex-1 overflow-hidden">
        <NewPlaylistDialog v-if="state.newPlaylistDialog" @close="closeDialog()" />
        <h1>Playlist parameter dump lmao</h1>
        <button class="px-3 py-2 rounded bg-gray-700 disabled:opacity-50 cursor-pointer" @click="state.newPlaylistDialog = !state.newPlaylistDialog;">
            Activate New Playlist Dialog
        </button>

        <div class="flex flex-row space-x-3" v-if="!state.isLoading && playlist.playlists" v-for="p in playlist.playlists">
            <p>{{ p.ID }}</p>
            <p>{{ p.Name }}</p>
            <p>{{ p.Description }}</p>
            <p>{{ p.Picture }}</p>
        </div>
        <div v-else-if="!state.isLoading">
            No playlists found
        </div>
        <div v-else>
            Loading...
        </div>
    </div>
</template>

<script setup lang="ts">
    import NewPlaylistDialog from '~/components/NewPlaylistDialog.vue';
    import { usePlaylistsStore } from '~/stores/playlists.stores';

    const playlist = usePlaylistsStore();

    const state = reactive({
        isLoading: true,
        newPlaylistDialog: false
    })

    onMounted(async () => {
        await playlist.loadPlaylists();
        state.isLoading = false;
    });

    const closeDialog = async () => {
        state.newPlaylistDialog = false;
        await playlist.loadPlaylists();
    }
</script>