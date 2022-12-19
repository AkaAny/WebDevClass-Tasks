import axios from "axios";

const AxiosInstance=axios.create({
    baseURL:'http://localhost:8083',
    cancelToken:undefined,
})
AxiosInstance.interceptors.request.use(
    function (requestConfig){
        const authHeaderValue=window.localStorage.getItem('accessToken');
        requestConfig.headers['Authorization']=authHeaderValue;
        return requestConfig;
    }
);

export default AxiosInstance;


