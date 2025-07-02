<template>
  <div class="container vh-100 d-flex justify-content-center align-items-center">
    <form class="w-100" style="max-width: 320px;" @submit.prevent="submit">
      <h2 class="mb-4 text-center">Log into WASAText</h2>
      <div class="mb-3">
        <label for="username" class="form-label">Username</label>
        <input
          id="username"
          type="text"
          class="form-control"
          placeholder="Enter your username"
          v-model="username"
          @input="validateUsername"
          required
        />
        <div class="text-danger" v-if="usernameError"> {{ usernameError }}</div>
        <div class="text-danger" v-if="loginError"> {{ loginError }}</div>
      </div>
      <button type="submit" class="btn btn-primary w-100" :disabled="!isUsernameValid">Sign In</button>
    </form> 
  </div> 
</template>

<script setup>
import { ref } from 'vue';
import axios from '../services/axios.js';
import { useValidation } from '../composables/useValidation.js';

const { useUsernameValidation } = useValidation();
const emit = defineEmits(['loginSuccess']);

// Use the validation composable
const { username, usernameError, validateUsername, isUsernameValid } = useUsernameValidation();

// Additional error state for login failures
const loginError = ref('');

const submit = async () => {
  validateUsername();
  if (isUsernameValid.value) {
    try {
      const response = await axios.post('/session', { username: username.value });
      console.log('Login successful:', response.data);
      localStorage.setItem('loggedInUser', JSON.stringify(response.data));
      emit('loginSuccess', response.data);
    } catch (err) {
      console.error('Login failed:', err.response?.data || err.message);
      loginError.value = 'Login failed: ' + (err.response?.data || err.message);
    }
  }
};
</script>