import { defineStore } from "pinia";
import { GetDuration, GetPosition, PauseMusic, StopPlayback, SetVolume, GetFilePath, GetMetadata, IsPlaying, Seek, PlayFile } from "~/wailsjs/go/main/App";
import { database } from '~/wailsjs/go/models';

export enum PlaybackSourceType {
    Library,
    Playlist
}

export enum RepeatType {
    Off,
    Repeat,
    RepeatOne
}

export const usePlaybackStore = defineStore("playback", {
    state: () => ({
        filePath: null as string | null,
        metadata: null as Record<string, string> | null,
        position: null as number | null,
        duration: null as number | null,
        playing: false as boolean,
        currentSong: null as database.SongWithDetails | null,
        volume: 0 as number,
        prevVolume: 0 as number,
        source: {
            type: null as PlaybackSourceType | null,
            id: null as number | null
        },
        queue: null as database.SongWithDetails[] | null,
        unshuffledQueue: null as database.SongWithDetails[] | null,
        repeat: RepeatType.Off as RepeatType,
        shuffle: false as boolean
    }),
    actions: {
        async beginPlayback(song: database.SongWithDetails, source: PlaybackSourceType) {
            try {
                const songs = useSongsStore();

                switch (source) {
                    case PlaybackSourceType.Library:
                        // If songs is null, get all songs in library
                        if (!songs.songs) await songs.getAllSongs();

                        // Find position in library
                        let pos = songs.songs?.findIndex((s) => s.ID === song.ID)

                        // Get new queue
                        this.queue = songs.getQueue<database.SongWithDetails>(songs.songs ? songs.songs : [], pos ? pos : song.ID, false);
                        this.unshuffledQueue = this.queue;

                        if (this.shuffle) {
                            this.queue = songs.getQueue<database.SongWithDetails>(songs.songs ? songs.songs : [], pos ? pos : song.ID, true);
                        }

                        console.log(this.queue);

                        // Begin playback of file
                        await this.playFile(song);
                        break;
                    case PlaybackSourceType.Playlist:
                        // TODO
                        break;
                    default:
                        throw new Error("Unsupported playback source type");
                }
            } catch (err) {
                console.error(err);
            }
        },

        async playFile(song: database.SongWithDetails) {
            await PlayFile(song.Path);
            this.currentSong = song;
            this.position = await GetPosition();
            this.duration = await GetDuration();
            this.filePath = await this.getFilePath();
            this.metadata = await GetMetadata();
            
            this.playing = true;
        },

        async queueStep(forward: boolean, userClick: boolean = false) {
            if (forward && this.queue) {
                let currentPos = this.queue.findIndex((s) => s.ID === this.currentSong?.ID);
                let newPos = currentPos ? currentPos + 1 : 1;
                let newSong;

                if (this.repeat === RepeatType.RepeatOne && !userClick) {
                    newSong = this.queue[currentPos];
                } else {
                    newSong = this.queue[newPos];
                    if (this.repeat === RepeatType.RepeatOne) {
                        this.repeat = RepeatType.Repeat;
                    }
                }

                if (!newSong && this.repeat !== RepeatType.Repeat) {
                    // Clear everything out, queue is done
                    await this.stopPlayback();
                } else if (!newSong && this.repeat === RepeatType.Repeat) {
                    // Loop back to first position in queue
                    newSong = this.queue[0];
                }

                await this.playFile(newSong);
            } else if (!forward && this.queue) {
                let currentPos = this.queue.findIndex((s) => s.ID === this.currentSong?.ID);
                let newPos = currentPos ? currentPos - 1 : 0;
                let newSong = this.queue[newPos];

                if (!newSong) newSong = this.queue[currentPos];

                if (this.repeat === RepeatType.RepeatOne) {
                    this.repeat = RepeatType.Repeat;
                }

                await this.playFile(newSong);
            } else {
                console.log("No queue exists");
            }
        },

        async cycleRepeat() {
            switch (this.repeat) {
                case RepeatType.Off:
                    this.repeat = RepeatType.Repeat;
                    break;
                case RepeatType.Repeat:
                    this.repeat = RepeatType.RepeatOne;
                    break;
                case RepeatType.RepeatOne:
                    this.repeat = RepeatType.Off;
                    break;
                default:
                    this.repeat = RepeatType.Off;
                    break;
            }
        },

        async reloadData() {
            try {
                this.filePath = await GetFilePath();
                this.position = await GetPosition();
                this.duration = await GetDuration();
                this.metadata = await GetMetadata();
                this.playing = await IsPlaying();
            } catch (err) {
                console.error("Error on reload: " + err);
            }
        },

        async togglePlayback() {
            await PauseMusic();
            this.playing = await this.isPlaying();
        },

        async getDuration() {
            this.duration = await GetDuration();
            return this.duration;
        },

        async getPosition() {
            this.position = await GetPosition();
            return this.position;
        },

        async getFilePath() {
            this.filePath = await GetFilePath();
            return this.filePath;
        },

        async isPlaying() {
            this.playing = await IsPlaying();
            return this.playing;
        },

        async changePosition(event: Event) {
            let newPosE = event.target as HTMLInputElement;
            let newPos = parseFloat(newPosE.value);
            
            await Seek(newPos);
            this.position = await this.getPosition();
        },

        async stopPlayback() {
            await StopPlayback();
            this.currentSong = null;
            this.playing = false;
            this.queue = null;
            this.source = {
                type: null,
                id: null
            };
            this.filePath = null;
            this.duration = null;
            this.position = null;
            this.metadata = null;
            this.repeat = RepeatType.Off;
        },
        
        async toggleShuffle() {
            const songs = useSongsStore();

            if (this.shuffle) {
                this.shuffle = false;
                if (this.queue && this.currentSong && this.unshuffledQueue) {
                    // Get new queue
                    this.queue = songs.getQueue<database.SongWithDetails>(this.unshuffledQueue, this.currentSong.ID, false);
                    console.log(this.queue);
                }
            } else {
                this.shuffle = true;
                if (this.queue && this.currentSong) {
                    // Get new queue
                    this.queue = songs.getQueue<database.SongWithDetails>(this.queue, this.currentSong.ID, true);
                    console.log(this.queue);
                }
            }
        },

        updateVolume(event: Event) {
            const target = event.target as HTMLInputElement;
            SetVolume(parseFloat(target.value));
        },

        setVolume(volume: number) {
            if (typeof volume === "string") {
                volume = parseFloat(volume);
            }
            
            this.volume = volume;
            SetVolume(volume);
        },

        toggleMute() {
            if (this.volume > -5) {
                this.prevVolume = this.volume;
                this.setVolume(-5);
            } else {
                this.setVolume(this.prevVolume);
            }
        }
    }
})