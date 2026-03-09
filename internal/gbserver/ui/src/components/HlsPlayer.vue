<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue'
import videojs from 'video.js'
import 'video.js/dist/video-js.css'

interface Props {
  url: string
}

const props = defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)
let player: any = null

const initPlayer = () => {
  if (!videoRef.value || !props.url) return
  
  if (player) {
    player.src({ src: props.url, type: 'application/x-mpegURL' })
    return
  }

  player = videojs(videoRef.value, {
    autoplay: true,
    controls: true,
    fluid: true,
    sources: [{
      src: props.url,
      type: 'application/x-mpegURL'
    }]
  })

  player.on('error', () => {
    console.error('Video playback error:', player.error())
  })
}

watch(() => props.url, (newUrl) => {
  if (newUrl && player) {
    player.src({ src: newUrl, type: 'application/x-mpegURL' })
    player.play()
  } else if (newUrl) {
    initPlayer()
  }
})

onMounted(() => {
  if (props.url) {
    initPlayer()
  }
})

onBeforeUnmount(() => {
  if (player) {
    player.dispose()
    player = null
  }
})
</script>

<template>
  <div class="hls-player">
    <video ref="videoRef" class="video-js vjs-default-skin vjs-big-play-centered"></video>
  </div>
</template>

<style scoped>
.hls-player {
  width: 100%;
  height: 100%;
}

.hls-player :deep(.video-js) {
  width: 100%;
  height: 100%;
}
</style>
