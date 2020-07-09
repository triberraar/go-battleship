import BoardManager from './board'

interface Coordinate {
  x: number
  y: number
}

interface HitMessage {
  type: string
  coordinate: Coordinate
}

interface missMessage {
  type: string
  coordinate: Coordinate
}

interface shipDestroyedMessage {
  type: string
  coordinate: Coordinate
  shipSize: number
  vertical: boolean
}

export default class CommunicationManager {
  private boardManager: BoardManager
  private ws: WebSocket
  constructor() {
    const loc = window.location
    var wsUri = ''
    if (loc.protocol === 'https:') {
      wsUri = 'wss:'
    } else {
      wsUri = 'ws:'
    }
    wsUri += '//' + loc.host + '/battleship'
    console.log(wsUri)
    if (loc.host.startsWith('localhost')) {
      wsUri = 'ws://localhost:10002/battleship'
    }

    this.ws = new WebSocket(wsUri)

    this.ws.onopen = this.onopen
    this.ws.onmessage = this.onmessage
    this.ws.onerror = this.onerror
  }

  onopen = () => {}

  onmessage = (ev: MessageEvent) => {
    const m = JSON.parse(ev.data)
    switch (m.type) {
      case 'HIT': {
        this.onHit(m)
        break
      }
      case 'MISS': {
        this.onMiss(m)
        break
      }
      case 'SHIP_DESTROYED': {
        this.onShipDestroyed(m)
        break
      }
    }
  }

  onerror = (ev: Event) => {
    console.log(`error ${ev}`)
  }

  setBoardManager(boardManager: BoardManager) {
    this.boardManager = boardManager
  }

  fire(x: number, y: number) {
    this.ws.send(JSON.stringify({ type: 'FIRE', coordinate: { x, y } }))
  }

  onHit(m: HitMessage) {
    this.boardManager.hit(m.coordinate.x, m.coordinate.y)
  }

  onMiss(m: missMessage) {
    this.boardManager.miss(m.coordinate.x, m.coordinate.y)
  }

  onShipDestroyed(m: shipDestroyedMessage) {
    this.boardManager.destoryShip(m.coordinate.x, m.coordinate.y, m.shipSize, m.vertical)
  }
}
