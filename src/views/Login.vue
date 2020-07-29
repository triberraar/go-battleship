<template>
  <div class="container">
    <section>
      <b-field label="Username">
        <b-input v-model="username" placeholder="Ove or something" rounded />
      </b-field>
      <div class="buttons">
        <button class="button is-success" @click="login">
          <span class="icon is-small">
            <i class="fas fa-check" />
          </span>
          <span>Log In</span>
        </button>
      </div>
    </section>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import User from '@/store/modules/user'

@Component({})
export default class Navigation extends Vue {
  userModule = getModule(User, this.$store)

  username = ''

  login(): void {
    this.userModule.logIn(this.username)
    this.$buefy.toast.open({
      message: `Welcome ${this.username}`,
      type: 'is-success'
    })
    this.$router.push({ name: 'Home' })
  }
}
</script>
