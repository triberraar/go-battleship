<template>
  <div class="game">
    <h1>This is an game page</h1>
  </div>
</template>

<script>
export default {
  name: 'Game',
  data() {
    return {
      downloaded: false,
      gameInstance: null,
      containerId: 'game-container'
    }
  },
  async mounted() {
    const game = await import(/* webpackChunkName: "game" */ '@/games/battleships/battleships')
    this.downloaded = true
    this.$nextTick(() => {
      this.gameInstance = game.launchBattleships(this.containerId)
    })
  },
  destroyed() {
    this.gameInstance.destroy(true)
  }
}
</script>
