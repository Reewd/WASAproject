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
        <div class="text-danger" v-if="errorMessage"> {{ errorMessage }}</div>
      </div>
      <button type="submit" class="btn btn-primary w-100" :disabled="error">Sign In</button>
    </form> 
  </div> 
</template>

<script setup>
import { ref } from 'vue';
import axios from '../services/axios.js';

const username = ref('');
const errorMessage = ref('');
const error = ref(false);
const emit = defineEmits(['loginSuccess']);

const validateUsername = () => {
  if (username.value.trim() === '') {
    errorMessage.value = 'Username is required';
    error.value = true;
  } else if (username.value.length < 3 || username.value.length > 16) {
    errorMessage.value = 'Username must be between 3 and 16 characters';
    error.value = true;
  } else {
    errorMessage.value = '';
    error.value = false;
  }
};

const submit = async () => {
  validateUsername();
  if (!error.value) {
    try {
      const response = await axios.post('/session', { username: username.value });
      console.log('Login successful:', response.data);
      localStorage.setItem('loggedInUser', JSON.stringify(response.data));
      emit('loginSuccess', response.data);
    } catch (err) {
      console.error('Login failed:', err.response?.data || err.message);
      errorMessage.value = 'Login failed: ' + (err.response?.data || err.message);
    }
  }
};
</script>