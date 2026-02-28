<script setup lang="ts">

import {ref} from "vue";
import router from "./router";

const activeIndex = ref('1');
const handleSelect = (key: string) => {
  console.log("current key: ",key);
  if (key === '1') {
    router.push('/');
  } else if (key === '2') {
    router.push('/video');
  } else if (key === '3') {
    router.push('/devices');
  }
  activeIndex.value = key;
};
</script>

<template>
  <div class="common-layout">
    <el-container>
      <el-header>
        <el-menu
            :default-active="activeIndex"
            class="el-menu-demo"
            mode="horizontal"
            :ellipsis="false"
            @select="handleSelect"
        >
          <el-menu-item index="1">首页</el-menu-item>
          <el-menu-item index="2">视频监控</el-menu-item>
          <el-menu-item index="3">设备列表</el-menu-item>
        </el-menu>
      </el-header>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>


<style scoped>
</style>
