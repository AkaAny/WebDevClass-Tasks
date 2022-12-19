<template>
  <h1>Book View</h1>
  <el-table :data="books">
    <el-table-column label="ID" prop="id" />
    <el-table-column label="Name" prop="name" />
    <el-table-column label="Author" prop="author" />
    <el-table-column label="Operation" align="right">
      <template #default="scope">
        <el-button size="small"
                   @click="handleRent(scope.$index, scope.row)">Rent</el-button>
        <el-button size="small"
                   @click="handleEdit(scope.$index, scope.row)">Edit</el-button>
        <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.$index, scope.row)">Delete</el-button>
      </template>
    </el-table-column>
  </el-table>
  <el-container>
    <el-button @click="doGetBook">Refresh</el-button>
    <el-button @click="handleAdd">Add</el-button>
  </el-container>
  <el-dialog v-model="addDialogVisible">
    <EditBookForm book-i-d="0"
                  :initial-book-form="currentBookItem"
    :on-submit="doAdd" :on-cancel="()=>{this.addDialogVisible=false}"></EditBookForm>
  </el-dialog>
  <el-dialog v-model="editDialogVisible">
    <EditBookForm book-i-d="0"
                  :initial-book-form="currentBookItem"
                  :on-submit="doEdit" :on-cancel="()=>{this.editDialogVisible=false}"></EditBookForm>
  </el-dialog>
</template>

<script>
import BookAPI from "../api/book_api";
import {reactive} from "vue";
import {ElMessage, ElMessageBox} from "element-plus";
import EditBookForm from "../components/EditBookForm.vue";
import {handleAxiosError} from "../api/handler";

export default {
  name: "BookView",
  components: {EditBookForm},
  data(){
    return {
      books:reactive([]),
      addDialogVisible:reactive(false),
      editDialogVisible:reactive(false),
      currentBookID:reactive(0),
      currentBookItem:reactive({name:'',author:''}),
    }
  },
  methods:{
    doGetBook(){
      BookAPI.listBook().then((response)=>{
        this.books=response.data.data;
      }).catch(handleAxiosError)
    },
    handleRent(i,row){
      const bookID=row.id;
      ElMessageBox.confirm(
          `Are you sure to rent book with id:${bookID}`, 'Warning',
          {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'warning',
          }).then(() => {
        this.doRent(bookID);
      }).catch((reason)=>{
        //ignore
      })
    },
    doRent(bookID){
      BookAPI.addRent({
        bookID:bookID,
      }).then((response)=>{
        ElMessage.success({
          message: `success`,
        })
      }).catch(handleAxiosError)
    },
    handleEdit(i,row){
      this.setCurrent(row);
      this.editDialogVisible=true;
    },
    setCurrent(row){
      this.currentBookID=row.id;
      this.currentBookItem=row;
    },
    clearCurrentAndDismissDialog(){
      this.currentBookID=0;
      this.currentBookItem={name:'',author:''};
      this.addDialogVisible=false;
      this.editDialogVisible=false;
    },
    doEdit(bookID,bookForm){
      BookAPI.editBook(bookID,bookForm).then((response)=>{
        ElMessage.success({
          message: `success`,
        })
        this.doGetBook();
        this.clearCurrentAndDismissDialog();
      }).catch(handleAxiosError)
    },

    handleAdd(){
      this.currentBookID=0;
      this.currentBookItem={name:'',author:''};
      this.addDialogVisible=true;
    },

    doAdd(bookID,bookForm){
      BookAPI.addBook(bookForm).then((response)=>{
        ElMessage.success({
          message: `success`,
        })
        this.doGetBook();
        this.clearCurrentAndDismissDialog();
      }).catch(handleAxiosError)
    },

    handleDelete(i,row){
      console.log(row);
      const bookID=row.id;
      ElMessageBox.confirm(
          `Are you sure to delete book with id:${bookID}`, 'Warning',
          {
            confirmButtonText: 'OK',
            cancelButtonText: 'Cancel',
            type: 'warning',
          }).then(() => {
            this.doDelete(bookID);
          }).catch((reason)=>{
            //ignore
          })
    },

    doDelete(bookID){
      BookAPI.deleteBook(bookID).then((response)=>{
        ElMessage.success({
          message: 'Delete completed',
        });
        this.doGetBook();
        this.clearCurrentAndDismissDialog();
      }).catch(handleAxiosError)
    }
  },
  mounted() {
    this.doGetBook();
  }
}
</script>

<style scoped>

</style>