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
                        {{ playback.metadata?.album ? playback.metadata.album : "Unknown Album" }} • {{ playback.metadata?.year }}
                    </p>
                </div>
            </div>
            <div class="flex-grow"></div>
            <div class="flex p-4 items-end">
                <div class="bg-[#252525] p-4 rounded albumart-shading w-86">
                    <h1 class="text-xl font-bold">Song Info</h1>
                    <div class="flex items-center">
                        <p class="font-bold mr-2">Year</p>
                        <div class="flex-grow"></div>
                        <p class="break-all text-right max-w-56">{{ playback.metadata?.year }}</p>
                    </div>
                    <div class="flex items-center">
                        <p class="font-bold mr-2">Genre</p>
                        <div class="flex-grow"></div>
                        <p class="break-all text-right max-w-56">{{ playback.metadata?.genre ? playback.metadata?.genre : 'Not Specified'}}</p>
                    </div>
                    <div class="flex items-start">
                        <p class="font-bold mr-2">Comment</p>
                        <div class="flex-grow"></div>
                        <p class="break-all text-right max-w-56">{{ playback.metadata?.comment ? playback.metadata?.comment : 'No Comment'}}</p>
                    </div>
                    <div class="flex items-start">
                        <p class="font-bold mr-2">Composer</p>
                        <div class="flex-grow"></div>
                        <p class="break-all text-right max-w-56">{{ playback.metadata?.composer ? playback.metadata?.composer : 'Not Specified'}}</p>
                    </div>
                    <div class="flex items-start">
                        <p class="font-bold mr-2">File Path</p>
                        <div class="flex-grow"></div>
                        <p class="break-all text-right max-w-56">{{ playback.filePath ? playback.filePath : 'Not Specified'}}</p>
                    </div>
                    <hr class="border-[#aaa] mt-3 mb-2"/>
                    <h1 class="text-xl font-bold mb-2">Next In Queue</h1>
                    <div class="flex">
                        <div class="flex items-center" v-if="nextInQueue">
                            <div class="w-[64px] min-w-[64px]">
                                <img class="w-[64px] shadow rounded" draggable="false" :src="nextInQueue.AlbumArt.String ? nextInQueue.AlbumArt.String : defaultArtwork" />
                            </div>
                            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                                <span class="font-bold truncate">{{ nextInQueue.Title ? nextInQueue.Title : "Unknown Title" }}</span>
                                <span class="truncate">{{ nextInQueue.ArtistName.String ? nextInQueue.ArtistName.String : "Unknown Artist" }}</span>
                            </div>
                        </div>
                        <div class="flex items-center" v-else>
                            <div class="w-[64px] min-w-[64px]">
                                <img class="w-[64px] shadow rounded" draggable="false" :src="defaultArtwork" />
                            </div>
                            <div class="flex ml-4 flex-col flex-1 overflow-hidden mr-4">
                                <span class="text-[#aaaaaa] italic">No more items in queue</span>
                            </div>
                        </div>
                    </div>
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

    const nextInQueue = await playback.getNextInQueue();
    console.log(nextInQueue);
</script>