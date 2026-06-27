<template>
  <div id="app">
    <LoginPage v-if="!user" @authenticated="onAuthenticated" />
    <HomePage v-else :user="user" @logout="onLogout" />
  </div>
</template>

<script>
import LoginPage from './pages/LoginPage.vue'
import HomePage from './pages/HomePage.vue'
import { getMe } from './api/documents.js'

export default {
  name: 'App',
  components: { LoginPage, HomePage },

  data() {
    return {
      user: null,
    }
  },

  async created() {
    const token = localStorage.getItem('kb_token')
    if (token) {
      const user = await getMe()
      if (user) {
        this.user = user
      } else {
        localStorage.removeItem('kb_token')
      }
    }
  },

  methods: {
    onAuthenticated(user) {
      this.user = user
    },

    onLogout() {
      this.user = null
    },
  },
}
</script>

<style>
*, *::before, *::after {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: 'Inter', 'Segoe UI', system-ui, sans-serif;
  background: #0f1117;
  color: #e8eaf0;
  min-height: 100vh;
}

#app {
  min-height: 100vh;
}
</style>