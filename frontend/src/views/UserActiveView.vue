<template>
  <h1>User Active</h1>
  <el-form>
    <el-form-item label="Code">
      <el-input v-model="activeForm.code"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="doSubmit">Submit</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import UserAPI from "../api/user_api";
import {ElMessageBox} from "element-plus";
import {handleAxiosError} from "../api/handler";

export default {
  name: "UserActiveView",
  data(){
    return {
      activeForm:{
        code:'',
      }
    }
  },
  methods:{
    doSubmit(){
      UserAPI.active(this.activeForm).then((response)=>{
        ElMessageBox.confirm("Active success, go to login?",
            'Active').then((data)=>{
          console.log(data)
          this.$router.push('/userLogin')
        }).catch((reason)=>{
          //ignore
        })
      }).catch(handleAxiosError)
    }
  }
}
</script>

<style scoped>

</style>