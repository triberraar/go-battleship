import { VuexModule, Module, Mutation, Action } from 'vuex-module-decorators'

@Module({ name: 'User' })
class User extends VuexModule {
  loggedIn = false

  username = ''

  @Mutation
  public LOG_IN(username: string): void {
    this.loggedIn = true
    this.username = username
  }

  @Action
  public logIn(username: string): void {
    this.context.commit('LOG_IN', username)
  }

  @Mutation
  public LOG_OUT(): void {
    this.loggedIn = false
  }

  @Action
  public logOut(): void {
    this.context.commit('LOG_OUT')
  }
}

export default User
