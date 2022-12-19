<template>
  <UserLoginForm :on-login-submit="doLogin"/>
</template>

<script>
import UserLoginForm from "../components/UserLoginForm.vue";
import UserAPI from "../api/user_api";
import {ElMessage, ElMessageBox} from "element-plus";
import md5 from "js-md5";


export default {
  name: "UserLoginView",
  components: {UserLoginForm},
  methods: {
    doLogin(loginForm) {
      loginForm.password=md5(loginForm.password)
      UserAPI.login(loginForm).then((response) => {
        const respData=response.data.data;
        window.localStorage.setItem('accessToken',respData.token);
        this.$router.push('/userInfo');
      }).catch((reason) => {
        const response = reason.response;
        const respData = response.data;
        ElMessageBox.alert(respData.msg, 'Failed to login', {
          // if you want to disable its autofocus
          // autofocus: false,
          confirmButtonText: 'OK',
        })
      })
    },

  }
}
</script>

<style scoped>

</style>