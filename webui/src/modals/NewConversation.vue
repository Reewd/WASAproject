<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>New Conversation</h2>
        <button @click="closeModal" class="close-button">✕</button>
      </div>

      <div class="modal-body">
        <!-- Search Section -->
        <div class="search-section">
          <label for="userSearch">Search Users:</label>
          <input
            id="userSearch"
            type="text"
            v-model="searchQuery"
            placeholder="Search users..."
            class="search-input"
          />
        </div>

        <!-- Users List Section -->
        <div class="users-section">
          <h3>Select Users ({{ selectedUsers.length }} selected)</h3>
          
          <div class="users-list">
            <div
              v-for="user in filteredUsers"
              :key="user.username"
              @click="toggleUserSelection(user.username)"
              :class="{ selected: selectedUsers.includes(user.username) }"
              class="user-item"
            >
              <img
                :src="getUserPhotoUrl(user)"
                :alt="`${user.username}'s photo`"
                class="user-photo"
              />
              <span class="user-name">{{ user.username }}</span>
              <span v-if="selectedUsers.includes(user.username)" class="selected-badge">✓</span>
            </div>
          </div>
        </div>

        <!-- Selected Users Preview -->
        <div v-if="selectedUsers.length > 0" class="selected-users">
          <h4>Selected Users:</h4>
          <div class="selected-list">
            <span 
              v-for="username in selectedUsers" 
              :key="username"
              class="selected-user"
            >
              {{ username }}
              <button @click="removeSelectedUser(username)" class="remove-selected">✕</button>
            </span>
          </div>
        </div>

        <!-- Conversation Type Section -->
        <div class="conversation-type-section">
          <div class="checkbox-container">
            <input
              type="checkbox"
              id="isGroupCheckbox"
              v-model="isGroup"
              class="group-checkbox"
            />
            <label for="isGroupCheckbox" class="checkbox-label">
              Create Group Conversation
            </label>
          </div>
          <p class="checkbox-hint">
            {{ isGroup ? 'Group conversations allow multiple participants' : 'Private conversation with selected user(s)' }}
          </p>
        </div>

        <!-- Group Name Section (only for groups) -->
        <div v-if="isGroup" class="group-name-section">
          <label for="groupName">Group Name:</label>
          <input
            id="groupName"
            type="text"
            v-model="groupName"
            placeholder="Enter group name"
            class="group-name-input"
          />
        </div>

        <!-- Action Buttons -->
        <div class="modal-actions">
          <div class="left-info">
            <span class="participants-info">
              {{ selectedUsers.length }} participant{{ selectedUsers.length !== 1 ? 's' : '' }} selected
            </span>
          </div>
          
          <div class="right-actions">
            <button @click="closeModal" class="cancel-button">
              Cancel
            </button>
            <button 
              @click="createConversation" 
              class="create-button"
              :disabled="!canCreateConversation"
            >
              {{ isGroup ? 'Create Group' : 'Create Conversation' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from '../services/axios.js';
import { useUser } from '../composables/useUser.js';
import { useImageUrl } from '../composables/useImageUrl.js';
import userDefaultIcon from "/assets/icons/user-default.png";

const { getUserId, getUsername } = useUser();
const { getImageUrl } = useImageUrl();

const emits = defineEmits(['close', 'conversationCreated']);
const props = defineProps({
  existingConversations: {
    type: Array,
    default: () => []
  }
});
// Reactive data
const searchQuery = ref('');
const users = ref([]);
const selectedUsers = ref([]);
const isGroup = ref(false);
const groupName = ref('');

// Computed properties
const currentUsername = computed(() => getUsername());

const filteredUsers = computed(() =>
  users.value.filter((user) => 
    user.username.toLowerCase().includes(searchQuery.value.toLowerCase()) && 
    user.username !== currentUsername.value
  )
);

const canCreateConversation = computed(() => {
  if (selectedUsers.value.length === 0) return false;
  if (isGroup.value && !groupName.value.trim()) return false;
  return true;
});

// Methods
const fetchUsers = async () => {
  try {
    const response = await axios.get('/users');
    console.log('API Response:', response.data);
    users.value = Array.isArray(response.data.users) ? response.data.users : [];
  } catch (error) {
    console.error('Error fetching users:', error);
    users.value = [];
  }
};

const getUserPhotoUrl = (user) => {
  if (user.photo?.path) {
    return getImageUrl(user.photo.path);
  }
  return userDefaultIcon;
};

const toggleUserSelection = (username) => {
  if (selectedUsers.value.includes(username)) {
    selectedUsers.value = selectedUsers.value.filter((u) => u !== username);
  } else {
    selectedUsers.value.push(username);
  }
  
  // Automatically set isGroup to true if more than 1 user is selected
  if (selectedUsers.value.length > 1 && !isGroup.value) {
    isGroup.value = true;
  }
};

const removeSelectedUser = (username) => {
  selectedUsers.value = selectedUsers.value.filter((u) => u !== username);
  
  // If only 1 or no users left, set to private conversation
  if (selectedUsers.value.length <= 1) {
    isGroup.value = false;
  }
};

const createConversation = async () => {
  if (!canCreateConversation.value) {
    return;
  }

  const userId = getUserId();
  if (!userId) {
    console.error('User ID not found');
    alert('Authentication error. Please try again.');
    return;
  }

  const requestBody = {
    participants: selectedUsers.value,
    isGroup: isGroup.value,
    ...(isGroup.value && { name: groupName.value.trim() })
  };

  try {
    await axios.post('/conversations', requestBody, {
      headers: {
        'Authorization': userId
      }
    });

    console.log(`${isGroup.value ? 'Group' : 'Private'} conversation created`);
    emits('conversationCreated'); // Emit event to notify parent component
    emits('close');
  } catch (error) {
    console.error('Error creating conversation:', error);
    alert('Failed to create conversation. Please try again.');
  }
};

const closeModal = () => {
  emits('close');
};

// Lifecycle
onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 600px;
  max-height: 85vh;
  overflow-y: auto;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #eee;
}

.modal-header h2 {
  margin: 0;
  font-size: 24px;
  color: #333;
}

.close-button {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #666;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  background-color: #f0f0f0;
}

.modal-body {
  padding: 20px;
}

.search-section {
  margin-bottom: 20px;
}

.search-section label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.search-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
}

