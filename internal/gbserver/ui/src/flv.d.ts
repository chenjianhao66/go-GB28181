declare module 'flv.js' {
  export interface Player {
    on(event: string, callback: (data1?: any, data2?: any) => void): void
    off(event: string, callback?: (data1?: any, data2?: any) => void): void
    play(): void
    pause(): void
    destroy(): void
    attachMediaElement(element: HTMLMediaElement): void
    load(): void
  }

  export interface CreatePlayerOptions {
    type: string
    url: string
    hasAudio?: boolean
    hasVideo?: boolean
    cors?: boolean
    isLive?: boolean
  }

  export function createPlayer(options: CreatePlayerOptions): Player
  export function isSupported(): boolean

  export const Events: {
    ERROR: string
    LOADING: string
    MEDIA_INFO: string
    STATISTICS_INFO: string
  }

  export const ErrorTypes: {
    NETWORK_ERROR: string
    MEDIA_ERROR: string
    OTHER_ERROR: string
  }

  export const ErrorDetails: {
    OK: string
    EOF: string
    TIMEOUT: string
    UNKNOWN: string
  }
}
