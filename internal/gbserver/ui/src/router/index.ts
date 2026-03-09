import {createRouter, createWebHistory} from 'vue-router'
import home from "../views/home.vue"
import video from "../views/Video.vue"
import deviceList from "../views/DeviceList.vue"
import channelList from "../views/ChannelList.vue"
import mediaConfig from "../views/MediaConfig.vue"

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            redirect: '/devices'
        },
        {
            path: '/video',
            name: 'video',
            component: video
        },
        {
            path: '/channels',
            name: 'channels',
            component: channelList
        },
        {
            path: '/media',
            name: 'media',
            component: mediaConfig
        }
    ]
})

export default router