import Vue from 'vue'
import VueRouter, { RouteConfig } from 'vue-router'
import Home from '../views/Home.vue'
import Games from '../views/Games.vue'
import Login from '../views/Login.vue'
import Battleships from '../views/Battleships.vue'
import RockPaperScissor from '../views/RockPaperScissor.vue'

Vue.use(VueRouter)

const routes: Array<RouteConfig> = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/accounting',
    name: 'Accounting',
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () => import(/* webpackChunkName: "about" */ '../views/About.vue')
  },
  {
    path: '/games',
    name: 'Games',
    component: Games
  },
  {
    path: '/bs',
    name: 'Battleships',
    component: Battleships
  },
  {
    path: '/rps',
    name: 'rps',
    component: RockPaperScissor
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

export default router
