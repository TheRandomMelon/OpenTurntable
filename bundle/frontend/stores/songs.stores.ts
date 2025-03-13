import { defineStore } from "pinia";
import { ChooseAndCreateSong, GetSongs } from "~/wailsjs/go/main/App";
import { database } from '~/wailsjs/go/models';

export const useSongsStore = defineStore("songs", {
    state: () => ({
        songs: null as database.Song[] | null
    }),
    actions: {
        async getAllSongs() {
            try {
                this.songs = await GetSongs();
            } catch (err) {
                console.error(err);
            }
        },

        async chooseSong() {
            await ChooseAndCreateSong();
            await this.getAllSongs();
        }
    }
})