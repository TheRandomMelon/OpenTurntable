<template>
    <div :class="`flex-1 overflow-hidden bg-no-repeat bg-cover bg-center`" :style="`background-image: url('${playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork}');`" v-if="playback.metadata"> 
        <div class="flex items-end w-full h-full backdrop-blur-[100px] backdrop-brightness-70">
            <div class="flex p-4 space-x-4 items-end">
                <div>
                    <img class="w-42 albumart-shading" :src="playback.metadata?.albumArt ? playback.metadata.albumArt : defaultArtwork" />
                </div>
                <div>
                    <h1 class="text-3xl font-bold text-shading">{{ playback.metadata?.title }}</h1>
                    <p class="text-shading">
                        {{ playback.metadata?.artist ? playback.metadata.artist : "Unknown Artist" }}<br/>
                        {{ playback.metadata?.album ? playback.metadata.album : "Unknown Album" }}
                    </p>
                </div>
            </div>
        </div>
    </div>
    <div class="flex-1 overflow-hidden" v-else> 
        <h1>You aren't currently playing anything.</h1>
    </div>
</template>

<style lang="css" scoped>
    .albumart-shading {
        box-shadow: 0px 0px 7px #111;
    }

    .text-shading {
        text-shadow: 0px 0px 7px #111;
    }
</style>

<script setup lang="ts">
    const playback = usePlaybackStore();
    import defaultArtwork from '@/assets/img/default_artwork.png';
</script>