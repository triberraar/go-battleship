import 'phaser'
import CommunicationManager from './communication'
import { Ship2, Ship1, Ship3, Ship4 } from './ship'

class SeaTile extends Phaser.GameObjects.GameObject {
  private hitSprite: Phaser.GameObjects.Sprite

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
    this.communicationManager.fire(this.x, this.y)
  }

  destoryShip(size: number, vertical: boolean) {
    switch (size) {
      case 1: {
        new Ship1(this.scene, this.x, this.y, vertical)
        break
      }
      case 2: {
        new Ship2(this.scene, this.x, this.y, vertical)
        break
      }
      case 3: {
        new Ship3(this.scene, this.x, this.y, vertical)
        break
      }
      case 4: {
        new Ship4(this.scene, this.x, this.y, vertical)
        break
      }
    }
  }

  miss() {
    const missSprite = this.scene.add.sprite(this.x * 48, this.y * 48, 'miss').setOrigin(0, 0)
    missSprite.anims.load('miss')
    missSprite.anims.play('miss')
  }

  hit() {
    this.hitSprite = this.scene.add.sprite(this.x * 48, this.y * 48, 'hit').setOrigin(0.15, 0.15)
    this.hitSprite.anims.load('hit')
    this.hitSprite.anims.play('hit')
  }

  hideHit() {
    if (this.hitSprite) this.hitSprite.setVisible(false)
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
  destoryShip(x: number, y: number, size: number, vertical: boolean) {
    for (var i = 0; i < size; i++) {
      if (vertical) {
        this.board[x][y + i].hideHit()
      } else {
        this.board[x + i][y].hideHit()
      }
    }
    this.board[x][y].destoryShip(size, vertical)
  }
  constructor(
    private scene: Phaser.Scene,
    private x: number,
    private y: number,
    communicationManager: CommunicationManager
  ) {
    for (var i = 0; i < 10; i++) {
      this.board[i] = []
      for (var j = 0; j < 10; j++) {
        this.board[i][j] = new SeaTile(this.scene, i, j, communicationManager)
      }
    }
  }
}
