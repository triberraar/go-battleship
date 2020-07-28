import 'phaser';

import libs from '@/assets/libs.png';
import lobo from '@/assets/phaser3-logo.png';

class Demo extends Phaser.Scene {
  constructor() {
    super('demo');
  }

  preload() {
    this.load.image('logo', lobo);
    this.load.image('libs', libs);
  }

  create() {
    // this.add.shader('Plasma', 0, 412, 800, 172).setOrigin(0);

    this.add.image(400, 300, 'libs');

    this.add.image(400, 70, 'logo');
  }
}

// const config = {
//   type: Phaser.AUTO,
//   backgroundColor: '#125555',
//   width: 800,
//   height: 600,
//   scene: Demo
// };

// const game = new Phaser.Game(config);

function launch(containerId: string) {
  return new Phaser.Game({
    type: Phaser.AUTO,
    width: 800,
    height: 600,
    parent: containerId,
    physics: {
      default: 'arcade',
      arcade: {
        gravity: { y: 300 },
        debug: false
      }
    },
    scene: Demo
  });
}

export default launch;
export { launch };
