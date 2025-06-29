<template>
  <div class="modal">
    <div class="modal-content">
      <h2>New Conversation</h2>
      <input
        type="text"
        v-model="searchQuery"
        placeholder="Search users..."
        class="search-bar"
      />
      <ul class="user-list">
        <li
          v-for="user in filteredUsers"
          :key="user.username"
          @click="toggleUserSelection(user.username)"
          :class="{ selected: selectedUsers.includes(user.username) }"
        >
          <UserPreview :user="user" />
        </li>
      </ul>
      <div>
        <input
          type="checkbox"
          id="isGroupCheckbox"
          v-model="isGroup"
        />
        <label for="isGroupCheckbox">Group Conversation</label>
      </div>
      <div v-if="isGroup">
        <input
          type="text"
          v-model="groupName"
          placeholder="Enter group name"
          class="search-bar"
        />
      </div>
      <button
        @click="isGroup ? promptGroupDetails() : createPrivateConversation()"
        class="continue-button"
      >
        Continue
      </button>
      <button @click="$emit('close')" class="close-button">Close</button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from '../services/axios.js';
import UserPreview from '../components/UserPreview.vue';
import { useUser } from '../composables/useUser.js'; // Import the useUser composable

const { getUserId, getUsername } = useUser(); // Get both userId and username
const groupName = ref(''); // Add this for group name
const emits = defineEmits(['close']); // Define the close event
const searchQuery = ref('');
const users = ref([]);
const selectedUsers = ref([]);
const isGroup = ref(false);
const currentUsername = ref(getUsername()); // Get the current username

const fetchUsers = async () => {
  try {
    const response = await axios.get('/users');
    console.log('API Response:', response.data);
    users.value = Array.isArray(response.data.users) ? response.data.users : [];
  } catch (error) {
    console.error('Error fetching users:', error);
    users.value = []; // Fallback to an empty array
  }
};

const filteredUsers = computed(() =>
  users.value.filter((user) => 
    user.username.toLowerCase().includes(searchQuery.value.toLowerCase()) && 
    user.username !== currentUsername.value
  )
);

const toggleUserSelection = (username) => {
  if (selectedUsers.value.includes(username)) {
    selectedUsers.value = selectedUsers.value.filter((u) => u !== username);
  } else {
    selectedUsers.value.push(username);
  }
  // Automatically set isGroup to true if at least 3 users are selected
  isGroup.value = selectedUsers.value.length >= 3;
};

const createPrivateConversation = async () => {
  if (selectedUsers.value.length === 0) {
    alert("Please select at least one user");
    return;
  }

  const userId = getUserId();
  if (!userId) {
    console.error('User ID not found');
    return;
  }

  // Make sure to include current user in participants
  const allParticipants = [...selectedUsers.value];

  const requestBody = {
    participants: allParticipants,
    isGroup: false
  };

  try {
    await axios.post('/conversations', requestBody, {
      headers: {
        'Authorization': userId
      }
    });
    
    console.log('Private conversation created');
    // Close the modal and potentially refresh the conversations list
    emits('close');
  } catch (error) {
    console.error('Error creating private conversation:', error);
    alert('Failed to create conversation. Please try again.');
  }
};

const promptGroupDetails = async () => {
  if (selectedUsers.value.length === 0) {
    alert("Please select at least one user");
    return;
  }
  
  if (isGroup.value && !groupName.value) {
    alert("Please enter a name for the group");
    return;
  }

  const userId = getUserId();
  if (!userId) {
    console.error('User ID not found');
    return;
  }

  // Make sure to include current user in participants
  const allParticipants = [...selectedUsers.value];

  const requestBody = {
    name: groupName.value,
    participants: allParticipants,
    isGroup: true
  };

  try {
    await axios.post('/conversations', requestBody, {
      headers: {
        'Authorization': userId
      }
    });
    
    console.log('Group conversation created');
    // Close the modal and potentially refresh the conversations list
    emits('close');
  } catch (error) {
    console.error('Error creating group conversation:', error);
    alert('Failed to create group conversation. Please try again.');
  }
};

onMounted(() => {
  fetchUsers();
});
</script>

<style scoped>
.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-content {
  background: white;
  padding: 20px;
  border-radius: 8px;
  width: 800px;
  height: 800px;
}

.search-bar {
  width: 100%;
  padding: 8px;
  margin-bottom: 10px;
}

.user-list {
  list-style: none;
  padding: 0;
  max-height: 600px;
  overflow-y: auto;
}

.user-list li {
  padding: 8px;
  cursor: pointer;
}

.user-list li.selected {
  background-color: #007bff;
  color: white;
}

.continue-button {
  margin-top: 10px;
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}

.close-button {
  margin-top: 10px;
  padding: 8px 16px;
  background-color: #ccc;
  border: none;
  cursor: pointer;
}
</style>