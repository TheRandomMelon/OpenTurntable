<template>
    <div class="flex flex-col w-full h-screen">
        <MainBar />
        <div class="flex space-x-2">
            <button @click="initializeAudio" class="px-3 py-2 rounded bg-gray-800 cursor-pointer">Load Music</button>
        </div>
        <div class="flex-grow"></div>
        <NowPlayingBar />
    </div>
</template>

<script lang="ts" setup>
    const playback = usePlaybackStore();

    const initializeAudio = async() => {
        try {
            await playback.beginPlayback();
        } catch (err) {
            console.error(err);
        }
    }
</script>