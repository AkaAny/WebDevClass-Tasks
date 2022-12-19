<template>
  <h1>Book Rent View</h1>
  <el-table :data="records">
    <el-table-column label="ID" prop="id" />
    <el-table-column label="Book ID" prop="book.id" />
    <el-table-column label="Book Name" prop="book.name" />
    <el-table-column label="Rent By" prop="rentBy"/>
    <el-table-column label="ExpiredAt" prop="expiredAt" />
    <el-table-column label="Returned At" prop="returnedAt" />
    <el-table-column label="Operation" align="right">
      <template #default="scope">
        <el-button v-if="checkIfNotReturn(scope.row)" size="small"
                   @click="handleReturn(scope.$index, scope.row)">Return</el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-container>
    <el-button @click="doGetRecords">Refresh</el-button>
  </el-container>
  <el-dialog v-model="addDialogVisible">

  </el-dialog>
</template>

<script>
import {reactive} from "vue";
import BookAPI from "../api/book_api";
import {handleAxiosError} from "../api/handler";
import {ElMessage, ElMessageBox} from "element-plus";

export default {
  name: "BookRentView",
  data(){
    return {
      exampleRecord:{
        id:0,
        book:{
          id:0,
          name:'',
        },
        rentBy:'',
        expiredAt:new Date(),
        returnedAt:new Date(),
      },
      records:reactive([]),
      addDialogVisible:reactive(false),
    }
  },
  methods:{
    doGetRecords(){
      return BookAPI.listRent().then((response)=>{
        this.records=response.data.data;
      }).catch(handleAxiosError)
    },
    checkIfNotReturn(row){
      const returnedAtAsDate= new Date(row.returnedAt); //0001-01-01T00:00:00Z
      return returnedAtAsDate.valueOf()<0;
    },
    handleReturn(i,row){
      const rentID=row.id;
      ElMessageBox.confirm(
          `Are you sure to return book with rent id:${rentID}`, 'Warning',
          {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'warning',
          }).then(() => {
        this.doReturn(rentID);
      }).catch((reason)=>{
        //ignore
      })
    },
    doReturn(rentID){
      BookAPI.returnRent(rentID).then((response)=>{
        ElMessage.success({
          message: `success`,
        })
        this.doGetRecords();
      }).catch(handleAxiosError)
    },
  },
  mounted() {
    this.doGetRecords();
  }
}
</script>

<style scoped>

</style>