import 'phaser'

import RexUIPlugin from 'phaser3-rex-plugins/templates/ui/ui-plugin.js'
import SeaScene from './scenes/seaScene'
import MenuScene from './scenes/menuScene'

function launch(containerId: string) {
  return new Phaser.Game({
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
    },
    parent: containerId
  })
}

export default launch
export { launch }
