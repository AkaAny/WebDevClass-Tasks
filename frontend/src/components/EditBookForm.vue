<template>
  <el-form :model="bookForm">
    <el-form-item label="ID">
      <el-input v-model="bookID_" disabled/>
    </el-form-item>
    <el-form-item label="Name">
      <el-input v-model="bookForm.name"/>
    </el-form-item>
    <el-form-item label="Author">
      <el-input v-model="bookForm.author"/>
    </el-form-item>
    <el-form-item>
      <el-button type="primary" @click="handleSubmit">Submit</el-button>
      <el-button @click="handleCancel">Cancel</el-button>
    </el-form-item>
  </el-form>
</template>

<script>
import {reactive} from "vue";

export default {
  name: "EditBookForm",
  props:{
    bookID:0,
    initialBookForm:{name:'',author:''},
    onSubmit:function (bookID,bookForm){},
    onCancel:function (bookID,bookForm){},
  },
  data(){
    return {
      bookID_:reactive(0),
      bookForm:reactive({name:'',author:''})
    }
  },
  methods:{
    handleSubmit(){
      this.onSubmit(this.bookID,this.bookForm);
    },
    handleCancel(){
      this.onCancel(this.bookID,this.bookForm);
    }
  },
  mounted() {
    this.bookID_=this.bookID;
    this.bookForm=this.initialBookForm;
  }
}
</script>

<style scoped>

</style>