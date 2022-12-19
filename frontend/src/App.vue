<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <el-container direction="horizontal">
    <el-aside>
      <el-menu>
        <el-sub-menu index="1">
          <template #title>
            <el-icon><message /></el-icon>Navigator One
          </template>
          <el-menu-item-group>
            <template #title>Book</template>
            <el-menu-item index="1-1" @click="this.$router.push('/bookInfo');">Book Info</el-menu-item>
            <el-menu-item index="1-2" @click="this.$router.push('/bookRent');">Book Rent</el-menu-item>
          </el-menu-item-group>
          <el-menu-item-group title="User">
            <el-menu-item index="1-3" @click="this.$router.push('/userLogin');">User Login</el-menu-item>
            <el-menu-item index="1-4" @click="this.$router.push('/userInfo');">User Info</el-menu-item>
            <el-menu-item index="1-5" @click="this.$router.push('/userRegister');">User Register</el-menu-item>
            <el-menu-item index="1-6" @click="this.$router.push('/userActive');">User Active</el-menu-item>
          </el-menu-item-group>
        </el-sub-menu>
      </el-menu>
    </el-aside>
    <el-main>
      <RouterView />
    </el-main>
  </el-container>
</template>

<style scoped>

</style>
