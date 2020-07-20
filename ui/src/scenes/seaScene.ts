import BoardManager from '../board'
import CommunicationManager from '../communication'
import FeedbackText from '../feedbackText'

export default class SeaScene extends Phaser.Scene {
  private boardManager: BoardManager
  private communicationManager: CommunicationManager

  constructor() {
    super('seaScene')
  }

  preload() {
    this.load.spritesheet('tile', 'assets/tile.png', { frameWidth: 48 })
    this.load.spritesheet('miss', 'assets/miss.png', { frameWidth: 50, frameHeight: 50 })
    this.load.spritesheet('hit', 'assets/fireOponent.png', { frameWidth: 70, frameHeight: 70 })
    this.load.spritesheet('ship1Destroyed', 'assets/fireShipOp_1.png', { frameWidth: 54, frameHeight: 45 })
    this.load.spritesheet('ship2Destroyed', 'assets/fireShipOp_2.png', { frameWidth: 53, frameHeight: 85 })
    this.load.spritesheet('ship3Destroyed', 'assets/fireShipOp_3.png', { frameWidth: 69, frameHeight: 138 })
    this.load.spritesheet('ship4Destroyed', 'assets/fireShipOp_4.png', { frameWidth: 69, frameHeight: 180 })
    this.load.spritesheet('opponentShip1Destroyed', 'assets/fireShip_1.png', { frameWidth: 54, frameHeight: 45 })
    this.load.spritesheet('opponentShip2Destroyed', 'assets/fireShip_2.png', { frameWidth: 53, frameHeight: 85 })
    this.load.spritesheet('opponentShip3Destroyed', 'assets/fireShip_3.png', { frameWidth: 69, frameHeight: 138 })
    this.load.spritesheet('opponentShip4Destroyed', 'assets/fireShip_4.png', { frameWidth: 69, frameHeight: 180 })
    this.load.atlas('fireworks', 'assets/flares.png', 'assets/flares.json')

    this.load.image('ship1', 'assets/1_ship.png')
    this.load.image('ship2', 'assets/2_ship.png')
    this.load.image('ship3', 'assets/3_ship.png')
    this.load.image('ship4', 'assets/4_ship.png')
    this.load.image('smokeParticle', 'assets/smoke.png')
  }

  create() {
    this.anims.create({ key: 'miss', frames: this.anims.generateFrameNumbers('miss', {}), repeat: 0, frameRate: 16 })
    this.anims.create({ key: 'hit', frames: this.anims.generateFrameNumbers('hit', {}), repeat: 0, frameRate: 16 })
    this.anims.create({
      key: 'ship1Destroyed',
      frames: this.anims.generateFrameNumbers('ship1Destroyed', {}),
      repeat: 0,
      frameRate: 16
    })
    this.anims.create({
      key: 'ship2Destroyed',
      frames: this.anims.generateFrameNumbers('ship2Destroyed', {}),
      repeat: 0,
      frameRate: 16
    })
    this.anims.create({
      key: 'ship3Destroyed',
      frames: this.anims.generateFrameNumbers('ship3Destroyed', {}),
      repeat: 0,
      frameRate: 16
    })
    this.anims.create({
      key: 'ship4Destroyed',
      frames: this.anims.generateFrameNumbers('ship4Destroyed', {}),
      repeat: 0,
      frameRate: 16
    })

    this.communicationManager = new CommunicationManager()
    this.boardManager = new BoardManager(this, 10, 10, this.communicationManager)
    this.communicationManager.setBoardManager(this.boardManager)
    this.communicationManager.setFeedbackText(new FeedbackText(this))
    this.communicationManager.play()
  }

  init(data) {
    console.log('init seascene')
  }
}
