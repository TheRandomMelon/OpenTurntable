<template>
	<NuxtLayout>
		<NuxtLoadingIndicator />
		<NuxtPage />
	</NuxtLayout>
</template>

<script lang="ts" setup>
	useHead({ titleTemplate: 'OpenTurntable' });

	console.log('[Frontend] Rendering application');
</script>
