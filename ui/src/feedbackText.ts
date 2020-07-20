import 'phaser'

export default class FeedbackText extends Phaser.GameObjects.GameObject {
  private text: Phaser.GameObjects.Text

  constructor(scene: Phaser.Scene) {
    super(scene, 'feedbackText')
    this.text = scene.add.text(500, 500, 'Awaiting players')
  }

  setText(m: string) {
    this.text.setText(m)
  }
}
