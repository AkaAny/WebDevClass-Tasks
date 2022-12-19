import {ElMessageBox} from "element-plus";

export function handleAxiosError(reason){
    const response=reason.response;
    const respData=response.data;
    ElMessageBox.alert(respData.msg, 'Error', {
        // if you want to disable its autofocus
        // autofocus: false,
        confirmButtonText: 'OK',
    })
}