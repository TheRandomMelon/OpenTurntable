import { defineStore } from "pinia";
import { ChooseAndCreateSong, GetSongs, ImportSongsFromDirectory, GetSongsWithDetails, GetPlaylistWithSongs, GetPlaylists, CreatePlaylist } from "~/wailsjs/go/main/App";
import { database } from '~/wailsjs/go/models';

export const usePlaylistsStore = defineStore("playlists", {
    state: () => ({
        playlists: [] as database.Playlist[],
        currentPlaylist: null as database.PlaylistWithSongs | null
    }),
    actions: {
        async loadPlaylists() {
            try {
                this.playlists = await GetPlaylists();
            } catch (err) {
                console.error(err);
            }
        },

        async getPlaylist(playlist_id: number) {
            try {
                this.currentPlaylist = await GetPlaylistWithSongs(playlist_id);
            } catch (err) {
                console.error(err);
            }
        },

        async createPlaylist(name: string, description: string, picture: string) {
            try {
                let playlist = {
                    ID: -1,
                    Name: name,
                    Description: description,
                    Picture: picture
                }

                let playlist_id = await CreatePlaylist(playlist);
                playlist.ID = playlist_id;
                
                this.playlists.push(playlist);
            } catch (err) {
                console.error(err);
            }
        },
    }
})