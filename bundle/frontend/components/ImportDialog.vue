<template>
	<div class="flex justify-center items-center min-h-screen z-1000">
		<div class="fixed inset-0 flex items-center justify-center bg-black/70 z-1000">
			<div class="bg-[#252525] p-6 rounded-lg shadow-lg w-96 max-w-96 z-1000">
				<h2 class="text-xl font-semibold">{{ $t('library.importing_songs') }}</h2>
				<p class="mt-2">{{ $t('library.currently_working') }}</p>
				<p class="mt-2 break-all max-w-92">{{ currentlyImporting }}</p>
			</div>
		</div>
	</div>
</template>

<script setup>
defineProps(['currentlyImporting']);
const emit = defineEmits(['close']);
</script>