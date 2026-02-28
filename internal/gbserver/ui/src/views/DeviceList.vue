<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const router = useRouter()

interface Device {
  deviceId: string
  name: string
  ip: string
  port: string
  keepalive: string
  manufacturer: string
  model: string
  firmware: string
  transport: string
  offline: number
  registerTime: string
  domain: string
}

const deviceList = ref<Device[]>([])
const loading = ref(false)
const error = ref('')
const detailVisible = ref(false)
const currentDevice = ref<Device | null>(null)

const fetchDevices = async () => {
  if (loading.value) return
  
  loading.value = true
  error.value = ''
  try {
    const response = await axios.get('/device/list')
    deviceList.value = response.data.data || []
  } catch (err: any) {
    error.value = err.message || '获取设备列表失败'
    console.error('获取设备列表失败:', err)
  } finally {
    loading.value = false
  }
}

const formatTime = (time: string) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 19)
}

const handleDetail = (row: Device) => {
  currentDevice.value = { ...row }
  detailVisible.value = true
  console.log('详情', row)
}

const handleSync = async (row: Device) => {
  try {
    await axios.post(`/device/catalog/${row.deviceId}`)
    ElMessage.success(`同步设备成功: ${row.name}`)
  } catch (err: any) {
    console.error('同步失败:', err)
    ElMessage.error(err.response?.data?.message || '同步设备失败')
  }
}

const handleDelete = (row: Device) => {
  console.log('删除', row)
  ElMessage.warning(`删除设备: ${row.name}`)
}

const handleChannel = (row: Device) => {
  router.push({ path: '/channels', query: { deviceId: row.deviceId } })
}

// 只在组件首次挂载时获取数据
fetchDevices()
</script>

<template>
  <div class="device-list">
    <div class="header">
      <h2>设备列表</h2>
      <el-button type="primary" @click="fetchDevices" :loading="loading">
        刷新
      </el-button>
    </div>

    <el-card v-if="error" class="error-card">
      <el-alert type="error" :title="error" show-icon />
    </el-card>

    <el-table 
      :data="deviceList" 
      v-loading="loading"
      stripe
      border
    >
      <el-table-column prop="deviceId" label="设备ID" width="200" />
      <el-table-column prop="name" label="设备名称" min-width="120" show-overflow-tooltip />
      <el-table-column prop="ip" label="IP" min-width="120" />
      <el-table-column prop="port" label="端口" width="80" />
      <el-table-column prop="offline" label="状态" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.offline === 1 ? 'success' : 'danger'" size="small">
            {{ row.offline === 1 ? '在线' : '离线' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最后心跳" min-width="160">
        <template #default="{ row }">
          {{ formatTime(row.keepalive) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="220" align="center">
        <template #default="{ row }">
          <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
          <el-button type="info" link @click="handleSync(row)">同步</el-button>
          <el-button type="success" link @click="handleChannel(row)">通道</el-button>
          <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-empty v-if="!loading && deviceList.length === 0" description="暂无设备" />

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="设备详情" width="600px">
      <el-form v-if="currentDevice" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="设备域">
              <el-input :value="currentDevice.domain || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="设备ID">
              <el-input :value="currentDevice.deviceId" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="设备名称">
              <el-input :value="currentDevice.name" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="IP地址">
              <el-input :value="currentDevice.ip" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="端口">
              <el-input :value="currentDevice.port" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="传输模式">
              <el-input :value="currentDevice.transport || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="制造商">
              <el-input :value="currentDevice.manufacturer || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="型号">
              <el-input :value="currentDevice.model || '-'" disabled />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="固件版本">
              <el-input :value="currentDevice.firmware || '-'" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="在线状态">
              <el-tag :type="currentDevice.offline === 1 ? 'success' : 'danger'">
                {{ currentDevice.offline === 1 ? '在线' : '离线' }}
              </el-tag>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="注册时间">
              <el-input :value="formatTime(currentDevice.registerTime)" disabled />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最后心跳">
              <el-input :value="formatTime(currentDevice.keepalive)" disabled />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.device-list {
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

.error-card {
  margin-bottom: 20px;
}
</style>