import 'phaser'
import SeaScene from './scenes/seaScene'

const config = {
  type: Phaser.AUTO,
  scale: {
    mode: Phaser.Scale.FIT,
    autoCenter: Phaser.Scale.CENTER_BOTH,
    width: 800,
    height: 600
  },
  backgroundColor: '#125555',
  width: 800,
  height: 600,
  scene: SeaScene
}

const game = new Phaser.Game(config)
