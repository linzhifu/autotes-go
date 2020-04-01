// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import ElementUI from 'element-ui'
import 'element-ui/lib/theme-chalk/index.css'
import axios from 'axios'
import './assets/css/main.css'
import './assets/icon/iconfont.css'

Vue.use(ElementUI)

Vue.config.productionTip = false
// 测试分支
// || window.location.origin.indexOf('10.2.40.232') !== -1
if (window.location.host === '127.0.0.1:8080' || window.location.host === 'localhost:8080' || window.location.origin.indexOf('10.2.40.232') !== -1) {
  Vue.prototype.url = 'http://localhost:8000'
} else {
  Vue.prototype.url = window.location.origin
}

// Vue.prototype.url = 'http://172.16.9.88/'
// Vue.prototype.url = 'http://127.0.0.1:8000/'
// Vue.prototype.url = 'http://47.106.111.62/'
// Vue.prototype.url = window.location.origin
Vue.prototype.axios = axios
Vue.prototype.storage = window.localStorage
Vue.prototype.path = ''
// axios.defaults.withCredentials = true

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
