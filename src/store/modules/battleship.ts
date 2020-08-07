import { VuexModule, Module, Mutation } from 'vuex-module-decorators'

class Statistics {
  private hits: number
  private misses: number
  private destroyed: number
  constructor(public username: string) {
    this.hits = 0
    this.misses = 0
    this.destroyed = 0
  }

  hit() {
    this.hits++
  }

  miss() {
    this.misses++
  }

  shipDestroyed() {
    this.hits++
    this.destroyed++
  }

  reset() {
    this.hits = 0
    this.misses = 0
    this.destroyed = 0
  }
}

@Module({ name: 'Battleship' })
class Battleship extends VuexModule {
  connected = ''
  statistics: Statistics[] = []

  @Mutation
  public CONNECTED(): void {
    this.connected = 'CONNECTED'
  }

  @Mutation
  public DISCONNECTED(): void {
    this.connected = 'DISCONNECTED'
  }

  @Mutation
  public RECONNECTING(): void {
    this.connected = 'RECONNECTING'
  }

  @Mutation
  public FAILED(): void {
    this.connected = 'FAILED'
  }

  @Mutation
  public HIT(username: string): void {
    let stat = this.statistics.find(it => it.username === username)
    if (!stat) {
      stat = new Statistics(username)
      this.statistics.push(stat)
    }
    stat.hit()
  }

  @Mutation
  public MISS(username: string): void {
    let stat = this.statistics.find(it => it.username === username)
    if (!stat) {
      stat = new Statistics(username)
      this.statistics.push(stat)
    }
    stat.miss()
  }

  @Mutation
  public SHIPDESTROYED(username: string): void {
    let stat = this.statistics.find(it => it.username === username)
    if (!stat) {
      stat = new Statistics(username)
      this.statistics.push(stat)
    }
    stat.shipDestroyed()
  }

  @Mutation
  public RESETSTATS(usernames: string[]) {
    this.statistics = []
    usernames.forEach(it => {
      this.statistics.push(new Statistics(it))
    })
  }
}

export default Battleship
