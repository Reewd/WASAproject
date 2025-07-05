<template>
	<div class="sidebar">
		<div class="profile-wrapper">
			<img
				v-if="user?.photo?.path"
				:src="getImageUrl(user.photo.path)"
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
		/>
	</div>
</template>

<script setup>
import { ref } from "vue";
import { useImageUrl } from "@/composables/useImageUrl.js";
import ProfileSettings from "@/modals/ProfileSettings.vue";
import { useAuth } from "@/composables/useAuth.js";
import userDefaultIcon from "/assets/icons/user-default.png";
import accountSettingsIcon from "/assets/icons/account-settings.png";
import logoutIcon from "/assets/icons/logout.png";

const showProfileSettings = ref(false);
const { getImageUrl } = useImageUrl();
const { user, logout: authLogout } = useAuth();

const openProfileSettings = () => {
  showProfileSettings.value = true;
};

const closeProfileSettings = () => {
  showProfileSettings.value = false;
};

const logout = () => {
  authLogout();
  window.location.reload();
};
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
	border: none;
	cursor: pointer;
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
	border: none;
	cursor: pointer;
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
	width: 30px;
	height: 30px;
}
</style>
