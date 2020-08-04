import BoardManager from './board' // eslint-disable-line
import FeedbackText from './feedbackText'
import UserStore from '@/store/modules/user'

interface Coordinate {
  x: number
  y: number
}

interface HitMessage {
  coordinate: Coordinate
}

interface MissMessage {
  coordinate: Coordinate
}

interface ShipDestroyedMessage {
  coordinate: Coordinate
  shipSize: number
  vertical: boolean
}

interface BoardMessage {
  shipSizes: number[]
}

interface GameStartedMessage {
  turn: boolean
  duration: number
}

interface TurnMessage {
  turn: boolean
  duration: number
}

interface TurnExtendedMessage {
  duration: number
}

export default class CommunicationManager {
  private boardManager: BoardManager

  private feedbackText: FeedbackText

  private ws: WebSocket

  private pingTimer: number

  constructor() {
    const loc = window.location
    let wsUri = ''
    if (loc.protocol === 'https:') {
      wsUri = 'wss:'
    } else {
      wsUri = 'ws:'
    }
    wsUri += `//${loc.host}/battleship`
    if (loc.host.startsWith('localhost')) {
      wsUri = 'ws://localhost:10002/battleship'
    }

    this.ws = new WebSocket(wsUri)
    this.ws.onmessage = this.onmessage
    this.ws.onerror = this.onerror
    this.ws.onopen = this.onopen
    this.ws.onclose = this.onclose
  }

  close() {
    this.ws.close(1000)
  }

  send(message: string) {
    console.log('semd')
    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message)
    } else if (this.ws.readyState === WebSocket.CONNECTING) {
      console.log('waiting for connection')
      setTimeout(() => this.send(message), 1000)
    } else if (
      this.ws.readyState === WebSocket.CLOSING ||
      this.ws.readyState === WebSocket.CLOSED
    ) {
      // reconnect
    }
  }

  onopen = () => {
    this.pingTimer = setInterval(() => this.send(JSON.stringify({ type: 'PING' })), 5000)
  }

  onclose = () => {
    clearInterval(this.pingTimer)
  }

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
      case 'VICTORY': {
        this.onVictory()
        break
      }
      case 'LOSS': {
        this.onLoss()
        break
      }
      case 'BOARD': {
        this.onBoard(m)
        break
      }
      case 'GAME_STARTED': {
        this.onGameStarted(m)
        break
      }
      case 'TURN': {
        this.onTurn(m)
        break
      }
      case 'TURN_EXTENDED': {
        this.onTurnExtended(m)
        break
      }
      default: {
        console.error(`unknowns message ${m.type}`)
        break
      }
    }
  }

  onerror = (ev: Event) => {
    console.log(`error ${ev}`)
    clearInterval(this.pingTimer)
  }

  setBoardManager(boardManager: BoardManager) {
    this.boardManager = boardManager
  }

  setFeedbackText(fd: FeedbackText) {
    this.feedbackText = fd
  }

  fire(x: number, y: number) {
    this.send(JSON.stringify({ type: 'FIRE', coordinate: { x, y } }))
  }

  play() {
    console.log(`player ${UserStore.state.username}`)
    this.send(JSON.stringify({ type: 'PLAY', username: UserStore.state.username }))
  }

  onHit(m: HitMessage) {
    this.boardManager.hit(m.coordinate.x, m.coordinate.y)
  }

  onMiss(m: MissMessage) {
    this.boardManager.miss(m.coordinate.x, m.coordinate.y)
  }

  onShipDestroyed(m: ShipDestroyedMessage) {
    this.boardManager.destoryShip(m.coordinate.x, m.coordinate.y, m.shipSize, m.vertical)
  }

  onVictory() {
    this.boardManager.victory()
    this.feedbackText.setText('You win')
    this.close()
  }

  onLoss() {
    this.boardManager.loss()
    this.feedbackText.setText('The other dummy won, loser')
  }

  onBoard(m: BoardMessage) {
    console.log(this.boardManager)
    this.boardManager.ships(m.shipSizes)
  }

  onGameStarted(m: GameStartedMessage) {
    if (m.turn) {
      this.feedbackText.setCountDownText('Your turn', m.duration)
    } else {
      this.feedbackText.setText('Waiting for the other dummy')
    }
  }

  onTurn(m: TurnMessage) {
    if (m.turn) {
      this.feedbackText.setCountDownText('Your turn', m.duration)
    } else {
      this.feedbackText.setText('Waiting for the other dummy')
    }
  }

  onTurnExtended(m: TurnExtendedMessage) {
    this.feedbackText.setCountDownText('Your turn', m.duration)
  }
}
