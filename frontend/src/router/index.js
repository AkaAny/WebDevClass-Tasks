import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path:'/bookInfo',
      name:'bookInfo',
      component: ()=> import('../views/BookInfoView.vue')
    },
    {
      path:'/bookRent',
      name:'bookRent',
      component: ()=> import('../views/BookRentView.vue')
    },
    {
      path:'/userRegister',
      name:'userRegister',
      component: ()=> import('../views/UserRegisterView.vue')
    },
    {
      path:'/userActive',
      name:'userActive',
      component: ()=> import('../views/UserActiveView.vue')
    },
    {
      path:'/userLogin',
      name:'userLogin',
      component: ()=> import('../views/UserLoginView.vue')
    },
    {
      path:'/userInfo',
      name:'userInfo',
      component: ()=> import('../views/UserInfoView.vue')
    }
  ]
})

export default router
