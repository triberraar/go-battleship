export default class MenuScene extends Phaser.Scene {
  private playButton: any // eslint-disable-line

  constructor() {
    super('menuScene')
  }

  create() {
    this.playButton = this.createButton('Play')
    // @ts-ignore
    const buttons = this.rexUI.add
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

  // eslint-disable-next-line
  buttonClicked(button: any) {
    if (button === this.playButton) {
      this.scene.start('seaScene')
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

  clear() {
    this.plugins.removeScenePlugin('rexUI')
  }
}
