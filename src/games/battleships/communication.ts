import BoardManager from './board' // eslint-disable-line
import FeedbackText from './feedbackText'
import UserStore from '@/store/modules/user'
import Store from '@/store/index'

interface Coordinate {
  x: number
  y: number
}

interface HitMessage {
  username: string
  coordinate: Coordinate
}

interface MissMessage {
  username: string
  coordinate: Coordinate
}

interface ShipDestroyedMessage {
  username: string
  coordinate: Coordinate
  shipSize: number
  vertical: boolean
}

interface BoardMessage {
  username: string
  shipSizes: number[]
}

interface GameStartedMessage {
  username: string
  turn: boolean
  duration: number
  usernames: string[]
}

interface TurnMessage {
  username: string
  turn: boolean
  duration: number
}

interface TurnExtendedMessage {
  username: string
  duration: number
}

interface BoardStateMessage {
  username: string
  board: BoardMessage
  destroys: ShipDestroyedMessage[]
  hits: HitMessage[]
  misses: MissMessage[]
}

interface VictoryMessage {
  username: string
}

interface LossMessage {
  username: string
}

interface OpponentDestroyedShip {
  username: string
}

export default class CommunicationManager {
  private boardManager: BoardManager
  private feedbackText: FeedbackText

  private ws: WebSocket
  // private pingTimer: number
  private reconnectAttempts = 0

  constructor() {
    this.reconnect()
  }

  close() {
    this.feedbackText.clear()
    this.ws.close(1000)
    console.log('closed socket')
  }

  reconnect() {
    if (this.reconnectAttempts > 15) {
      Store.commit('FAILED')
      return
    }
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
    wsUri += `?username=${UserStore.state.username}`

    this.reconnectAttempts++
    this.ws = new WebSocket(wsUri)
    this.ws.onmessage = this.onmessage
    this.ws.onerror = this.onerror
    this.ws.onopen = this.onopen
    this.ws.onclose = this.onclose
  }

  send(message: string) {
    if (this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(message)
    } else if (this.ws.readyState === WebSocket.CONNECTING) {
      console.log('waiting for connection')
      setTimeout(() => this.send(message), 1000)
    }
  }

  onopen = () => {
    Store.commit('CONNECTED')
    this.reconnectAttempts = 0
    this.play()
  }

  onerror = (ev: Event) => {
    console.log(`error ${ev}`)
    // clearInterval(this.pingTimer)
    Store.commit('RECONNECTING')
  }

  onclose = (ev: CloseEvent) => {
    // clearInterval(this.pingTimer)
    console.log(`onclose ${ev}`)
    if (!ev.wasClean) {
      Store.commit('RECONNECTING')
      setTimeout(() => this.reconnect(), 1000)
    } else {
      Store.commit('DISCONNECTED')
    }
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
        this.onVictory(m)
        break
      }
      case 'LOSS': {
        this.onLoss(m)
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
      case 'BOARD_STATE': {
        this.onBoardState(m)
        break
      }
      case 'OPPONENT_DESTROYED_SHIP': {
        this.onOpponentDestroyedShip(m)
        break
      }
      case 'CANCELLED': {
        this.onCancelled()
        break
      }
      default: {
        console.error(`unknowns message ${m.type}`)
        break
      }
    }
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
    this.send(JSON.stringify({ type: 'PLAY', username: UserStore.state.username }))
  }

  onHit(m: HitMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.hit(m.coordinate.x, m.coordinate.y)
    }
    Store.commit('HIT', m.username)
  }

  onMiss(m: MissMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.miss(m.coordinate.x, m.coordinate.y)
    }
    Store.commit('MISS', m.username)
  }

  onShipDestroyed(m: ShipDestroyedMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.destoryShip(m.coordinate.x, m.coordinate.y, m.shipSize, m.vertical)
    }
    Store.commit('SHIPDESTROYED', m.username)
  }

  onOpponentDestroyedShip(m: OpponentDestroyedShip) {
    Store.commit('SHIPDESTROYED', m.username)
  }

  onVictory(m: VictoryMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.victory()
      this.feedbackText.setText('You win')
      this.close()
    }
  }

  onLoss(m: LossMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.loss()
      this.feedbackText.setText('The other dummy won, loser')
      this.close()
    }
  }

  onCancelled() {
    this.boardManager.loss()
    this.feedbackText.setText('The game got cancelled :(')
    this.close()
  }

  onBoard(m: BoardMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.ships(m.shipSizes)
    }
  }

  onGameStarted(m: GameStartedMessage) {
    if (m.username === UserStore.state.username) {
      if (m.turn) {
        this.feedbackText.setCountDownText('Your turn', m.duration)
      } else {
        this.feedbackText.setText('Waiting for the other dummy')
      }
    }
    Store.commit('RESETSTATS', m.usernames)
  }

  onTurn(m: TurnMessage) {
    if (m.username === UserStore.state.username) {
      if (m.turn) {
        this.feedbackText.setCountDownText('Your turn', m.duration)
      } else {
        this.feedbackText.setText('Waiting for the other dummy')
      }
    }
  }

  onTurnExtended(m: TurnExtendedMessage) {
    if (m.username === UserStore.state.username) {
      this.feedbackText.setCountDownText('Your turn', m.duration)
    }
  }

  onBoardState(m: BoardStateMessage) {
    if (m.username === UserStore.state.username) {
      this.boardManager.ships(m.board.shipSizes)
      m.destroys.forEach(d => this.onShipDestroyed(d))
      m.hits.forEach(h => this.onHit(h))
      m.misses.forEach(ms => this.onMiss(ms))
    }
  }
}
