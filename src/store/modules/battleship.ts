import { VuexModule, Module, Mutation } from 'vuex-module-decorators'

@Module({ name: 'Battleship' })
class Battleship extends VuexModule {
  connected = ''

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
}

export default Battleship
