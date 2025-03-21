import { defineStore } from "pinia";
import { ChooseAndCreateSong, GetSongs, ImportSongsFromDirectory, GetSongsWithDetails } from "~/wailsjs/go/main/App";
import { database } from '~/wailsjs/go/models';

export const useSongsStore = defineStore("songs", {
    state: () => ({
        songs: null as database.SongWithDetails[] | null,
        importing: false as boolean,
        arrangement: {
            key: "id",
            asc: true
        }
    }),
    actions: {
        async getAllSongs() {
            try {
                this.songs = await GetSongsWithDetails();
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
        },

        async ascOrReset() {
            if (this.arrangement.asc) {
                this.arrangement.asc = false;
            } else {
                this.arrangement.key = "id";
                this.arrangement.asc = true;
            }
        },

        async rearrangeSongs(by: string) {
            // Figure out exact rearrangement
            switch (by) {
                case "title":
                    if (this.arrangement.key === "title") {
                        if (this.arrangement.asc) {
                            this.arrangement.asc = false;
                        } else {
                            this.arrangement.key = "artist";
                            this.arrangement.asc = true;
                        }
                    } else {
                        this.arrangement.key = "title";
                        this.arrangement.asc = true;
                    }
                    break;
                case "artist":
                    if (this.arrangement.key === "artist") {
                        await this.ascOrReset();
                    } else {
                        this.arrangement.key = "artist";
                        this.arrangement.asc = true;
                    }
                    break;
                case "album":
                    if (this.arrangement.key === 'album') {
                        await this.ascOrReset();
                    } else {
                        this.arrangement.key = 'album';
                        this.arrangement.asc = true;
                    }
                    break;
                case "genre":
                    if (this.arrangement.key === 'genre') {
                        await this.ascOrReset();
                    } else {
                        this.arrangement.key = 'genre';
                        this.arrangement.asc = true;
                    }
                    break;
                case "year":
                    if (this.arrangement.key === 'year') {
                        await this.ascOrReset();
                    } else {
                        this.arrangement.key = 'year';
                        this.arrangement.asc = true;
                    }
                    break;
                case "id":
                    this.arrangement.key = 'id';
                    this.arrangement.asc = !this.arrangement.asc;
                    break;
                default:
                    this.arrangement.key = "id";
                    this.arrangement.asc = true;
                    break;
            }

            // Execute resorting
            switch (this.arrangement.key) {
                case "title":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * a.Title.localeCompare(b.Title));
                    break;
                case "artist":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * a.ArtistName.String.localeCompare(b.ArtistName.String));
                    break;
                case "album":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * a.AlbumName.String.localeCompare(b.AlbumName.String));
                    break;
                case "genre":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * a.Genre.String.localeCompare(b.Genre.String));
                    break;
                case "year":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * (a.Year > b.Year ? 1 : -1));
                    break;
                case "id":
                    this.songs?.sort((a, b) => (this.arrangement.asc ? 1 : -1) * (a.ID > b.ID ? 1 : -1));
                    break;
            }

            return this.arrangement;
        },

        getQueue<T>(arr: T[], startIndex: number): T[] {
            if (arr.length === 0) return [];
            const adjustedIndex = startIndex % arr.length;
            const safeIndex = adjustedIndex >= 0 ? adjustedIndex : adjustedIndex + arr.length;
            return [...arr.slice(safeIndex), ...arr.slice(0, safeIndex)];
        }
    }
})