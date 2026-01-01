import { FFmpeg } from '@ffmpeg/ffmpeg'
import { fetchFile } from '@ffmpeg/util'
import { VIDEO_THUMBNAIL_TIME } from './constants'

let ffmpegInstance: FFmpeg | null = null
let ffmpegLoading: Promise<FFmpeg> | null = null

/**
 * Lazy-load and initialize FFmpeg instance
 */
async function getFFmpeg(): Promise<FFmpeg> {
  if (ffmpegInstance) return ffmpegInstance

  if (ffmpegLoading) return ffmpegLoading

  ffmpegLoading = (async () => {
    const ffmpeg = new FFmpeg()

    // Load FFmpeg core
    await ffmpeg.load({
      coreURL: '/ffmpeg-core.js',
      wasmURL: '/ffmpeg-core.wasm',
    })

    ffmpegInstance = ffmpeg
    ffmpegLoading = null
    return ffmpeg
  })()

  return ffmpegLoading
}

/**
 * Extract a thumbnail frame from a video file using FFmpeg
 */
export async function extractVideoThumbnail(
  file: File,
  time: number = VIDEO_THUMBNAIL_TIME,
): Promise<Blob> {
  try {
    const ffmpeg = await getFFmpeg()

    // Write input file
    const inputName = 'input.mp4'
    const outputName = 'thumbnail.jpg'
    await ffmpeg.writeFile(inputName, await fetchFile(file))

    // Extract frame at specified time
    // -ss: seek to time
    // -i: input file
    // -vframes 1: extract one frame
    // -q:v 2: quality (2 is high quality)
    await ffmpeg.exec([
      '-ss',
      time.toString(),
      '-i',
      inputName,
      '-vframes',
      '1',
      '-q:v',
      '2',
      outputName,
    ])

    // Read output file
    const data = await ffmpeg.readFile(outputName)

    // Clean up
    await ffmpeg.deleteFile(inputName)
    await ffmpeg.deleteFile(outputName)

    // Convert Uint8Array to Blob (explicit cast for TypeScript)
    return new Blob([data as BlobPart], { type: 'image/jpeg' })
  } catch (error) {
    console.error('Failed to extract video thumbnail:', error)
    throw new Error('Video thumbnail extraction failed')
  }
}

/**
 * Generate video thumbnail with fallback to generic icon
 */
export async function generateVideoThumbnail(file: File): Promise<Blob | null> {
  try {
    return await extractVideoThumbnail(file)
  } catch (error) {
    console.warn(
      'Video thumbnail generation failed, will use generic icon',
      error,
    )
    return null
  }
}
