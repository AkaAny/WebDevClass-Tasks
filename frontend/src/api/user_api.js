import AxiosInstance from "./axios";

const UserAPI={
    login:function (data){
        return AxiosInstance.post('/task2/token/',data)
    },
   userInfo:function (){
        return AxiosInstance.get('/task2/user/self/json')
   },
    register: function (userID,registerForm){
        return AxiosInstance.post(`/task2/user/${userID}`,registerForm)
    },
    active:function (activeForm){
        return AxiosInstance.put('/task2/user/active',activeForm)
    }
}

export default UserAPI;