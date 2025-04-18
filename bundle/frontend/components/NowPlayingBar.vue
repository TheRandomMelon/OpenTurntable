<template>
    <div class="flex items-center justify-center w-full shadow-lg space-x-2 p-2 bg-[#252525]">
        <!-- Left hand side (track info) -->
        <div class="flex w-1/3 max-w-1/3 items-center" v-if="playback.filePath">
            <div class="w-[64px] min-w-[64px]">
                <img class="w-[64px] shadow rounded" draggable="false" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
            </div>
            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                <span class="font-bold truncate">{{ playback.metadata?.title ? playback.metadata.title : $t('defaults.title') }}</span>
                <span class="truncate">{{ playback.metadata?.artist ? playback.metadata.artist : $t('defaults.artist') }}</span>
            </div>
        </div>
        <div class="flex w-1/3 max-w-1/3 items-center" v-else>
            <div class="w-[64px] min-w-[64px]">
                <img class="w-[64px] shadow rounded" draggable="false" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
            </div>
            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                <span class="text-[#aaaaaa] italic">{{ $t('general.not_playing') }}</span>
            </div>
        </div>
        
        <!-- Middle (playback controls) -->
        <div class="flex flex-col space-y-1 justify-center items-center w-1/3 max-w-1/3">
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
            <div class="flex items-center space-x-6">
                <fa icon="shuffle" :class="playback.shuffle ? 'text-blue-600 cursor-pointer text-xl' : 'text-[#bbb] cursor-pointer text-xl'" @click="playback.toggleShuffle()"></fa>
                <fa icon="backward-step" class="text-[#bbb] cursor-pointer text-xl" @click="playback.queueStep(false, true)"></fa>
                <fa :icon="playback.playing ? 'circle-pause' : 'circle-play'" class="text-white cursor-pointer text-3xl" @click="playback.togglePlayback()"></fa>
                <fa icon="forward-step" class="text-[#bbb] cursor-pointer text-xl" @click="playback.queueStep(true, true)"></fa>
                <div class="relative inline-block">
                    <fa icon="repeat" :class="playback.repeat === RepeatType.Off ? 'text-[#bbb] cursor-pointer text-xl' : 'text-blue-600 cursor-pointer text-xl'" @click="playback.cycleRepeat()"></fa>
                    <span v-if="playback.repeat === RepeatType.RepeatOne" class="absolute -top-1 -right-2 bg-blue-600 text-white text-xs font-bold w-4 h-4 flex items-center justify-center rounded-full">1</span>
                </div>
            </div>
        </div>

        <!-- End (volume/misc controls) -->
        <div class="flex items-center space-x-3 w-1/3 max-w-1/3 justify-end mr-2">
            <div class="relative">
                <fa icon="gauge" class="text-[#bbb] cursor-pointer" @click="state.speedDropdown = !state.speedDropdown"></fa>

                <div class="absolute bottom-full mb-2 bg-[#333] shadow-lg rounded-md p-2 w-24" v-if="state.speedDropdown">
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 0.25 ? 'bg-[#444]' : ''}`" @click="changeSpeed(0.25)">0.25x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 0.5 ? 'bg-[#444]' : ''}`" @click="changeSpeed(0.5)">0.5x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 0.75 ? 'bg-[#444]' : ''}`" @click="changeSpeed(0.75)">0.75x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 1 ? 'bg-[#444]' : ''}`" @click="changeSpeed(1)">1x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 1.25 ? 'bg-[#444]' : ''}`" @click="changeSpeed(1.25)">1.25x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 1.5 ? 'bg-[#444]' : ''}`" @click="changeSpeed(1.5)">1.5x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 1.75 ? 'bg-[#444]' : ''}`" @click="changeSpeed(1.75)">1.75x</div>
                    <div :class="`cursor-pointer px-3 py-1 rounded-md hover:bg-[#444] ${playback.speed === 2 ? 'bg-[#444]' : ''}`" @click="changeSpeed(2)">2x</div>
                </div>
            </div>

            <fa icon="volume-high" class="text-[#bbb] cursor-pointer" @click="playback.toggleMute()"></fa>
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
    import { RepeatType } from '~/stores/playback.stores';
    import { SecondsToDuration } from '~/utils/format';

    const state = reactive({
        speedDropdown: false
    });

    const changeSpeed = async (speed: number) => {
        state.speedDropdown = false;
        await playback.setSpeed(speed);
    }
    
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