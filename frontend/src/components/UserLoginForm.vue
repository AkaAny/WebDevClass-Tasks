<template>
  <el-form>
    <el-form-item label="Mail">
      <el-input v-model="loginForm.mail"/>
    </el-form-item>
    <el-form-item label="Password">
      <el-input type="password" v-model="loginForm.password"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="loginClick">Submit</el-button>
    </el-form-item>
  </el-form>
</template>

<script>

export default {
  name: "UserLoginForm",
  props:{
    onLoginSubmit:function (loginForm){},
  },
  data(){
    return {
      loginForm:{
        mail:'',
        password:'',
      }
    }
  },
  methods:{
    loginClick(){
      this.onLoginSubmit(this.loginForm);
    }
  }
}
</script>

<style scoped>

</style>