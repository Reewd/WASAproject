<template>
  <div class="sidebar">
    <div class="profile-wrapper">
      <img
        v-if="loggedInUser?.photo?.path"
        :src="getImageUrl(loggedInUser.photo.path)"
        alt="Profile picture"
        class="profile-picture"
      />
      <img
        v-else
        src="/assets/icons/user-default.png"
        alt="Default profile picture"
        class="profile-picture"
      />
    </div>
    <button class="settings-button">
      <img src="/assets/icons/account-settings.png" alt="Settings" class="settings-icon" />
    </button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useImageUrl } from '@/composables/useImageUrl.js';

const loggedInUser = ref(null);
const { getImageUrl } = useImageUrl();

onMounted(() => {
  const userData = localStorage.getItem('loggedInUser');
  if (userData) {
    loggedInUser.value = JSON.parse(userData);
    console.log('Logged in user:', loggedInUser.value);
  }
});
</script>

<style>
.sidebar {
  width: 80px;
  height: 100vh;
  background: #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 20px;
}

.profile-wrapper {
  width: 50px;
  height: 50px;
  margin-bottom: 20px;
}

.profile-picture {
  width: 100%;
  height: 100%;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #ddd;
}

.profile-picture:hover {
  border-color: #aaa; /* Change border color on hover */
}

.settings-button {
  width: 50px;
  height: 50px;
  background: none;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%; /* Apply circular shape by default */
  transition: background-color 0.3s ease; /* Smooth transition */
}

.settings-button:hover {
  background-color: #bfbfbf; /* Dark background on hover */
}

.settings-icon {
  width: 30px;
  height: 30px;
}
</style>