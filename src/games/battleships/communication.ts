import BoardManager from './board'
import FeedbackText from './feedbackText'

interface Coordinate {
  x: number;
  y: number;
}

interface HitMessage {
  coordinate: Coordinate;
}

interface missMessage {
  coordinate: Coordinate;
}

interface shipDestroyedMessage {
  coordinate: Coordinate;
  shipSize: number;
  vertical: boolean;
}

interface boardMessage {
  shipSizes: number[];
}

interface gameStartedMessage {
  turn: boolean;
  duration: number;
}

interface turnMessage {
  turn: boolean;
  duration: number;
}

interface turnExtendedMessage {
  duration: number;
}

export default class CommunicationManager {
  private boardManager: BoardManager

  private feedbackText: FeedbackText

  private ws: WebSocket

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
    this.ws.onopen = this.onopen
    this.ws.onmessage = this.onmessage
    this.ws.onerror = this.onerror
    this.waitForConnection()
  }

  waitForConnection() {}

  close() {
    this.ws.close(1000)
  }

  send(message: string) {
    if (this.ws.readyState !== 1) {
      console.log('waiting for connection')
      setTimeout(() => this.send(message), 5)
    } else {
      this.ws.send(message)
    }
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
    }
  }

  onerror = (ev: Event) => {
    console.log(`error ${ev}`)
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
    this.send(JSON.stringify({ type: 'PLAY' }))
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

  onVictory() {
    this.boardManager.victory()
    this.feedbackText.setText('You win')
    this.close()
  }

  onLoss() {
    this.boardManager.loss()
    this.feedbackText.setText('The other dummy won, loser')
  }

  onBoard(m: boardMessage) {
    console.log(this.boardManager)
    this.boardManager.ships(m.shipSizes)
  }

  onGameStarted(m: gameStartedMessage) {
    if (m.turn) {
      this.feedbackText.setCountDownText('Your turn', m.duration)
    } else {
      this.feedbackText.setText('Waiting for the other dummy')
    }
  }

  onTurn(m: turnMessage) {
    if (m.turn) {
      this.feedbackText.setCountDownText('Your turn', m.duration)
    } else {
      this.feedbackText.setText('Waiting for the other dummy')
    }
  }

  onTurnExtended(m: turnExtendedMessage) {
    this.feedbackText.setCountDownText('Your turn', m.duration)
  }
}
