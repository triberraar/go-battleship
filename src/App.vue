<template>
  <span>
    <navigation />
    <router-view />
  </span>
</template>

<script>
import Navigation from '@/components/Navigation.vue'

export default {
  name: 'App',
  components: {
    Navigation
  }
}
</script>
