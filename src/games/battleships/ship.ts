import 'phaser'

class Ship extends Phaser.GameObjects.GameObject {
  constructor(
    scene: Phaser.Scene,
    type: string,
    private x: number,
    private y: number,
    private size: number,
    private vertical: boolean
  ) {
    super(scene, type)

    if (vertical) {
      // const ship = scene.add.image(24 + x * 48, this.size * 24 + y * 48, this.type)
      const sp = this.scene.add.sprite(
        24 + x * 48,
        this.size * 24 + y * 48,
        `${this.type}Destroyed`
      )
      sp.anims.load(`${this.type}Destroyed`)
      sp.anims.play(`${this.type}Destroyed`)
      sp.on('animationcomplete', this.destructionDone, this)
    } else {
      // const ship = scene.add.image(this.size * 24 + x * 48, 24 + y * 48, this.type).setAngle(90)
      const sp = this.scene.add
        .sprite(this.size * 24 + x * 48, 24 + y * 48, `${this.type}Destroyed`)
        .setAngle(90)
      sp.anims.load(`${this.type}Destroyed`)
      sp.anims.play(`${this.type}Destroyed`)
      sp.on('animationcomplete', this.destructionDone, this)
    }
  }

  destructionDone() {
    const p = this.scene.add.particles('smokeParticle')
    const emiterConfig = {
      speed: 50,
      blendMode: Phaser.BlendModes.ADD,
      scale: 0.1,
      frequency: 1,
      quantity: 1
    } as Phaser.Types.GameObjects.Particles.ParticleEmitterConfig
    for (let i = 0; i < this.size; i++) {
      if (this.vertical) {
        emiterConfig.x = this.x * 48 + 24
        emiterConfig.y = (this.y + i) * 48 + 24
        const e = p.createEmitter(emiterConfig)
        e.setDeathZone({
          type: 'onLeave',
          source: new Phaser.Geom.Circle(this.x * 48 + 24, (this.y + i) * 48 + 24, 24)
        })
      } else {
        emiterConfig.x = (this.x + i) * 48 + 24
        emiterConfig.y = this.y * 48 + 24
        const e = p.createEmitter(emiterConfig)
        e.setDeathZone({
          type: 'onLeave',
          source: new Phaser.Geom.Circle((this.x + i) * 48 + 24, this.y * 48 + 24, 24)
        })
      }
    }
  }
}

export class Ship1 extends Ship {
  constructor(scene: Phaser.Scene, x: number, y: number, vertical: boolean) {
    super(scene, 'ship1', x, y, 1, vertical)
  }
}

export class Ship2 extends Ship {
  constructor(scene: Phaser.Scene, x: number, y: number, vertical: boolean) {
    super(scene, 'ship2', x, y, 2, vertical)
  }
}

export class Ship3 extends Ship {
  constructor(scene: Phaser.Scene, x: number, y: number, vertical: boolean) {
    super(scene, 'ship3', x, y, 3, vertical)
  }
}

export class Ship4 extends Ship {
  constructor(scene: Phaser.Scene, x: number, y: number, vertical: boolean) {
    super(scene, 'ship4', x, y, 4, vertical)
  }
}

export class OpponentShip extends Phaser.GameObjects.GameObject {
  private destroyed: boolean

  private sprite: Phaser.GameObjects.Sprite

  constructor(
    scene: Phaser.Scene,
    type: string,
    private x: number,
    private y: number,
    private size: number
  ) {
    super(scene, type)
    this.destroyed = false

    this.sprite = this.scene.add
      .sprite(x, y, `ship${size}`)
      .setAngle(-90)
      .setOrigin(0, 0)
  }

  matches(size: number, destroyed: boolean) {
    return this.size === size && destroyed === this.destroyed
  }

  hide() {
    this.destroyed = true
    this.sprite.setVisible(false)
  }
}
