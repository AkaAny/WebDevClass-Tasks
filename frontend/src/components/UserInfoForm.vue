<template>
  <el-descriptions title="User Info"
                   direction="vertical"
                   :column="2"
                   :size="size"
                   border>
    <el-descriptions-item label="UserID">{{ userInfo.userID }}</el-descriptions-item>
    <el-descriptions-item label="Mail">{{ userInfo.mail }}</el-descriptions-item>
    <el-descriptions-item label="Roles">
      <el-tag v-for="roleID in this.userInfo.roles" :key="index" size="small">{{ roleID }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item label="Status">{{ userInfo.status }}</el-descriptions-item>
  </el-descriptions>
</template>

<script>
export default {
  name: "UserInfoForm",
  props:{
    userInfo:{userID:"",mail:"",roles:[],status:0},
  }
}
</script>

<style scoped>

</style>