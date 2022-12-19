<template>
  <h1>User Register</h1>
  <el-form>
    <el-form-item label="UserID">
      <el-input v-model="registerForm.userID"/>
    </el-form-item>
    <el-form-item label="Mail">
      <el-input v-model="registerForm.mail"/>
    </el-form-item>
    <el-form-item label="Password">
      <el-input type="password" v-model="registerForm.password"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="doRegister">Submit</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import UserAPI from "../api/user_api";
import {handleAxiosError} from "../api/handler";
import {ElMessageBox} from "element-plus";
import md5 from "js-md5";

export default {
  name: "UserRegisterView",
  data(){
    return {
      registerForm:{
        userID:'',
        mail:'',
        password:'',
      }
    }
  },
  methods:{
    doRegister(){
      this.registerForm.password=md5(this.registerForm.password);
      UserAPI.register(this.registerForm.userID,this.registerForm).then((response)=>{
        ElMessageBox.alert("Register success, please check your email to active account",
            'Register').then((data)=>{
          console.log(data)
          this.$router.push('/userActive')
        })
      }).catch(handleAxiosError)
    }
  },
}
</script>

<style scoped>

</style>