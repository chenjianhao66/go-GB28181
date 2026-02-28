<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

interface Channel {
  id: number
  deviceId: string
  name: string
  manufacturer: string
  model: string
  owner: string
  civilCode: string
  address: string
  parental: string
  parentId: string
  status: string
}

const route = useRoute()
const router = useRouter()
const deviceId = route.query.deviceId as string

const channelList = ref<Channel[]>([])
const loading = ref(false)
const error = ref('')

const fetchChannels = async () => {
  if (loading.value) return
  
  loading.value = true
  error.value = ''
  try {
    const response = await axios.get(`/channel/list/${deviceId}`)
    channelList.value = response.data.data || []
  } catch (err: any) {
    error.value = err.message || '获取通道列表失败'
    console.error('获取通道列表失败:', err)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push('/devices')
}

onMounted(() => {
  fetchChannels()
})
</script>

<template>
  <div class="channel-list">
    <div class="header">
      <div class="title">
        <el-button @click="goBack" type="primary">返回</el-button>
        <h2>通道列表 - 设备: {{ deviceId }}</h2>
      </div>
      <el-button type="primary" @click="fetchChannels" :loading="loading">
        刷新
      </el-button>
    </div>

    <el-card v-if="error" class="error-card">
      <el-alert type="error" :title="error" show-icon />
    </el-card>

    <el-table 
      :data="channelList" 
      v-loading="loading"
      stripe
      border
    >
      <el-table-column prop="deviceId" label="通道ID" width="200" />
      <el-table-column prop="name" label="通道名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="manufacturer" label="制造商" min-width="120" show-overflow-tooltip />
      <el-table-column prop="model" label="型号" min-width="120" show-overflow-tooltip />
      <el-table-column prop="owner" label="归属" min-width="120" show-overflow-tooltip />
      <el-table-column prop="civilCode" label="行政区域" min-width="120" show-overflow-tooltip />
      <el-table-column prop="address" label="安装地址" min-width="150" show-overflow-tooltip />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.status === 'ON' ? 'success' : 'info'" size="small">
            {{ row.status || '-' }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && channelList.length === 0" description="暂无通道" />
  </div>
</template>

<style scoped>
.channel-list {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.title {
  display: flex;
  align-items: center;
  gap: 15px;
}

.title h2 {
  margin: 0;
  font-size: 24px;
  color: #303133;
}

.error-card {
  margin-bottom: 20px;
}
</style>