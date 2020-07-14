import CommunicationManager from '../communication'
import SeaScene from './seaScene'

export default class MenuScene extends Phaser.Scene {
  private playButton: any
  private communicationManager: CommunicationManager

  constructor() {
    super('menuScene')
  }

  init() {
    this.communicationManager = new CommunicationManager()
  }
  create() {
    this.playButton = this.createButton('Play')
    // @ts-ignore
    var buttons = this.rexUI.add
      .buttons({
        x: 400,
        y: 300,
        width: 300,
        orientation: 'x',
        buttons: [this.playButton],

        expand: true
      })
      .layout()

    buttons.on('button.click', this.buttonClicked, this)
  }

  buttonClicked(button: any) {
    if (button === this.playButton) {
      this.scene.start('seaScene', { communicationManager: this.communicationManager })
    }
  }

  createButton(text: string) {
    // @ts-ignore
    return this.rexUI.add.label({
      width: 40,
      height: 40,
      // @ts-ignore
      background: this.rexUI.add.roundRectangle(0, 0, 0, 0, 20, 0x7b5e57),
      text: this.add.text(0, 0, text, {
        fontSize: 18
      }),
      space: {
        left: 10,
        right: 10
      },
      align: 'center'
    })
  }
}
