<script setup lang="ts">
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface MediaDetail {
  ID: string
  ip: string
  hookIp: string
  sdpIp: string
  streamIp: string
  status: boolean
  httpPort: number
  httpSSlPort: number
  rtmpPort: number
  rtmpSSlPort: number
  rtpProxyPort: number
  rtspPort: number
  rtspSSLPort: number
  rtpEnable: boolean
  rtpPortRange: string
  secret: string
  default: boolean
  hookAliveInterval: number
  createTime: string
  updateTime: string
  lastKeepaliveTime: string
}

const mediaList = ref<MediaDetail[]>([])
const loading = ref(false)
const dialogVisible = ref(false)
const currentMedia = ref<MediaDetail | null>(null)

const fetchMediaList = async () => {
  loading.value = true
  try {
    const response = await axios.get('/media/list')
    mediaList.value = response.data.data || []
  } catch (err: any) {
    ElMessage.error('获取ZLM列表失败')
    console.error('获取ZLM列表失败:', err)
  } finally {
    loading.value = false
  }
}

const handleCardClick = (media: MediaDetail) => {
  currentMedia.value = { ...media }
  dialogVisible.value = true
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

onMounted(() => {
  fetchMediaList()
})
</script>

<template>
  <div class="media-config">
    <div class="header">
      <h2>ZLM配置</h2>
      <el-button type="primary" @click="fetchMediaList" :loading="loading">
        刷新
      </el-button>
    </div>

    <div class="card-container" v-loading="loading">
      <el-card 
        v-for="media in mediaList" 
        :key="media.ID" 
        class="media-card"
        :class="{ 'online': media.status, 'offline': !media.status }"
        @click="handleCardClick(media)"
      >
        <div class="card-image">
          <div class="placeholder-icon"><img src="../assets/zlm.png"/></div>
        </div>
        <div class="card-content">
          <h3 class="media-id">{{ media.ID }}</h3>
          <p class="media-ip">IP: {{ media.ip }}</p>
          <div class="media-status">
            <el-tag :type="media.status ? 'success' : 'danger'" size="small">
              {{ media.status ? '在线' : '离线' }}
            </el-tag>
          </div>
        </div>
      </el-card>

      <el-empty v-if="!loading && mediaList.length === 0" description="暂无ZLM配置" />
    </div>

    <el-dialog v-model="dialogVisible" title="ZLM详情" width="800px">
      <el-form v-if="currentMedia" label-width="120px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="ZLM ID">
              <el-input :value="currentMedia.ID" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="状态">
              <el-tag :type="currentMedia.status ? 'success' : 'danger'">
                {{ currentMedia.status ? '在线' : '离线' }}
              </el-tag>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="IP">
              <el-input :value="currentMedia.ip" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Hook IP">
              <el-input :value="currentMedia.hookIp || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="SDP IP">
              <el-input :value="currentMedia.sdpIp || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="Stream IP">
              <el-input :value="currentMedia.streamIp || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="HTTP端口">
              <el-input :value="currentMedia.httpPort" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="HTTPS端口">
              <el-input :value="currentMedia.httpSSlPort || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="RTMP端口">
              <el-input :value="currentMedia.rtmpPort" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="RTMPS端口">
              <el-input :value="currentMedia.rtmpSSlPort || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="RTSP端口">
              <el-input :value="currentMedia.rtspPort" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="RTSPS端口">
              <el-input :value="currentMedia.rtspSSLPort || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="RTP代理端口">
              <el-input :value="currentMedia.rtpProxyPort || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="RTP端口范围">
              <el-input :value="currentMedia.rtpPortRange || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="RTP启用">
              <el-tag :type="currentMedia.rtpEnable ? 'success' : 'info'">
                {{ currentMedia.rtpEnable ? '是' : '否' }}
              </el-tag>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="默认">
              <el-tag :type="currentMedia.default ? 'success' : 'info'">
                {{ currentMedia.default ? '是' : '否' }}
              </el-tag>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="API密钥">
              <el-input :value="currentMedia.secret || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="心跳间隔">
              <el-input :value="currentMedia.hookAliveInterval ? currentMedia.hookAliveInterval + '秒' : '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="创建时间">
              <el-input :value="formatTime(currentMedia.createTime)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="更新时间">
              <el-input :value="formatTime(currentMedia.updateTime)" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="24">
            <el-form-item label="最后心跳">
              <el-input :value="formatTime(currentMedia.lastKeepaliveTime)" disabled />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.media-config {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.card-container {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  min-height: 200px;
}

.media-card {
  width: 280px;
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.media-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.media-card.online {
  border-left: 4px solid #67c23a;
}

.media-card.offline {
  border-left: 4px solid #f56c6c;
}

.card-image {
  height: 120px;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 15px;
  border-radius: 4px;
}

.placeholder-icon {
  font-size: 48px;
}

.card-image img {
  max-height: 100%;
  max-width: 100%;
  object-fit: contain;
}

.card-content {
  text-align: center;
}

.media-id {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: #303133;
  word-break: break-all;
}

.media-ip {
  margin: 5px 0;
  color: #606266;
  font-size: 14px;
}

.media-status {
  margin: 10px 0;
}
</style>