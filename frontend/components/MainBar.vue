<template>
    <div class="flex items-center w-full shadow-lg space-x-2 p-2" style="background: #252525;">
		<div class="flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing" style="background: #151515;">
			<fa icon="book" size="lg"></fa>
			<span class="font-bold">Library</span>
		</div>
		<div class="flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing">
			<fa icon="circle-play" size="lg"></fa>
			<span class="font-bold">Now Playing</span>
		</div>
		<div class="flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing">
			<fa icon="file-audio" size="lg"></fa>
			<span class="font-bold">Playlists</span>
		</div>
		<div class="flex-grow"></div>
		<div class="flex items-center justify-center space-x-2 px-4 py-2 rounded-md cursor-pointer thing">
			<fa icon="bars" size="lg"></fa>
		</div>
	</div>
</template>

<style scoped>
	.thing {
		transition: all 200ms ease;
	}

	.thing:hover {
		background: #151515;
	}
</style>