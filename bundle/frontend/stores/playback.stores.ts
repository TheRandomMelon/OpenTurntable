import { defineStore } from "pinia";
import { GetDuration, GetPosition, PauseMusic, SelectAndPlayFile, SetVolume, GetFilePath, GetMetadata, IsPlaying, Seek } from "~/wailsjs/go/main/App";

export const usePlaybackStore = defineStore("playback", {
    state: () => ({
        filePath: null as string | null,
        metadata: null as Record<string, string> | null,
        position: null as number | null,
        duration: null as number | null,
        playing: false as boolean,
        volume: 0 as number,
        prevVolume: 0 as number
    }),
    actions: {
        async beginPlayback() {
            try {
                await SelectAndPlayFile();
                this.position = await GetPosition();
                this.duration = await GetDuration();
                this.filePath = await this.getFilePath();
                this.metadata = await GetMetadata();
                
                this.playing = true;
            } catch (err) {
                console.error(err);
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