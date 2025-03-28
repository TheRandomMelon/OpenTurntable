<template>
    <div class="flex items-center w-full shadow-lg space-x-2 p-2" style="background: #252525;">
		<div @click="navigateTo('/')" :class="`flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing ${route.path === '/' ? 'bg-[#151515]' : ''}`">
			<fa icon="book" size="lg"></fa>
			<span class="font-bold">{{ $t('general.library') }}</span>
		</div>
		<div @click="navigateTo('/nowplaying')" :class="`flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing ${route.path === '/nowplaying' ? 'bg-[#151515]' : ''}`">
			<fa icon="circle-play" size="lg"></fa>
			<span class="font-bold">{{ $t('general.now_playing') }}</span>
		</div>
		<div @click="navigateTo('/playlists')" :class="`flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing ${route.path === '/playlists' ? 'bg-[#151515]' : ''}`">
			<fa icon="file-audio" size="lg"></fa>
			<span class="font-bold">{{ $t('general.playlist', 2) }}</span>
		</div>
		<div class="flex-grow"></div>
		<div class="flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing">
			<fa icon="bars" size="lg"></fa>
		</div>
	</div>
</template>

<script lang="ts" setup>
	const route = useRoute();
</script>

<style scoped>
	.thing {
		transition: all 200ms ease;
	}

	.thing:hover {
		background: #151515;
	}
</style>