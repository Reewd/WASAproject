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
				:src="userDefaultIcon"
				alt="Default profile picture"
				class="profile-picture"
			/>
		</div>
		<button class="settings-button" @click="openProfileSettings">
			<img
				:src="accountSettingsIcon"
				alt="Settings"
				class="settings-icon"
			/>
		</button>

		<button class="logout-button" @click="logout">
			<img
				:src="logoutIcon"
				alt="Logout"
				class="logout-icon"
			/>
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
import { ref, onMounted } from "vue";
import { useImageUrl } from "@/composables/useImageUrl.js";
import ProfileSettings from "@/modals/ProfileSettings.vue";
import userDefaultIcon from "/assets/icons/user-default.png";
import accountSettingsIcon from "/assets/icons/account-settings.png";
import logoutIcon from "/assets/icons/logout.png";

const loggedInUser = ref(null);
const showProfileSettings = ref(false);
const { getImageUrl } = useImageUrl();

const loadUserData = () => {
	const userData = localStorage.getItem("loggedInUser");
	if (userData) {
		loggedInUser.value = JSON.parse(userData);
		console.log("Logged in user:", loggedInUser.value);
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

const logout = () => {
  // Clear user data from localStorage
  localStorage.removeItem('loggedInUser');
  localStorage.removeItem('userId');
  
  // Reset user state
  loggedInUser.value = null;
  window.location.reload();

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

.settings-button {
	width: 50px;
	height: 50px;
	background: none;
	border: none;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	transition: background-color 0.3s ease;
}

.settings-button:hover {
	background-color: #bfbfbf;
}

.settings-icon {
	width: 30px;
	height: 30px;
}

.logout-button {
	width: 50px;
	height: 50px;
	background: none;
	border: none;
	cursor: pointer;
	display: flex;
	align-items: center;
	justify-content: center;
	border-radius: 50%;
	transition: background-color 0.3s ease;
	margin-top: auto;
	margin-bottom: 20px;
}

.logout-button:hover {
	background-color: #ffcccc;
}

.logout-icon {
	width: 28px;
	height: 28px;
}
</style>
