import AxiosInstance from "./axios";

const BookAPI={
    addBook:function (data){
        return AxiosInstance.post('/task3/book/',data)
    },
    deleteBook:function (bookID){
        return AxiosInstance.delete(`/task3/book/${bookID}`)
    },
    editBook:function (bookID,data){
        return AxiosInstance.put(`/task3/book/${bookID}`,data)
    },
    listBook:function (){
        return AxiosInstance.get(`/task3/book/`)
    },
    listRent:function (){
        return AxiosInstance.get('/task3/book/rent/user/self')
    },
    addRent:function(rentForm){
        return AxiosInstance.post('/task3/book/rent/',rentForm)
    },
    returnRent:function (rentID){
        return AxiosInstance.put(`/task3/book/rent/return/${rentID}`)
    }
}

export default BookAPI;