<template>
    <div class="flex flex-col w-full h-screen">
        <MainBar />
        <div class="p-8">
            <div class="flex space-x-2">
                <button @click="songs.chooseSong" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">Add Song</button>
            </div>
            <div v-if="!state.isLoading && songs.songs !== null">
                <table class="w-full">
                    <thead>
                        <tr>
                            <th class="play-btn"></th>
                            <th>ID</th>
                            <th>Path</th>
                            <th>Title</th>
                            <th>Artist ID</th>
                            <th>Album ID</th>
                            <th>Comment</th>
                            <th>Genre</th>
                            <th>Year</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="song in songs.songs">
                            <td class="play-btn"><fa icon="circle-play" class="text-white cursor-pointer text-xl" @click="playback.beginPlayback(song.Path)"></fa></td>
                            <td>{{ song.ID }}</td>
                            <td>{{ song.Path }}</td>
                            <td>{{ song.Title }}</td>
                            <td>{{ song.Artist_ID }}</td>
                            <td>{{ song.Album_ID }}</td>
                            <td>{{ song.Comment }}</td>
                            <td>{{ song.Genre }}</td>
                            <td>{{ song.Year }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div v-else-if="!state.isLoading">
                <i>You don't currently have any songs in your library</i>
            </div>
            <div v-else>
                <i>Loading library...</i>
            </div>
        </div>
        <div class="flex-grow"></div>
        <NowPlayingBar />
    </div>
</template>

<script lang="ts" setup>
    import { database } from '~/wailsjs/go/models';

    const playback = usePlaybackStore();
    const songs = useSongsStore();

    const state = reactive({
        isLoading: true
    })

    onMounted(async () => {
        await songs.getAllSongs();
        state.isLoading = false;
    });
</script>