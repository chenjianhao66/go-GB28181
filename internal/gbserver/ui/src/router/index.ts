import {createRouter, createWebHistory} from 'vue-router'
import home from "../views/home.vue"
import video from "../views/Video.vue"
import deviceList from "../views/DeviceList.vue"
import channelList from "../views/ChannelList.vue"

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: home
        },
        {
            path: '/video',
            name: 'video',
            component: video
        },
        {
            path: '/devices',
            name: 'devices',
            component: deviceList
        },
        {
            path: '/channels',
            name: 'channels',
            component: channelList
        }
    ]
})

export default router