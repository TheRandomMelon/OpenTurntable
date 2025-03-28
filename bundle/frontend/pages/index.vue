<template>
   <div class="flex-1 overflow-hidden">
        <ImportDialog v-if="songs.importing" :currentlyImporting="state.currentlyImporting" />
        <div class="p-4 h-full flex flex-col gap-4">
            <div class="flex space-x-2">
                <button @click="songs.chooseSong" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">
                    <fa icon="square-plus" class="mr-1"></fa> {{ $t("library.add_song") }}
                </button>
                <button @click="songs.importSongsFromDirectory" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">
                    <fa icon="folder-plus" class="mr-1"></fa> {{ $t("library.import_songs") }}
                </button>
            </div>
            <div class="flex-1 overflow-auto" v-if="!state.isLoading && songs.songs !== null && !songs.importing">
                <table class="w-full">
                    <thead>
                        <tr>
                            <th class="cursor-pointer" @click="rearrangeSongs('id')">
                                <div class="flex flex-row items-center space-x-2">
                                    <span>#</span>
                                    <SortIcon :ascending="songs.arrangement.asc" v-if="songs.arrangement.key === 'id'" />
                                </div>
                            </th>
                            <th class="cursor-pointer" @click="rearrangeSongs(state.sortTitle ? 'title' : 'artist')">
                                <div class="flex flex-row items-center space-x-2">
                                    <span>{{ songs.arrangement.key === 'artist' ? $t('general.artist') : $t('general.title')  }}</span>
                                    <SortIcon :ascending="songs.arrangement.asc" v-if="songs.arrangement.key === 'title' || songs.arrangement.key === 'artist'" />
                                </div>
                            </th>
                            <th class="cursor-pointer" @click="rearrangeSongs('album')">
                                <div class="flex flex-row items-center space-x-2">
                                    <span>{{ $t('general.album') }}</span>
                                    <SortIcon :ascending="songs.arrangement.asc" v-if="songs.arrangement.key === 'album'" />
                                </div>
                            </th>
                            <th class="cursor-pointer" @click="rearrangeSongs('genre')">
                                <div class="flex flex-row items-center space-x-2">
                                    <span>{{ $t('general.genre') }}</span>
                                    <SortIcon :ascending="songs.arrangement.asc" v-if="songs.arrangement.key === 'genre'" />
                                </div>
                            </th>
                            <th class="cursor-pointer" @click="rearrangeSongs('year')">
                                <div class="flex flex-row items-center space-x-2">
                                    <span>{{ $t('general.year') }}</span>
                                    <SortIcon :ascending="songs.arrangement.asc" v-if="songs.arrangement.key === 'year'" />
                                </div>
                            </th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="song in songs.songs">
                            <td class="play-btn">
                                <span class="id-label">{{ song.ID }}</span>
                                <div class="play-button-icon">
                                    <fa icon="circle-play" class="text-white cursor-pointer text-xl" @click="playback.beginPlayback(song, PlaybackSourceType.Library)"></fa>
                                </div>
                            </td>
                            <td>
                                <div class="flex flex-row space-x-3 items-center">
                                    <div class="w-[42px] min-w-[42px]">
                                        <img class="w-[42px] shadow rounded" draggable="false" :src="song.AlbumArt.String ? song.AlbumArt.String : defaultArtwork" />
                                    </div>
                                    <div class="flex flex-col">
                                        <b>{{ song.Title }}</b>
                                        <span>{{ song.ArtistName.String }}</span>
                                    </div>
                                </div>
                            </td>
                            <td>{{ song.AlbumName.String }}</td>
                            <td>{{ song.Genre.String }}</td>
                            <td>{{ song.Year.String == "0" ? "" : song.Year.String }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div v-else-if="!state.isLoading && !songs.importing">
                <i>{{ $t('library.no_songs') }}</i>
            </div>
            <div v-else>
                <i>{{ $t('library.loading') }}</i>
            </div>
        </div>
    </div>
</template>

<style lang="css" scoped>
table th {
    position: sticky;
    top: 0;
    background: var(--component);
}

td .play-button-icon {
    display: none;
}

td:hover .play-button-icon {
    display: initial;
}

td .id-label {
    display: initial;
}

td:hover .id-label {
    display: none;
}
</style>

<script lang="ts" setup>
    import { database } from '~/wailsjs/go/models';
    import { EventsOn } from '~/wailsjs/runtime';
    import defaultArtwork from '@/assets/img/default_artwork.png';
    import { PlaybackSourceType } from '~/stores/playback.stores';
    import ImportDialog from '~/components/ImportDialog.vue';

    const playback = usePlaybackStore();
    const songs = useSongsStore();

    const state = reactive({
        isLoading: true,
        currentlyImporting: "",
        sortTitle: true,
        importDialog: false
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
        
        EventsOn("playbackComplete", async () => {
			playback.queueStep(true);
		});
    });

    const rearrangeSongs = async (by: string) => {
        let arrangement = await songs.rearrangeSongs(by);
        console.log(arrangement);

        if (arrangement.key === "title") {
            state.sortTitle = true;
        } else if (arrangement.key === "artist") {
            state.sortTitle = false;
        } else {
            state.sortTitle = true;
        }
    }
</script>