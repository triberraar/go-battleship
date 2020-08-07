<template>
  <div class="tile is-ancestor">
    <div class="tile is-3" v-if="!gameLoading">
      <div :id="containerId"></div>
    </div>
    <div class="tile is-3" v-else>
      <b-progress />
    </div>
    <div class="tile is-3" v-if="!gameLoading">
      <b-table :data="battleshipModule.statistics" :columns="columns" />
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import User from '@/store/modules/user'
import Battleship from '@/store/modules/battleship'

@Component({})
export default class Battleships extends Vue {
  userModule = getModule(User, this.$store)
  battleshipModule = getModule(Battleship, this.$store)

  gameInstance: Phaser.Game
  gameLoading = true

  containerId = 'battleshipsContainer'

  async mounted() {
    if (!this.userModule.loggedIn) {
      this.$buefy.notification.open({
        message: 'You need to be logged in to play this game',
        type: 'is-danger',
        hasIcon: true
      })
      this.$router.push({ name: 'Games' })
      return
    }
    const game = await import(/* webpackChunkName: "game" */ '@/games/battleships/battleships')
    this.gameLoading = false
    this.$nextTick(() => {
      this.gameInstance = game.launchBattleships(this.containerId)
    })
  }

  destroyed() {
    if (this.gameInstance) {
      this.gameInstance.destroy(true)
    }
  }

  @Watch('battleshipModule.connected')
  connectionChanged(n: string) {
    if (n === 'RECONNECTING') {
      this.$buefy.toast.open({
        message: 'Disconnected, trying to reconnect',
        type: 'is-warning'
      })
    } else if (n === 'CONNECTED') {
      this.$buefy.toast.open({
        message: 'Connected',
        type: 'is-success'
      })
    } else if (n === 'FAILED') {
      this.$buefy.toast.open({
        message: 'Failed to connect, curling up and dieing',
        type: 'is-danger'
      })
    }
  }

  tddd = [
    {
      username: 'sdf',
      hits: 1,
      misses: 10,
      destroyed: 5
    },
    {
      username: 'sdsdff',
      hits: 1,
      misses: 10,
      destroyed: 5
    }
  ]

  columns = [
    {
      field: 'username',
      label: 'username'
    },
    {
      field: 'hits',
      label: 'Hits'
    },
    {
      field: 'misses',
      label: 'Misses'
    },
    {
      field: 'destroyed',
      label: 'Ships destroyed'
    }
  ]
}
</script>
