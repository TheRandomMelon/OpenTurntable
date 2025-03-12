<template>
    <div class="flex items-center w-full shadow-lg space-x-2 p-2" style="background: #252525;">
		<div>
            <img class="w-[64px] shadow rounded" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
        </div>
        <div class="flex ml-2 flex-col">
            <span class="font-bold">{{ playback.metadata ? playback.metadata.title : "Unknown Title" }}</span>
            <span class="">{{ playback.metadata?.artist ? playback.metadata.artist : "Unknown Artist" }}</span>
        </div>
        <div class="flex-grow"></div>
        <div class="flex flex-col space-y-2 justify-center items-center">
            <div class="flex space-x-4">
                <span>{{ playback.position ? SecondsToDuration(playback.position) : "0:00" }}</span>
                <input
                    type="range"
                    min="0"
                    :max="playback.duration ? playback.duration : 0"
                    step="0.1"
                    v-model="playback.position"
                    @input="playback.changePosition"
                    class="cursor-pointer w-64"
                    :disabled="playback.duration ? false : true"
                />
                <span>{{ playback.duration ? SecondsToDuration(playback.duration) : "0:00" }}</span>
            </div>
            <div>
                <fa :icon="playback.playing ? 'circle-pause' : 'circle-play'" class="text-white cursor-pointer text-3xl" @click="playback.togglePlayback()"></fa>
            </div>
        </div>
        <div class="flex-grow"></div>
        <div class="flex space-x-2">
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