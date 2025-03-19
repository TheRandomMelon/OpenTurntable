<template>
    <div>
        <fa icon="caret-up" size="lg" v-if="ascending"></fa>
        <fa icon="caret-down" size="lg" v-else></fa>
    </div>
</template>

<script lang="ts" setup>
    import { toRefs } from 'vue'

    const props = defineProps({
        ascending: {
            type: Boolean,
            required: true
        }
    })
    const { ascending } = toRefs(props);
</script>