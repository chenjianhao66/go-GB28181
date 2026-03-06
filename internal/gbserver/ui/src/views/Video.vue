<script setup lang="ts">
import { ref, onMounted, shallowRef } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'
import HlsPlayer from '../components/HlsPlayer.vue'

interface Device {
  deviceId: string
  name: string
}

interface Channel {
  id: number
  deviceId: string
  name: string
  status: string
}

interface TreeNode {
  id: string
  label: string
  children?: TreeNode[]
  isDevice: boolean
  deviceId?: string
  channelData?: Channel
}

interface StreamInfo {
  mediaServerId: string
  app: string
  ip: string
  deviceID: string
  channelId: string
  stream: string
  rtmp: string
  rtsp: string
  flv: string
  httpsFlv: string
  hls: string
}

const loading = ref(false)
const treeData = ref<TreeNode[]>([])
const selectedChannel = ref<Channel | null>(null)
const leftWidth = ref(300)
const streamUrl = ref('')
const playerRef = shallowRef()
const playing = ref(false)

// 获取设备列表
const fetchDevices = async () => {
  loading.value = true
  try {
    const response = await axios.get('/device/list')
    const devices: Device[] = response.data.data || []
    
    // 构建树形数据
    treeData.value = devices.map(device => ({
      id: device.deviceId,
      label: device.name || device.deviceId,
      isDevice: true,
      deviceId: device.deviceId,
      children: []  // 初始为空，点击时加载
    }))
  } catch (err: any) {
    ElMessage.error('获取设备列表失败')
    console.error('获取设备列表失败:', err)
  } finally {
    loading.value = false
  }
}

// 加载设备的通道
const loadChildren = async (node: TreeNode) => {
  if (!node.isDevice || !node.deviceId) return
  
  try {
    const response = await axios.get(`/channel/list/${node.deviceId}`)
    const channels: Channel[] = response.data.data || []
    
    node.children = channels.map(channel => ({
      id: `${node.deviceId}-${channel.id}`,
      label: channel.name || channel.deviceId,
      isDevice: false,
      deviceId: node.deviceId,
      channelData: channel
    }))
  } catch (err: any) {
    ElMessage.error(`获取通道列表失败: ${node.label}`)
    console.error('获取通道列表失败:', err)
  }
}

// 播放视频
const playVideo = async (channel: Channel, deviceId: string) => {
  try {
    const response = await axios.post(`/play/start/${deviceId}/${channel.deviceId}`)
    const streamInfo: StreamInfo = response.data.data
    
    if (!streamInfo) {
      ElMessage.error('获取播放地址失败')
      return
    }
    
    // 优先使用 HLS 地址播放
    if (streamInfo.hls) {
      streamUrl.value = streamInfo.hls
    } else if (streamInfo.flv) {
      streamUrl.value = streamInfo.flv
    } else if (streamInfo.fmp4) {
      streamUrl.value = streamInfo.fmp4
    } else {
      ElMessage.error('无可用的播放地址')
      return
    }
    
    playing.value = true
  } catch (err: any) {
    ElMessage.error(err.response?.data?.message || '播放失败')
    console.error('播放失败:', err)
  }
}

// 点击节点处理
const handleNodeClick = async (data: TreeNode) => {
  if (data.isDevice) {
    // 设备节点，展开/收起子节点
    return
  } else {
    // 通道节点，选中并播放
    selectedChannel.value = data.channelData || null
    if (data.channelData && data.deviceId) {
      await playVideo(data.channelData, data.deviceId)
    }
  }
}

// 加载子节点
const handleNodeExpand = (data: TreeNode) => {
  if (data.isDevice && data.children && data.children.length === 0) {
    loadChildren(data)
  }
}

onMounted(() => {
  fetchDevices()
})
</script>

<template>
  <div class="video-container">
    <div class="left-panel" :style="{ width: leftWidth + 'px' }">
      <div class="panel-header">
        <h3>设备列表</h3>
      </div>
      <el-tree
        :data="treeData"
        :props="{
          children: 'children',
          label: 'label'
        }"
        :load="(node: TreeNode, resolve: Function) => {
          if (node.level === 0) {
            resolve(treeData)
          } else if (node.data.isDevice) {
            loadChildren(node.data).then(() => {
              resolve(node.data.children || [])
            })
          } else {
            resolve([])
          }
        }"
        lazy
        node-key="id"
        @node-click="handleNodeClick"
        @node-expand="handleNodeExpand"
        v-loading="loading"
      >
        <template #default="{ node, data }">
          <span class="tree-node">
            <span class="node-icon">
              {{ data.isDevice ? '📱' : '📹' }}
            </span>
            <span class="node-label">{{ node.label }}</span>
          </span>
        </template>
      </el-tree>
    </div>
    
    <div class="resize-handle" @mousedown="startResize"></div>
    
    <div class="right-panel">
      <div v-if="playing && streamUrl" class="player-container">
        <HlsPlayer ref="playerRef" :url="streamUrl" />
      </div>
      <div v-else-if="selectedChannel" class="channel-info">
        <h3>已选择通道: {{ selectedChannel.name }}</h3>
        <p>通道ID: {{ selectedChannel.deviceId }}</p>
        <p>状态: {{ selectedChannel.status }}</p>
        <p>正在获取播放地址...</p>
      </div>
      <div v-else class="empty-hint">
        <el-empty description="请选择通道" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
// 拖拽调整宽度
export default {
  methods: {
    startResize(e: MouseEvent) {
      const startX = e.clientX
      const startWidth = (this as any).leftWidth
      
      const doDrag = (e: MouseEvent) => {
        (this as any).leftWidth = Math.max(200, Math.min(600, startWidth + e.clientX - startX))
      }
      
      const stopDrag = () => {
        document.removeEventListener('mousemove', doDrag)
        document.removeEventListener('mouseup', stopDrag)
      }
      
      document.addEventListener('mousemove', doDrag)
      document.addEventListener('mouseup', stopDrag)
    }
  }
}
</script>

<style scoped>
.video-container {
  display: flex;
  height: calc(100vh - 60px);
  overflow: hidden;
}

.left-panel {
  min-width: 200px;
  max-width: 600px;
  border-right: 1px solid #dcdfe6;
  overflow: auto;
  background: #fff;
}

.panel-header {
  padding: 15px;
  border-bottom: 1px solid #dcdfe6;
}

.panel-header h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.tree-node {
  display: flex;
  align-items: center;
  gap: 5px;
}

.node-icon {
  font-size: 14px;
}

.node-label {
  font-size: 14px;
}

.resize-handle {
  width: 5px;
  cursor: col-resize;
  background: transparent;
  transition: background 0.2s;
}

.resize-handle:hover {
  background: #409eff;
}

.right-panel {
  flex: 1;
  padding: 20px;
  overflow: auto;
  background: #000;
}

.player-container {
  width: 100%;
  height: 100%;
}

.channel-info h3 {
  margin: 0 0 15px 0;
  color: #fff;
}

.channel-info p {
  margin: 8px 0;
  color: #ccc;
}

.empty-hint {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}
</style>