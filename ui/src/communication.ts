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
    if (m.type == 'HIT') {
      const hm = m as HitMessage
      this.onHit(hm.coordinate)
    } else if (m.type == 'MISS') {
      const mm = m as missMessage
      this.onMiss(mm.coordinate)
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

  onHit(coordinate: Coordinate) {
    this.boardManager.hit(coordinate.x, coordinate.y)
  }

  onMiss(coordinate: Coordinate) {
    this.boardManager.miss(coordinate.x, coordinate.y)
  }
}
