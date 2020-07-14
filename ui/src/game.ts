import 'phaser'
import SeaScene from './scenes/seaScene'
import MenuScene from './scenes/menuScene'
import RexUIPlugin from 'phaser3-rex-plugins/templates/ui/ui-plugin.js'

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
  scene: [MenuScene, SeaScene],
  plugins: {
    scene: [
      {
        key: 'rexUI',
        plugin: RexUIPlugin,
        mapping: 'rexUI'
      }
    ]
  }
}

const game = new Phaser.Game(config)
