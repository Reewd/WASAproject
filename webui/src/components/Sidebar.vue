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
    <button class="settings-button" @click="openProfileSettings">
      <img src="/assets/icons/account-settings.png" alt="Settings" class="settings-icon" />
    </button>

    <!-- Profile Settings Modal -->
    <ProfileSettings 
      v-if="showProfileSettings" 
      @close="closeProfileSettings"
      @updated="handleProfileUpdated"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useImageUrl } from '@/composables/useImageUrl.js';
import ProfileSettings from '@/modals/ProfileSettings.vue';

const loggedInUser = ref(null);
const showProfileSettings = ref(false);
const { getImageUrl } = useImageUrl();

const loadUserData = () => {
  const userData = localStorage.getItem('loggedInUser');
  if (userData) {
    loggedInUser.value = JSON.parse(userData);
    console.log('Logged in user:', loggedInUser.value);
  }
};

const openProfileSettings = () => {
  showProfileSettings.value = true;
};

const closeProfileSettings = () => {
  showProfileSettings.value = false;
};

const handleProfileUpdated = () => {
  // Reload user data from localStorage after profile update
  loadUserData();
};

onMounted(() => {
  loadUserData();
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