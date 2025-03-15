<template>
    <div class="flex flex-col w-full h-screen">
        <MainBar />
        <div class="flex-1 overflow-y-auto">
            <div class="p-8">
                <div class="flex space-x-2">
                    <button @click="songs.chooseSong" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">Add Song</button>
                    <button @click="songs.importSongsFromDirectory" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">Add Folder of Songs</button>
                </div>
                <div v-if="!state.isLoading && songs.songs !== null && !songs.importing">
                    <table class="w-full">
                        <thead>
                            <tr>
                                <th class="play-btn"></th>
                                <th>ID</th>
                                <th>Path</th>
                                <th>Title</th>
                                <th>Artist</th>
                                <th>Album</th>
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
                                <td>{{ song.ArtistName.String }}</td>
                                <td>{{ song.AlbumName.String }}</td>
                                <td>{{ song.Comment.String }}</td>
                                <td>{{ song.Genre.String }}</td>
                                <td>{{ song.Year.String == "0" ? "" : song.Year.String }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>
                <div v-else-if="!state.isLoading && !songs.importing">
                    <i>You don't currently have any songs in your library</i>
                </div>
                <div v-else-if="songs.importing">
                    <i>Importing songs...</i><br/>
                    <i>Currently working on {{ state.currentlyImporting }}</i>
                </div>
                <div v-else>
                    <i>Loading library...</i>
                </div>
            </div>
        </div>
        <NowPlayingBar />
    </div>
</template>

<script lang="ts" setup>
    import { database } from '~/wailsjs/go/models';
    import { EventsOn } from '~/wailsjs/runtime';

    const playback = usePlaybackStore();
    const songs = useSongsStore();

    const state = reactive({
        isLoading: true,
        currentlyImporting: ""
    })

    onMounted(async () => {
        await songs.getAllSongs();
        state.isLoading = false;

        EventsOn("toggleImporting", async () => {
			songs.importing = !songs.importing;
		});
        
        EventsOn("currentImportFileWorking", async (path: string) => {
			state.currentlyImporting = path;
		});
    });
</script>