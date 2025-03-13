export namespace database {
	
	export class Song {
	    ID: number;
	    Path: string;
	    Title: string;
	    Artist_ID: number;
	    Album_ID: number;
	    Composer: string;
	    Comment: string;
	    Genre: string;
	    Year: string;
	
	    static createFrom(source: any = {}) {
	        return new Song(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Path = source["Path"];
	        this.Title = source["Title"];
	        this.Artist_ID = source["Artist_ID"];
	        this.Album_ID = source["Album_ID"];
	        this.Composer = source["Composer"];
	        this.Comment = source["Comment"];
	        this.Genre = source["Genre"];
	        this.Year = source["Year"];
	    }
	}

}

