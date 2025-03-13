// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {database} from '../models';

export function ChooseAndCreateSong():Promise<number>;

export function CreateSong(arg1:database.Song):Promise<number>;

export function GetDuration():Promise<number>;

export function GetFilePath():Promise<string>;

export function GetMetadata():Promise<Record<string, string>>;

export function GetPosition():Promise<number>;

export function GetSongs():Promise<Array<database.Song>>;

export function IsPlaying():Promise<boolean>;

export function PauseMusic():Promise<void>;

export function PlayFile(arg1:string):Promise<void>;

export function Seek(arg1:number):Promise<void>;

export function SelectAndPlayFile():Promise<void>;

export function SetVolume(arg1:number):Promise<void>;
