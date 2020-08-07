import 'phaser'
import CommunicationManager from './communication' // eslint-disable-line
import { Ship2, Ship1, Ship3, Ship4, OpponentShip } from './ship'

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

  /* eslint-disable no-new */
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
      default: {
        console.error('unknown ship size')
        break
      }
    }
  }
  /* eslint-enable no-new */

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

  private opponentShips: OpponentShip[] = []

  private fireworksEmitter1: Phaser.GameObjects.Particles.ParticleEmitter

  private fireworksEmitter2: Phaser.GameObjects.Particles.ParticleEmitter

  constructor(
    private scene: Phaser.Scene,
    private x: number,
    private y: number,
    private communicationManager: CommunicationManager
  ) {
    for (let i = 0; i < 10; i++) {
      this.board[i] = []
      for (let j = 0; j < 10; j++) {
        this.board[i][j] = new SeaTile(this.scene, i, j, this.communicationManager)
      }
    }
  }

  miss(x: number, y: number) {
    this.board[x][y].miss()
  }

  hit(x: number, y: number) {
    this.board[x][y].hit()
  }

  destoryShip(x: number, y: number, size: number, vertical: boolean) {
    for (let i = 0; i < size; i++) {
      if (vertical) {
        this.board[x][y + i].hideHit()
      } else {
        this.board[x + i][y].hideHit()
      }
    }
    this.board[x][y].destoryShip(size, vertical)
    const os = this.opponentShips.find(s => s.matches(size, false))
    if (os) {
      os.hide()
    }
  }

  victory() {
    const particles = this.scene.add.particles('fireworks')
    this.fireworksEmitter1 = particles.createEmitter({
      frame: ['red', 'green'],
      lifespan: 4000,
      angle: { min: -0, max: 360 },
      speed: { min: 0, max: 300 },
      scale: { start: 0.6, end: 0 },
      gravityY: 300,
      bounce: 0.9,

      collideTop: false,
      collideBottom: false,
      blendMode: 'ADD',
      on: false
    })
    this.fireworksEmitter2 = particles.createEmitter({
      frame: ['yellow', 'blue', 'white'],
      lifespan: 4000,
      angle: { min: -0, max: 360 },
      speed: { min: 0, max: 300 },
      scale: { start: 0.6, end: 0 },
      gravityY: 300,
      bounce: 0.9,

      collideTop: false,
      collideBottom: false,
      blendMode: 'ADD',
      on: false
    })
    for (let i = 0; i < 3; i++) {
      setTimeout(
        () => this.fireworksEmitter1.explode(
          150,
          Math.floor(Math.random() * (750 - 50)) + 50,
          Math.floor(Math.random() * (550 - 50)) + 50
        ),
        Math.random() * (5000 - 1000) + 1000
      )
      setTimeout(
        () => this.fireworksEmitter2.explode(
          150,
          Math.floor(Math.random() * (750 - 50)) + 50,
          Math.floor(Math.random() * (550 - 50)) + 50
        ),
        Math.random() * (5000 - 1000) + 1000
      )
    }
    this.fireworksEmitter1.explode(
      150,
      Math.floor(Math.random() * (750 - 50)) + 50,
      Math.floor(Math.random() * (550 - 50)) + 50
    )
    this.fireworksEmitter2.explode(
      150,
      Math.floor(Math.random() * (750 - 50)) + 50,
      Math.floor(Math.random() * (550 - 50)) + 50
    )
    setTimeout(() => this.backToMenu(), 6000)
  }

  loss() {
    setTimeout(() => this.backToMenu(), 6000)
  }

  ships(shipSizes: number[]) {
    shipSizes.sort()
    const os = shipSizes.map(
      (s, i) => new OpponentShip(this.scene, 'opponentship', 500, 50 + i * 50, s)
    )
    this.opponentShips = os
  }

  backToMenu() {
    this.scene.scene.stop('seaScene')
    this.scene.scene.start('menuScene')
  }
}
