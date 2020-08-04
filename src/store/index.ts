import Vue from 'vue'
import Vuex from 'vuex'
import User from './modules/user'
import Battleship from './modules/battleship'
import { config } from 'vuex-module-decorators'

config.rawError = true

Vue.use(Vuex)

export default new Vuex.Store({
  modules: { User, Battleship }
})
