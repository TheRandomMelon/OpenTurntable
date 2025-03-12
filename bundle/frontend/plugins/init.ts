// Code to run on initialization

export default defineNuxtPlugin(async nuxtApp => {
    nuxtApp.hook('app:mounted', async () => {
        console.log('[Frontend] Running init code');

        const { usePlaybackStore } = await import('~/stores/playback.stores'); // Import dynamically
        const playback = usePlaybackStore(nuxtApp.$pinia); // Pass the pinia instance

        await playback.reloadData();
    });
});