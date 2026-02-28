import axios from "axios";
import {ElMessage} from "element-plus";

const service = axios.create({
    timeout: 60000, // 请求超时时间
    // headers: {
    //     'Content-Type': 'application/json; charset=utf-8'
    // }
});

// request拦截器
service.interceptors.request.use(
    config => {

        return config;
    },
    error => {
        // Do something with request error
        console.log(error); // for debug
        Promise.reject(error);
    }
);


// response 拦截器
service.interceptors.response.use(
    response => {
        let { code, data, msg } = response.data;
        if (code === 200) {
            // ElMessage.success(msg || '请求成功!');
            return data;
        }  else {
            if (['arraybuffer', 'blob', 'text'].includes(response.config.responseType)) {
                return response.data;
            } else {
                ElMessage.error(msg || '请求失败!');
                return Promise.reject(new Error(msg));
            }
        }
    },
    error => {
        if (error.message === '路由跳转取消请求') {
            console.log('路由跳转取消请求' + error);
        } else {
            ElMessage.error(error.message);
            return Promise.reject(error);
        }
    }
);


export { service };