.search-input:focus {
  outline: none;
  border-color: #007aff;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.users-section {
  margin-bottom: 20px;
}

.users-section h3 {
  margin: 0 0 15px 0;
  font-size: 18px;
  color: #333;
}

.users-list {
  max-height: 250px;
  overflow-y: auto;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 10px;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px 8px;
  border-radius: 6px;
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
  width: 40px;
  height: 40px;
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
  font-size: 16px;
  color: #333;
  font-weight: 500;
}

.selected-badge {
  background-color: #007aff;
  color: white;
  padding: 4px 8px;
  border-radius: 50%;
  font-size: 12px;
  font-weight: bold;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.selected-users {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.selected-users h4 {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: #333;
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
  padding: 6px 12px;
  border-radius: 20px;
  font-size: 14px;
  color: #1976d2;
  font-weight: 500;
}

.remove-selected {
  background: none;
  border: none;
  color: #1976d2;
  cursor: pointer;
  margin-left: 6px;
  font-size: 12px;
  font-weight: bold;
}

.remove-selected:hover {
  color: #d32f2f;
}

.conversation-type-section {
  margin-bottom: 20px;
  padding: 15px;
  border: 1px solid #eee;
  border-radius: 8px;
}

.checkbox-container {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.group-checkbox {
  margin-right: 10px;
  width: 18px;
  height: 18px;
}

.checkbox-label {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  cursor: pointer;
}

.checkbox-hint {
  font-size: 14px;
  color: #666;
  margin: 0;
  font-style: italic;
}

.group-name-section {
  margin-bottom: 25px;
}

.group-name-section label {
  display: block;
  margin-bottom: 8px;
  font-weight: bold;
  color: #333;
}

.group-name-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  font-size: 16px;
  box-sizing: border-box;
}

.group-name-input:focus {
  outline: none;
  border-color: #007aff;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.modal-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.left-info {
  flex: 1;
}

.participants-info {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.right-actions {
  display: flex;
  gap: 12px;
}

.cancel-button, .create-button {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  min-width: 120px;
}

.cancel-button {
  background-color: #f5f5f5;
  color: #333;
}

.cancel-button:hover {
  background-color: #e0e0e0;
}

.create-button {
  background-color: #007aff;
  color: white;
}

.create-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.create-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #ccc;
}
</style>