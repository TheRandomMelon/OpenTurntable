<template>
    <MainBar />
    <div class="flex space-x-2">
        <button @click="initializeAudio" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">Load Music</button>
        <button @click="pause" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">{{ state.isPlaying ? "Pause" : "Play" }}</button>
        <input
            type="range"
            min="-5"
            max="0"
            step="0.1"
            v-model="state.volume"
            @input="updateVolume"
        />
    </div>
</template>

<script lang="ts" setup>
    const state = reactive({
        isPlaying: false,
        volume: 0
    });

    import { PlayMusic, PauseMusic, SetVolume } from '../wailsjs/go/main/App.js';

    const initializeAudio = async() => {
        try {
            await PlayMusic();
            state.isPlaying = true;
        } catch (err) {
            console.error(err);
        }
    }

    const pause = async() => {
        await PauseMusic();
        state.isPlaying = !state.isPlaying;
    }

    const updateVolume = (event: Event) => {
        const target = event.target as HTMLInputElement;
        console.log(parseFloat(target.value));
        SetVolume(parseFloat(target.value));
    };
</script>