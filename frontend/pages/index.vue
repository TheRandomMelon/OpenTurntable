<template>
    <MainBar />
    <button @click="test()">Test</button>
</template>

<script lang="ts" setup>
    import { Greet } from '../wailsjs/go/main/App.js';

    const test = async() => {
        Greet("Hello").then((thing) => {
            console.log(thing);
        })
    }
</script>