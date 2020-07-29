import { VuexModule, Module, Mutation, Action } from 'vuex-module-decorators'

@Module({ name: 'User' })
class User extends VuexModule {
  loggedIn = false

  userName = ''

  @Mutation
  public LOG_IN(userName: string): void {
    this.loggedIn = true
    this.userName = userName
  }

  @Action
  public logIn(userName: string): void {
    this.context.commit('LOG_IN', userName)
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
