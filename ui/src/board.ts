import 'phaser'
import CommunicationManager from './communication'

class SeaTile extends Phaser.GameObjects.GameObject {
  constructor(
    scene: Phaser.Scene,
    private x: number,
    private y: number,
    private communicationManager: CommunicationManager
  ) {
    super(scene, 'tile')
    scene.add
      .image(x * 48, y * 48, 'tile', 0)
      .setOrigin(0, 0)
      .setInteractive()
      .on('pointerdown', this.handleClick, this)
  }

  handleClick() {
    // if (this.x % 2 == 0) {
    //   this.communicationManager.onHit(this.x, this.y)
    // } else {
    //   this.communicationManager.onMiss(this.x, this.y)
    // }
    this.communicationManager.fire(this.x, this.y)
  }

  miss() {
    const missSprite = this.scene.add.sprite(this.x * 48, this.y * 48, 'miss').setOrigin(0, 0)
    missSprite.anims.load('miss')
    missSprite.anims.play('miss')
  }

  hit() {
    const hitSprite = this.scene.add.sprite(this.x * 48, this.y * 48, 'hit').setOrigin(0.15, 0.15)
    hitSprite.anims.load('hit')
    hitSprite.anims.play('hit')
  }
}

export default class BoardManager {
  private board: SeaTile[][] = []
  miss(x: number, y: number) {
    this.board[x][y].miss()
  }
  hit(x: number, y: number) {
    this.board[x][y].hit()
  }
  constructor(
    private scene: Phaser.Scene,
    private x: number,
    private y: number,
    private communicationManager: CommunicationManager
  ) {
    for (var i = 0; i < 10; i++) {
      this.board[i] = []
      for (var j = 0; j < 10; j++) {
        this.board[i][j] = new SeaTile(this.scene, i, j, communicationManager)
      }
    }
  }
}
