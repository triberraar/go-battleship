import BoardManager from '../board'
import CommunicationManager from '../communication'
import FeedbackText from '../feedbackText'

import tile from '@/assets/tile.png'
import miss from '@/assets/miss.png'
import hit from '@/assets/fireOponent.png'
import ship1Destroyed from '@/assets/fireShipOp_1.png'
import ship2Destroyed from '@/assets/fireShipOp_2.png'
import ship3Destroyed from '@/assets/fireShipOp_3.png'
import ship4Destroyed from '@/assets/fireShipOp_4.png'
import opponentShip1Destroyed from '@/assets/fireShip_1.png'
import opponentShip2Destroyed from '@/assets/fireShip_2.png'
import opponentShip3Destroyed from '@/assets/fireShip_3.png'
import opponentShip4Destroyed from '@/assets/fireShip_4.png'
import ship1 from '@/assets/1_ship.png'
import ship2 from '@/assets/2_ship.png'
import ship3 from '@/assets/3_ship.png'
import ship4 from '@/assets/4_ship.png'
import smokeParticle from '@/assets/smoke.png'
import flaresParticles from '@/assets/flares'
import flares from '@/assets/flares.png'

export default class SeaScene extends Phaser.Scene {
  private boardManager: BoardManager

  private communicationManager: CommunicationManager

  constructor() {
    super('seaScene')
  }

  preload() {
    this.load.spritesheet('tile', tile, { frameWidth: 48 })
    this.load.spritesheet('miss', miss, { frameWidth: 50, frameHeight: 50 })
    this.load.spritesheet('hit', hit, { frameWidth: 70, frameHeight: 70 })
    this.load.spritesheet('ship1Destroyed', ship1Destroyed, {
      frameWidth: 54,
      frameHeight: 45
    })
    this.load.spritesheet('ship2Destroyed', ship2Destroyed, {
      frameWidth: 53,
      frameHeight: 85
    })
    this.load.spritesheet('ship3Destroyed', ship3Destroyed, {
      frameWidth: 69,
      frameHeight: 138
    })
    this.load.spritesheet('ship4Destroyed', ship4Destroyed, {
      frameWidth: 69,
      frameHeight: 180
    })
    this.load.spritesheet('opponentShip1Destroyed', opponentShip1Destroyed, {
      frameWidth: 54,
      frameHeight: 45
    })
    this.load.spritesheet('opponentShip2Destroyed', opponentShip2Destroyed, {
      frameWidth: 53,
      frameHeight: 85
    })
    this.load.spritesheet('opponentShip3Destroyed', opponentShip3Destroyed, {
      frameWidth: 69,
      frameHeight: 138
    })
    this.load.spritesheet('opponentShip4Destroyed', opponentShip4Destroyed, {
      frameWidth: 69,
      frameHeight: 180
    })
    this.load.atlas('fireworks', flares, flaresParticles)
    this.load.image('ship1', ship1)
    this.load.image('ship2', ship2)
    this.load.image('ship3', ship3)
    this.load.image('ship4', ship4)
    this.load.image('smokeParticle', smokeParticle)
  }

  create() {
    this.anims.create({
      key: 'miss',
      frames: this.anims.generateFrameNumbers('miss', {}),
      repeat: 0,
      frameRate: 16
    })
    this.anims.create({
      key: 'hit',
      frames: this.anims.generateFrameNumbers('hit', {}),
      repeat: 0,
      frameRate: 16
    })
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
  }

  clear() {
    if (this.communicationManager) {
      this.communicationManager.close()
    }
  }
}
