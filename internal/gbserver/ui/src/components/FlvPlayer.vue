<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import flvjs from 'flv.js'

interface Props {
  url?: string
}

const props = defineProps<Props>()

const videoRef = ref<HTMLVideoElement | null>(null)
const errorMsg = ref('')
let flvPlayer: flvjs.Player | null = null

const initPlayer = () => {
  if (!props.url || !videoRef.value) return
  
  errorMsg.value = ''
  
  if (flvPlayer) {
    flvPlayer.destroy()
    flvPlayer = null
  }
  
  try {
    flvPlayer = flvjs.createPlayer({
      type: 'flv',
      url: props.url,
      hasAudio: true
    })
    
    flvPlayer.attachMediaElement(videoRef.value)
    flvPlayer.load()
    flvPlayer.play()
    
    flvPlayer.on(flvjs.Events.ERROR, (errType: string, errDetail: string) => {
      console.error('FLV播放错误:', errType, errDetail)
      if (errType === flvjs.ErrorTypes.NETWORK_ERROR) {
        errorMsg.value = '网络错误，请检查流媒体服务'
      } else if (errType === flvjs.ErrorTypes.MEDIA_ERROR) {
        errorMsg.value = '媒体格式不支持'
      } else {
        errorMsg.value = '播放错误: ' + errType
      }
    })
  } catch (e: any) {
    console.error('创建播放器失败:', e)
    errorMsg.value = '创建播放器失败: ' + e.message
  }
}

const destroyPlayer = () => {
  if (flvPlayer) {
    flvPlayer.destroy()
    flvPlayer = null
  }
}

watch(() => props.url, (newUrl) => {
  if (newUrl) {
    initPlayer()
  } else {
    destroyPlayer()
  }
})

onMounted(() => {
  if (props.url) {
    initPlayer()
  }
})

onBeforeUnmount(() => {
  destroyPlayer()
})

defineExpose({
  destroyPlayer
})
</script>

<template>
  <div class="video-player">
    <video 
      ref="videoRef" 
      controls 
      muted
      autoplay
      class="player-video"
    >
    </video>
    <div v-if="errorMsg" class="error-overlay">
      {{ errorMsg }}
    </div>
  </div>
</template>

<style scoped>
.video-player {
  width: 100%;
  height: 100%;
  background: #000;
  position: relative;
}

.player-video {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.error-overlay {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(255, 0, 0, 0.8);
  color: white;
  padding: 20px;
  border-radius: 8px;
  text-align: center;
}
</style>