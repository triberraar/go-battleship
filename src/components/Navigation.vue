<template>
  <b-navbar>
    <template slot="brand">
      <b-navbar-item>
        <img
          src="https://storage.googleapis.com/hexigames-logos/logo_dark_low_res.png"
          alt="Hexigames"
        />
      </b-navbar-item>
    </template>
    <template slot="start">
      <b-navbar-item class="is-dark" tag="router-link" :to="{ name: 'Home' }">Home</b-navbar-item>
      <b-navbar-item tag="router-link" :to="{ name: 'Accounting' }">Accounting</b-navbar-item>
      <b-navbar-item tag="router-link" :to="{ name: 'Games' }">Games</b-navbar-item>
    </template>

    <template slot="end">
      <b-navbar-item tag="router-link" :to="{ name: 'Login' }" v-if="!userModule.loggedIn">
        <div class="buttons">
          <a class="button is-light">Log in</a>
        </div>
      </b-navbar-item>
      <b-navbar-item v-else>
        <div class="buttons">
          <a class="button is-light" @click="logOut">Log out</a>
        </div>
      </b-navbar-item>
    </template>
  </b-navbar>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import User from '@/store/modules/user'

@Component({})
export default class Navigation extends Vue {
  userModule = getModule(User, this.$store)

  logOut(): void {
    this.userModule.logOut()
    this.$router.push({ name: 'Home' })
  }
}
</script>
