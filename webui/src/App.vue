<template>
  <component
    :is="isLoggedIn ? HomeView : LoginView"
    @login-success="login"
  />
</template>

<script setup>
import { useAuth } from '@/composables/useAuth'
import LoginView from './views/LoginView.vue'
import HomeView from './views/HomeView.vue'

const { isLoggedIn, login} = useAuth()
</script>