import BoardManager from '../board'
import CommunicationManager from '../communication'

export default class SeaScene extends Phaser.Scene {
  private boardManager: BoardManager
  constructor() {
    super('sea')
  }

  preload() {
    this.load.spritesheet('tile', 'assets/tile.png', { frameWidth: 48 })
    this.load.spritesheet('miss', 'assets/miss.png', { frameWidth: 50 })
    this.load.spritesheet('hit', 'assets/fireOponent.png', { frameWidth: 70 })

    this.load.image('ship1', 'assets/1_ship.png')
    this.load.image('ship2', 'assets/2_ship.png')
    this.load.image('ship3', 'assets/3_ship.png')
    this.load.image('ship4', 'assets/4_ship.png')
  }

  create() {
    this.anims.create({ key: 'miss', frames: this.anims.generateFrameNumbers('miss', {}), repeat: 0, frameRate: 16 })
    this.anims.create({ key: 'hit', frames: this.anims.generateFrameNumbers('hit', {}), repeat: 0, frameRate: 16 })

    const communicationManager = new CommunicationManager()
    this.boardManager = new BoardManager(this, 10, 10, communicationManager)
    communicationManager.setBoardManager(this.boardManager)
  }
}
