import { defineStore } from "pinia";
import { ChooseAndCreateSong, GetSongs, ImportSongsFromDirectory } from "~/wailsjs/go/main/App";
import { database } from '~/wailsjs/go/models';

export const useSongsStore = defineStore("songs", {
    state: () => ({
        songs: null as database.Song[] | null,
        importing: false as boolean
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
        },

        async importSongsFromDirectory() {
            await ImportSongsFromDirectory();
            await this.getAllSongs();
        }
    }
})