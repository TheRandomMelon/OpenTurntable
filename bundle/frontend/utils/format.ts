export function SecondsToDuration(seconds: number) {
    const roundedSeconds = Math.round(seconds); // Round to the nearest whole number
    const minutes = Math.floor(roundedSeconds / 60);
    const remainingSeconds = roundedSeconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
}