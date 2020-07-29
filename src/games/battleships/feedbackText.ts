import 'phaser'

export default class FeedbackText extends Phaser.GameObjects.GameObject {
  private text: Phaser.GameObjects.Text

  private countDownTimer: number

  constructor(scene: Phaser.Scene) {
    super(scene, 'feedbackText')
    this.text = scene.add.text(500, 500, 'Awaiting players')
  }

  setText(m: string) {
    clearInterval(this.countDownTimer)
    this.text.setText(m)
  }

  setCountDownText(m: string, i: number) {
    let counter = i
    clearInterval(this.countDownTimer)
    this.countDownTimer = setInterval(() => {
      counter -= 1
      if (counter === 0) {
        clearInterval(this.countDownTimer)
      } else {
        this.text.setText(`${m} (${counter} secs)`)
      }
    }, 1000)
  }
}
