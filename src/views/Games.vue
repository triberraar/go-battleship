<template>
  <div class="container">
    <div class="tile is-ancestor">
      <div class="tile is-parent">
        <article class="tile is-child box is-4" @click="playBattleships">
          <p class="title">Battleships</p>
          <figure class="image is-3by2">
            <img
              src="https://hips.hearstapps.com/pop.h-cdn.co/assets/17/20/2048x1335/gallery-1495140565-gettyimages-78402598.jpg?resize=480:320"
            />
          </figure>
        </article>
        <article class="tile is-child box is-4" @click="playRps">
          <p class="title">RPS</p>
        </article>
        <article class="tile is-child box is-4">
          <p class="title">Other game</p>
        </article>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import User from '@/store/modules/user'

@Component({})
export default class Games extends Vue {
  userModule = getModule(User, this.$store)

  playBattleships(): void {
    if (this.userModule.loggedIn) {
      this.$router.push({ name: 'Battleships' })
    } else {
      this.$buefy.notification.open({
        message: 'You need to be logged in to play this game',
        type: 'is-danger',
        hasIcon: true,
      })
    }
  }

  playRps(): void {
    if (this.userModule.loggedIn) {
      this.$router.push({ name: 'rps' })
    } else {
      this.$buefy.notification.open({
        message: 'You need to be logged in to play this game',
        type: 'is-danger',
        hasIcon: true,
      })
    }
  }
}
</script>
