<template>
  <div :id="containerId"></div>
</template>

<script>
export default {
  name: 'Games',
  data() {
    return {
      downloaded: false,
      gameInstance: null,
      containerId: 'battleshipsContainer'
    }
  },
  async mounted() {
    const game = await import(
      /* webpackChunkName: "game" */ '@/games/battleships/battleships'
    )
    this.$nextTick(() => {
      this.gameInstance = game.launchBattleships(this.containerId)
    })
  },
  destroyed() {
    this.gameInstance.destroy(true)
  }
}
</script>
