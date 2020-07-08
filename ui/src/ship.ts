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
      const ship = scene.add.image(24 + x * 48, this.size * 24 + y * 48, this.type)
    } else {
      const ship = scene.add.image(this.size * 24 + x * 48, 24 + y * 48, this.type).setAngle(90)
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
