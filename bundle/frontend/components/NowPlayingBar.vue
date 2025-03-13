<template>
    <div class="flex items-center justify-center w-full shadow-lg space-x-2 p-2" style="background: #252525;">
        <!-- Left hand side (track info) -->
        <div class="flex w-1/3 max-w-1/3 items-center" v-if="playback.filePath">
            <div class="w-[64px] min-w-[64px]">
                <img class="w-[64px] shadow rounded" draggable="false" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
            </div>
            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                <span class="font-bold truncate">{{ playback.metadata?.title ? playback.metadata.title : "Unknown Title" }}</span>
                <span class="truncate">{{ playback.metadata?.artist ? playback.metadata.artist : "Unknown Artist" }}</span>
            </div>
        </div>
        <div class="flex w-1/3 max-w-1/3 items-center" v-else>
            <div class="w-[64px] min-w-[64px]">
                <img class="w-[64px] shadow rounded" draggable="false" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
            </div>
            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                <span class="text-[#aaaaaa] italic">Not Playing</span>
            </div>
        </div>
        
        <!-- Middle (playback controls) -->
        <div class="flex flex-col space-y-2 justify-center items-center w-1/3 max-w-1/3">
            <div class="flex space-x-4 w-full">
                <span>{{ playback.position ? SecondsToDuration(playback.position) : "0:00" }}</span>
                <input
                    type="range"
                    min="0"
                    :max="playback.duration ? playback.duration : 0"
                    step="0.1"
                    v-model="playback.position"
                    @input="playback.changePosition"
                    class="cursor-pointer w-full"
                    :disabled="playback.duration ? false : true"
                />
                <span>{{ playback.duration ? SecondsToDuration(playback.duration) : "0:00" }}</span>
            </div>
            <div>
                <fa :icon="playback.playing ? 'circle-pause' : 'circle-play'" class="text-white cursor-pointer text-3xl" @click="playback.togglePlayback()"></fa>
            </div>
        </div>

        <!-- End (volume/misc controls) -->
        <div class="flex space-x-2 w-1/3 max-w-1/3 justify-end mr-2">
            <fa icon="volume-high" class="text-[#767676] cursor-pointer" @click="playback.toggleMute()"></fa>
            <input
                type="range"
                min="-5"
                max="0"
                step="0.1"
                v-model="playback.volume"
                @input="playback.updateVolume"
                class="cursor-pointer"
            />
        </div>
	</div>
</template>

<script lang="ts" setup>
    const playback = usePlaybackStore();
    import defaultArtwork from '@/assets/img/default_artwork.png';
    import { SecondsToDuration } from '~/utils/format';
    
    setInterval(async () => {
        if (playback.playing) {
            let pos = await playback.getPosition();
            let dur = await playback.getDuration();

            if (pos >= dur) {
                playback.playing = false;
            }
        }
    }, 10);
</script>