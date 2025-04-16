export namespace database {
	
	export class Playlist {
	    ID: number;
	    Name: string;
	    Description: string;
	    Picture: string;
	
	    static createFrom(source: any = {}) {
	        return new Playlist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Description = source["Description"];
	        this.Picture = source["Picture"];
	    }
	}
	export class SongWithDetails {
	    ID: number;
	    Path: string;
	    Title: string;
	    Artist_ID: sql.NullInt64;
	    Album_ID: sql.NullInt64;
	    Composer: sql.NullString;
	    Comment: sql.NullString;
	    Genre: sql.NullString;
	    Year: sql.NullString;
	    ArtistName: sql.NullString;
	    ArtistPFP: sql.NullString;
	    AlbumName: sql.NullString;
	    AlbumArt: sql.NullString;
	
	    static createFrom(source: any = {}) {
	        return new SongWithDetails(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Path = source["Path"];
	        this.Title = source["Title"];
	        this.Artist_ID = this.convertValues(source["Artist_ID"], sql.NullInt64);
	        this.Album_ID = this.convertValues(source["Album_ID"], sql.NullInt64);
	        this.Composer = this.convertValues(source["Composer"], sql.NullString);
	        this.Comment = this.convertValues(source["Comment"], sql.NullString);
	        this.Genre = this.convertValues(source["Genre"], sql.NullString);
	        this.Year = this.convertValues(source["Year"], sql.NullString);
	        this.ArtistName = this.convertValues(source["ArtistName"], sql.NullString);
	        this.ArtistPFP = this.convertValues(source["ArtistPFP"], sql.NullString);
	        this.AlbumName = this.convertValues(source["AlbumName"], sql.NullString);
	        this.AlbumArt = this.convertValues(source["AlbumArt"], sql.NullString);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlaylistEntryWithSong {
	    ID: number;
	    Playlist_ID: number;
	    ListOrder: number;
	    Song: SongWithDetails;
	
	    static createFrom(source: any = {}) {
	        return new PlaylistEntryWithSong(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Playlist_ID = source["Playlist_ID"];
	        this.ListOrder = source["ListOrder"];
	        this.Song = this.convertValues(source["Song"], SongWithDetails);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlaylistWithSongs {
	    Playlist: Playlist;
	    Entries: PlaylistEntryWithSong[];
	
	    static createFrom(source: any = {}) {
	        return new PlaylistWithSongs(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Playlist = this.convertValues(source["Playlist"], Playlist);
	        this.Entries = this.convertValues(source["Entries"], PlaylistEntryWithSong);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class Song {
	    ID: number;
	    Path: string;
	    Title: string;
	    Artist_ID: sql.NullInt64;
	    Album_ID: sql.NullInt64;
	    Composer: sql.NullString;
	    Comment: sql.NullString;
	    Genre: sql.NullString;
	    Year: sql.NullString;
	
	    static createFrom(source: any = {}) {
	        return new Song(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Path = source["Path"];
	        this.Title = source["Title"];
	        this.Artist_ID = this.convertValues(source["Artist_ID"], sql.NullInt64);
	        this.Album_ID = this.convertValues(source["Album_ID"], sql.NullInt64);
	        this.Composer = this.convertValues(source["Composer"], sql.NullString);
	        this.Comment = this.convertValues(source["Comment"], sql.NullString);
	        this.Genre = this.convertValues(source["Genre"], sql.NullString);
	        this.Year = this.convertValues(source["Year"], sql.NullString);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace sql {
	
	export class NullInt64 {
	    Int64: number;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullInt64(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Int64 = source["Int64"];
	        this.Valid = source["Valid"];
	    }
	}
	export class NullString {
	    String: string;
	    Valid: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NullString(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.String = source["String"];
	        this.Valid = source["Valid"];
	    }
	}

}

