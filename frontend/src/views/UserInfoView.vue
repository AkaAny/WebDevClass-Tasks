<template>
  <h1>User Info</h1>
  <UserInfoForm :user-info="userInfo"></UserInfoForm>
  <el-button @click="doGetUserInfo">Refresh</el-button>
</template>

<script>
import UserInfoForm from "../components/UserInfoForm.vue";
import {ElMessage, ElMessageBox} from "element-plus";
import UserAPI from "../api/user_api";
import {reactive} from "vue";
export default {
  name: "UserInfoView",
  components: {UserInfoForm},
  data(){
    return {
      userInfo:reactive({}),
    }
  },
  methods:{
    doGetUserInfo(){
      UserAPI.userInfo().then((response)=>{
        const respData=response.data;
          this.userInfo=respData.data;
      }).catch((reason)=>{
        const response=reason.response;
        const respData=response.data;
        if(response.status===401){
          ElMessageBox.alert(respData.msg, 'Unauthorized', {
            // if you want to disable its autofocus
            // autofocus: false,
            confirmButtonText: 'OK',
            callback: (action) => {
              this.$router.push('/userLogin')
            },
          })
        }else{
          ElMessageBox.alert(respData.msg, 'Failed to get user info', {
            // if you want to disable its autofocus
            // autofocus: false,
            confirmButtonText: 'OK',
          })
        }
      })
    }
  },
  mounted() {
    this.doGetUserInfo();
  }
}
</script>

<style scoped>

</style>