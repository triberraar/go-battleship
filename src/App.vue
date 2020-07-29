<template>
  <div>
    <b-navbar>
      <template slot="brand">
        <b-navbar-item>
          <img
            src="https://storage.googleapis.com/hexigames-logos/logo_dark_low_res.png"
            alt="Hexigames"
          />
        </b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item class="is-dark" tag="router-link" :to="{ name: 'Home' }">
          Home
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ name: 'Accounting' }">
          Accounting
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ name: 'Game' }">
          Games
        </b-navbar-item>
      </template>

      <template slot="end">
        <b-navbar-item tag="div">
          <div class="buttons">
            <a class="button is-light">
              Log in
            </a>
          </div>
        </b-navbar-item>
      </template>
    </b-navbar>
    <router-view />
  </div>
</template>

<style></style>
