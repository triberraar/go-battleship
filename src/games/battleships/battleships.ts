import 'phaser'

import RexUIPlugin from 'phaser3-rex-plugins/templates/ui/ui-plugin'
import SeaScene from './scenes/seaScene'
import MenuScene from './scenes/menuScene'

function launchBattleships(containerId: string): Phaser.Game {
  const g = new Phaser.Game({
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

  g.events.on('destroy', () => {
    const ms = g.scene.getScene('menuScene') as MenuScene
    ms.clear()
    const ss = g.scene.getScene('seaScene') as SeaScene
    ss.clear()
  })
  return g
}

export default launchBattleships
export { launchBattleships }
