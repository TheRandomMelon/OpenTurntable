<template>
	<div class="flex justify-center items-center min-h-screen z-1000">
		<div class="fixed inset-0 flex items-center justify-center bg-black/70 z-1000">
			<div class="bg-[#252525] p-6 rounded-lg shadow-lg w-96 max-w-96 z-1000">
                <div class="flex flex-row items-center">
				    <h2 class="text-xl font-semibold">{{ $t('playlists.new') }}</h2>
                    <div class="flex-grow"></div>
                    <fa icon="x" @click="emit('close')" class="cursor-pointer"></fa>
                </div>
                <div class="" v-if="state.error">
                    <p class="text-red-500">{{ state.error }}</p>
                </div>
                <div class="flex flex-col mt-4 space-y-4">
                    <div class="flex flex-col">
                        <p class="font-bold mb-1">{{ $t('general.name') }}</p>
                        <input type="text" :placeholder="$t('general.name')" v-model="name" @input="validateName()" />
                    </div>
                    <div class="flex flex-col">
                        <p class="font-bold mb-1">{{ $t('general.description') }}</p>
                        <input type="text" :placeholder="$t('general.description')" v-model="description" />
                    </div>
                </div>
                <div class="mt-6">
                    <button @click="createPlaylist()" class="px-3 py-2 rounded bg-gray-700 disabled:opacity-50 cursor-pointer" :disabled="state.buttonDisabled">
                        <fa icon="square-plus" class="mr-1"></fa> {{ $t("playlists.new") }}
                    </button>
                </div>
			</div>
		</div>
	</div>
</template>

<style lang="css" scoped>
    input[type="text"] {
        background: #151515;
        padding: 8px;
        outline: none;
        border-radius: 5px;
    }
</style>

<script setup>
const emit = defineEmits(['close']);
const playlist = usePlaylistsStore();

const state = reactive({
    buttonDisabled: true,
    error: ""
});

const name = defineModel('name');
const description = defineModel('description');

const validateName = async () => {
    if (name.value && name.value !== "" && name.value !== null) {
        state.buttonDisabled = false;
    } else {
        state.buttonDisabled = true;
    }
}

const createPlaylist = async () => {
    state.buttonDisabled = true;
    state.error = "";

    if (!name.value || name.value == "" || name.value == null) {
        state.error = "Please specify a playlist name.";
        state.buttonDisabled = false;
        return;
    }

    try {
        await playlist.createPlaylist(name.value, description.value, "");

        name.value = "";
        description.value = "";
        state.buttonDisabled = false;
        emit("close");
    } catch (err) {
        console.error(err);
        state.error = "This should not occur. Please report this bug.";
    }
}
</script>