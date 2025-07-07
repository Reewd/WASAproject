<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>New Conversation</h2>
        <button @click="closeModal" class="close-button">✕</button>
      </div>

      <div class="modal-body">
        <!-- Users List Section -->
        <div class="users-section">
          <h3>Select Users ({{ selectedUsers.length }} selected)</h3>
          
          <UserSelection
            :allowMultiple="true"
            @update:selectedUsers="handleSelectedUsersUpdate"
          />
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
        </div>

        <!-- Group Name Section (only for groups) -->
        <div v-if="isGroup" class="group-name-section">
          <label for="groupName">Group Name:</label>
          <input
            id="groupName"
            type="text"
            v-model="groupName"
            @input="validateGroupName"
            placeholder="Enter group name"
            class="group-name-input"
          />
          <div class="text-danger" v-if="groupNameError">{{ groupNameError }}</div>
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
import { useAuth } from "../composables/useAuth.js";
import { useValidation } from '../composables/useValidation.js';
import UserSelection from '../components/UserSelection.vue';

const { user : currentUser} = useAuth();
const { useGroupNameValidation } = useValidation();

const emits = defineEmits(['close', 'conversationCreated']);
const props = defineProps({
  existingConversations: {
    type: Array,
    default: () => []
  }
});
// Reactive data
const selectedUsers = ref([]);
const isGroup = ref(false);

// Use the validation composable
const { groupName, groupNameError, validateGroupName, isGroupNameValid } = useGroupNameValidation();

// Computed properties
const canCreateConversation = computed(() => {
  if (selectedUsers.value.length === 0) return false;
  if (isGroup.value && (!groupName.value.trim() || !isGroupNameValid.value)) return false;
  return true;
});

// Methods
const handleSelectedUsersUpdate = (selectedUsersList) => {
  selectedUsers.value = selectedUsersList;
  
  // Automatically set isGroup to true if more than 1 user is selected
  if (selectedUsers.value.length > 1 && !isGroup.value) {
    isGroup.value = true;
  }
  // If only 1 or no users left, set to private conversation
  else if (selectedUsers.value.length <= 1) {
    isGroup.value = false;
  }
};

const createConversation = async () => {
  if (!canCreateConversation.value) {
    return;
  }

  // Validate group name if it's a group
  if (isGroup.value) {
    validateGroupName();
    if (groupNameError.value) {
      return;
    }
  }

  const userId = currentUser.value.userId;
  if (!userId) {
    console.error('User ID not found');
    alert('Authentication error. Please try again.');
    return;
  }

  const requestBody = {
    participants: selectedUsers.value.map(user => user.username), // Extract usernames for backend
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

.users-section {
  margin-bottom: 20px;
}

.users-section h3 {
  margin: 0 0 15px 0;
  font-size: 18px;
  color: #333;
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
}

.group-checkbox {
  margin-right: 10px;
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.checkbox-label {
  font-size: 16px;
  font-weight: 500;
  color: #333;
  cursor: pointer;
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

.text-danger {
  color: #dc3545;
  font-size: 14px;
  margin-top: 5px;
}
</style>