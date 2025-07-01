<template>
	<div class="user-selection">
		<!-- Search Section -->
		<div class="search-section">
			<input
				type="text"
				v-model="searchQuery"
				placeholder="Search users..."
				class="search-input"
			/>
		</div>

		<!-- Users List Section -->
		<div class="users-list-container">
			<div class="users-list" v-if="availableUsers.length > 0">
				<div
					v-for="user in availableUsers"
					:key="user.userId || user.username"
					@click="toggleUserSelection(user)"
					:class="{
						'user-selection-item': true,
						'user-selection-item-selected': isUserSelected(user),
					}"
				>
					<img
						:src="getUserPhotoUrl(user)"
						:alt="`${user.username}'s photo`"
						class="user-photo"
					/>
					<span class="user-name">{{ user.username }}</span>
					<span v-if="isUserSelected(user)" class="selected-badge"
						>✓</span
					>
				</div>
			</div>
		</div>

		<!-- Selected Users Preview -->
		<div v-if="selectedUsersList.length > 0" class="selected-users">
			<div class="selected-list">
				<span
					v-for="user in selectedUsersList"
					:key="user.userId || user.username"
					class="selected-user"
				>
					{{ user.username }}
					<button
						@click.stop="removeSelectedUser(user)"
						class="remove-selected"
					>
						✕
					</button>
				</span>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import axios from "../services/axios.js";
import { useUser } from "../composables/useUser.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import userDefaultIcon from "/assets/icons/user-default.png";

const { getUserId, getUsername } = useUser();
const { getImageUrl } = useImageUrl();

const props = defineProps({
	// Array of users to exclude from selection (e.g., current participants)
	excludeUsers: {
		type: Array,
		default: () => [],
	},
	// Pre-selected users
	initialSelectedUsers: {
		type: Array,
		default: () => [],
	},
	allowMultiple: {
		type: Boolean,
		default: true,
	},
});

const emit = defineEmits(["update:selectedUsers"]);

// State
const searchQuery = ref("");
const allUsers = ref([]);
const selectedUsers = ref([...props.initialSelectedUsers]);

// Current username to exclude from the list
const currentUsername = computed(() => getUsername());

const availableUsers = computed(() => {
	return allUsers.value.filter((user) => {
		const matchesSearch =
			searchQuery.value.trim() === "" ||
			user.username
				.toLowerCase()
				.includes(searchQuery.value.toLowerCase());

		const isNotCurrentUser = user.username !== currentUsername.value; // ← add .value

		const isNotExcluded = !props.excludeUsers.some(
			(excludedUser) => excludedUser.username === user.username
		);

		return matchesSearch && isNotCurrentUser && isNotExcluded;
	});
});

const logUserData = () => {
	console.log("All users:", allUsers.value);
	console.log("Current username:", currentUsername.value);
	console.log("Excluded users:", props.excludeUsers);
	console.log("Available users count:", availableUsers.value.length);
};

// Make sure fetchUsers logs what it finds
const fetchUsers = async () => {
	try {
		const response = await axios.get("/users", {
			headers: {
				Authorization: getUserId(),
			},
		});

		if (Array.isArray(response.data.users)) {
			allUsers.value = response.data.users;
			logUserData();
		} else {
			console.error("Expected users array but got:", response.data);
			allUsers.value = [];
		}
	} catch (error) {
		console.error("Error fetching users:", error);
		allUsers.value = [];
	}
};

const selectedUsersList = computed(() => {
	return selectedUsers.value;
});

const getUserPhotoUrl = (user) => {
	if (user.photo?.path) {
		return getImageUrl(user.photo.path);
	}
	return userDefaultIcon;
};

const isUserSelected = (user) => {
	return selectedUsers.value.some(
		(selectedUser) => selectedUser.username === user.username
	);
};
const toggleUserSelection = (user) => {
  if (props.allowMultiple) {
    // Original multi-select behavior
    const index = selectedUsers.value.findIndex(
      (selectedUser) => selectedUser.username === user.username
    );

    if (index !== -1) {
      // User is already selected, remove them
      selectedUsers.value.splice(index, 1);
    } else {
      // User is not selected, add them
      selectedUsers.value.push(user);
    }
  } else {
    // Single-select behavior
    if (selectedUsers.value.length === 1 && 
        selectedUsers.value[0].username === user.username) {
      // Clicking the selected user deselects them
      selectedUsers.value = [];
    } else {
      // Replace the current selection with the new user
      selectedUsers.value = [user];
    }
  }

  emit("update:selectedUsers", selectedUsers.value);
};

const removeSelectedUser = (user) => {
	selectedUsers.value = selectedUsers.value.filter(
		(selectedUser) => selectedUser.username !== user.username
	);
	emit("update:selectedUsers", selectedUsers.value);
};

// Lifecycle
onMounted(() => {
	fetchUsers();
});
</script>

<style scoped>
.user-selection-item {
	display: flex;
	align-items: center;
	padding: 10px 8px;
	cursor: pointer;
	transition: background-color 0.2s;
	border-bottom: 1px solid #f0f0f0;
}

.user-selection-item:last-child {
	border-bottom: none;
}

.user-selection-item-selected {
	background-color: #e3f2fd;
	border-color: #007aff;
}

.user-selection-item-selected .user-photo {
	border-color: #007aff;
}

.search-section {
	margin-bottom: 15px;
}

.search-input {
	width: 100%;
	padding: 10px;
	border: 1px solid #ddd;
	border-radius: 8px;
	font-size: 14px;
	box-sizing: border-box;
}

.search-input:focus {
	outline: none;
	border-color: #007aff;
	box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.1);
}

.users-list-container {
	margin-bottom: 15px;
}

.users-list {
	max-height: 200px;
	overflow-y: auto;
	border: 1px solid #eee;
	border-radius: 8px;
}

.empty-users-list {
	padding: 15px;
	text-align: center;
	color: #666;
	font-style: italic;
}

.user-item {
	display: flex;
	align-items: center;
	padding: 10px 8px;
	cursor: pointer;
	transition: background-color 0.2s;
	border-bottom: 1px solid #f0f0f0;
}

.user-item:last-child {
	border-bottom: none;
}

.user-item:hover {
	background-color: #f8f9fa;
}

.user-item.selected {
	background-color: #e3f2fd;
	border-color: #007aff;
}

.user-photo {
	width: 36px;
	height: 36px;
	border-radius: 50%;
	margin-right: 12px;
	object-fit: cover;
	border: 2px solid #ddd;
}

.user-item.selected .user-photo {
	border-color: #007aff;
}

.user-name {
	flex: 1;
	font-size: 14px;
	color: #333;
	font-weight: 500;
}

.selected-badge {
	background-color: #007aff;
	color: white;
	padding: 4px;
	border-radius: 50%;
	font-size: 12px;
	width: 16px;
	height: 16px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.selected-users {
	margin-top: 10px;
}

.selected-list {
	display: flex;
	flex-wrap: wrap;
	gap: 8px;
}

.selected-user {
	display: inline-flex;
	align-items: center;
	background-color: #e3f2fd;
	padding: 5px 10px;
	border-radius: 16px;
	font-size: 13px;
	color: #1976d2;
	font-weight: 500;
}

.remove-selected {
	background: none;
	border: none;
	color: #1976d2;
	cursor: pointer;
	margin-left: 4px;
	font-size: 12px;
}

.remove-selected:hover {
	color: #d32f2f;
}
</style>
