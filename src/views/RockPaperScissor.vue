<template>
  <div class="about">
    <h1>This is rock papaer scissor</h1>
    <b-button type="is-primary" v-if="!url" @click="play">Play</b-button>
    <div class="buttons" v-if="url">
      <b-button type="is-primary" @click="rock">Rock</b-button>
      <b-button type="is-primary" @click="paper">Paper</b-button>
      <b-button type="is-primary" @click="scissor">Scissor</b-button>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator'
import { getModule } from 'vuex-module-decorators'
import User from '@/store/modules/user'

interface Result {
  Result: string
}

@Component({})
export default class RPSGame extends Vue {
  userModule = getModule(User, this.$store)
  url = null

  mounted() {
    if (!this.userModule.loggedIn) {
      this.$buefy.notification.open({
        message: 'You need to be logged in to play this game',
        type: 'is-danger',
        hasIcon: true,
      })
      this.$router.push({ name: 'Games' })
    }
  }

  play(): void {
    console.log('play')
    fetch(`./rps/play?username=${this.userModule.username}`).then(
      (response) => {
        response.json().then((data) => {
          console.log(data)
          if (data.URL) {
            this.url = data.URL
          } else {
            console.log('no url')
          }
        })
      }
    )
  }

  doResult(data: Result): void {
    if (data.Result === 'c') {
      this.$buefy.notification.open({
        message: 'Game is cancelled',
        type: 'is-danger',
        hasIcon: true,
      })
    } else if (data.Result === 'l') {
      this.$buefy.notification.open({
        message: 'Loser',
        type: 'is-danger',
        hasIcon: true,
      })
    } else if (data.Result === 'd') {
      this.$buefy.notification.open({
        message: 'Draw',
        type: 'is-warning',
        hasIcon: true,
      })
    } else if (data.Result === 'w') {
      this.$buefy.notification.open({
        message: 'Win',
        type: 'is-success',
        hasIcon: true,
      })
    }
    this.url = null
  }

  rock(): void {
    fetch(`http://${this.url}/play`, {
      method: 'POST',
      body: JSON.stringify({ username: this.userModule.username, move: 'r' }),
    }).then((response) => {
      response.json().then((data) => {
        this.doResult(data)
      })
    })
  }

  paper(): void {
    fetch(`http://${this.url}/play`, {
      method: 'POST',
      body: JSON.stringify({ username: this.userModule.username, move: 'p' }),
    }).then((response) => {
      response.json().then((data) => {
        this.doResult(data)
      })
    })
  }

  scissor(): void {
    fetch(`http://${this.url}/play`, {
      method: 'POST',
      body: JSON.stringify({ username: this.userModule.username, move: 's' }),
    }).then((response) => {
      response.json().then((data) => {
        this.doResult(data)
      })
    })
  }
}
</script>
