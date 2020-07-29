import { VuexModule, Module, Mutation, Action } from 'vuex-module-decorators' // eslint-disable-line

@Module({ name: 'User' })
class User extends VuexModule {
  public loggedIn = false

  @Mutation
  public LOG_IN(): void {
    this.loggedIn = true
  }

  @Action
  public logIn(): void {
    this.context.commit('LOG_IN')
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
