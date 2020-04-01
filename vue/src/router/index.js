import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/components/Login'
import Home from '@/components/Home'
import User from '@/components/User'
import Project from '@/components/Project'
import Test from '@/components/Test'
import Log from '@/components/Log'
import Report from '@/components/Report'
import Src from '@/components/Src'
import Mpcloud from '@/components/template/Mpcloud'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      redirect: '/login/'
    },
    {
      path: '/login/',
      name: 'Login',
      component: Login
    },
    {
      path: '/home/',
      component: Home,
      children: [
        {
          path: '',
          redirect: 'user'
        },
        {
          path: 'user',
          component: User
        },
        {
          path: 'project',
          component: Project
        },
        {
          path: 'test',
          component: Test
        },
        {
          path: 'log/',
          component: Log
        },
        {
          path: 'report',
          component: Report
        },
        {
          path: 'src',
          component: Src
        },
        {
          path: 'mpcloud',
          component: Mpcloud
        }
      ]
    }
  ]
})
