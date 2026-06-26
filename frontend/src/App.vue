<template>
  <div id="app">
    <HomePage />
  </div>
</template>

<script>
import HomePage from './pages/HomePage.vue'

export default {
  name: 'App',
  components: { HomePage },
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