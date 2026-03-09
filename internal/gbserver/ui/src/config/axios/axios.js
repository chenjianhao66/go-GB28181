import {service} from "@/config/axios/index.js";

const request = (option, headerParams = {}) => {
    const { url, method, params, data, headersType, responseType } = option;
    return service({
        url,
        method,
        params,
        data,
        responseType,
        headers: {
            'Content-Type': headersType ||'application/json',
            ...headerParams
        }
    });
};

/***
 * headers:头部入参
 */
function getFn(option, headers) {
    return request({ method: 'get', ...option }, headers);
}

function postFn(option, headers) {
    return request({ method: 'post', ...option }, headers);
}

function deleteFn(option, headers) {
    return request({ method: 'delete', ...option }, headers);
}

function putFn(option, headers) {
    return request({ method: 'put', ...option }, headers);
}

export const useAxios = () => {
    return {
        get: getFn,
        post: postFn,
        delete: deleteFn,
        put: putFn
    };
};
