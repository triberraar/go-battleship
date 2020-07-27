import 'phaser'

export default class FeedbackText extends Phaser.GameObjects.GameObject {
  private text: Phaser.GameObjects.Text
  private countDownTimer: NodeJS.Timer

  constructor(scene: Phaser.Scene) {
    super(scene, 'feedbackText')
    this.text = scene.add.text(500, 500, 'Awaiting players')
  }

  setText(m: string) {
    clearInterval(this.countDownTimer)
    this.text.setText(m)
  }

  setCountDownText(m: string, i: number) {
    clearInterval(this.countDownTimer)
    this.countDownTimer = setInterval(() => {
      console.log(i)
      i = i - 1
      if (i === 0) {
        clearInterval(this.countDownTimer)
      } else {
        this.text.setText(`${m} (${i} secs)`)
      }
    }, 1000)
  }
}
