<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Forward Message</h2>
        <button @click="closeModal" class="close-button">✕</button>
      </div>

      <div class="modal-body">
        <!-- Message Preview -->
        <div class="message-preview">
          <h3>Forwarding:</h3>
          <div class="preview-message">
            <div class="preview-sender">
              <img 
                :src="getImageUrl(message.sentBy.photo?.path) || '/assets/icons/user-default.png'" 
                alt="Sender Photo" 
                class="preview-photo"
              />
              <span class="preview-username">{{ message.sentBy.username }}</span>
            </div>
            <div class="preview-content">
              <img 
                v-if="message.photo" 
                :src="getImageUrl(message.photo.path)" 
                alt="Message Photo" 
                class="preview-image"
              />
              <p v-if="message.text" class="preview-text">{{ message.text }}</p>
              <span v-if="!message.text && !message.photo" class="no-content">Empty message</span>
            </div>
          </div>
        </div>

        <!-- Search Bar -->
        <div class="search-section">
          <input
            type="text"
            v-model="searchQuery"
            placeholder="Search conversations..."
            class="search-input"
          />
        </div>

        <!-- Conversations List -->
        <div class="conversations-section">
          <h3>Select Conversation ({{ filteredConversations.length }})</h3>
          
          <div class="conversations-list">
            <div
              v-for="conversation in filteredConversations"
              :key="conversation.conversationId"
              @click="selectConversation(conversation)"
              :class="{ 
                'conversation-item': true, 
                'selected': selectedConversation?.conversationId === conversation.conversationId,
                'current-conversation': conversation.conversationId === currentConversationId
              }"
            >
              <img
                :src="getConversationPhotoUrl(conversation)"
                :alt="getConversationName(conversation)"
                class="conversation-photo"
              />
              <div class="conversation-details">
                <span class="conversation-name">{{ getConversationName(conversation) }}</span>
                <span v-if="conversation.conversationId === currentConversationId" class="current-badge">
                  Current conversation
                </span>
                <span v-else-if="conversation.lastMessage" class="last-message">
                  {{ conversation.lastMessage.sentBy.username }}: {{ conversation.lastMessage.text || '📷 Photo' }}
                </span>
                <span v-else class="no-messages">No messages yet</span>
              </div>
              <span v-if="selectedConversation?.conversationId === conversation.conversationId" class="selected-badge">✓</span>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="modal-actions">
          <div class="left-info">
            <span v-if="selectedConversation" class="selection-info">
              Forward to: {{ getConversationName(selectedConversation) }}
            </span>
            <span v-else class="selection-info">
              Select a conversation to forward to
            </span>
          </div>
          
          <div class="right-actions">
            <button @click="closeModal" class="cancel-button">
              Cancel
            </button>
            <button 
              @click="forwardMessage" 
              class="forward-button"
              :disabled="!selectedConversation || selectedConversation.conversationId === currentConversationId"
            >
              Forward
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

const { getUserId, getUsername } = useUser();
const { getImageUrl } = useImageUrl();

const props = defineProps({
  message: {
    type: Object,
    required: true,
  },
  currentConversationId: {
    type: Number,
    required: true,
  },
});

const emits = defineEmits(['close', 'forwarded']);

// Reactive data
const searchQuery = ref('');
const conversations = ref([]);
const selectedConversation = ref(null);
const isForwarding = ref(false);

// Computed properties
const currentUsername = computed(() => getUsername());

const filteredConversations = computed(() => {
  if (!searchQuery.value) return conversations.value;
  
  return conversations.value.filter(conversation => {
    const name = getConversationName(conversation).toLowerCase();
    return name.includes(searchQuery.value.toLowerCase());
  });
});

// Methods
const fetchConversations = async () => {
	try {
		const userId = getUserId(); // Retrieve the userId using the composable

		if (!userId) {
			console.error("User ID not found");
			return;
		}

		const response = await axios.get("/conversations", {
			headers: {
				Authorization: userId,
			},
		});
		conversations.value = response.data.conversations;
		console.log("Conversations fetched:", conversations.value);
	} catch (error) {
		console.error("Error fetching conversations:", error);
	}
};

const getConversationName = (conversation) => {
  if (conversation.isGroup) {
    return conversation.name;
  } else {
    const otherParticipant = conversation.participants?.find(
      participant => participant.username !== currentUsername.value
    );
    return otherParticipant ? otherParticipant.username : conversation.name;
  }
};

const getConversationPhotoUrl = (conversation) => {
  if (conversation.isGroup) {
    if (conversation.photo?.path) {
      return getImageUrl(conversation.photo.path);
    }
    return '/assets/icons/group-default.png';
  } else {
    const otherParticipant = conversation.participants?.find(
      participant => participant.username !== currentUsername.value
    );
    if (otherParticipant?.photo?.path) {
      return getImageUrl(otherParticipant.photo.path);
    }
    return '/assets/icons/user-default.png';
  }
};

const selectConversation = (conversation) => {
  // Don't allow selecting the current conversation
  if (conversation.conversationId === props.currentConversationId) {
    return;
  }
  
  if (selectedConversation.value?.conversationId === conversation.conversationId) {
    selectedConversation.value = null; // Deselect if clicking the same conversation
  } else {
    selectedConversation.value = conversation;
  }
};

const forwardMessage = async () => {
  if (!selectedConversation.value || selectedConversation.value.conversationId === props.currentConversationId) {
    return;
  }

  isForwarding.value = true;

  try {
    await axios.post(
      `/conversations/${selectedConversation.value.conversationId}/forwarded_messages`,
      {
        forwardTo: selectedConversation.value.conversationId,
        messageId: props.message.messageId
      },
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: getUserId(),
        },
      }
    );

    console.log('Message forwarded successfully');
    emits('forwarded', {
      message: props.message,
      forwardedTo: selectedConversation.value
    });
    closeModal();
  } catch (error) {
    console.error('Error forwarding message:', error);
    alert('Failed to forward message. Please try again.');
  } finally {
    isForwarding.value = false;
  }
};

const closeModal = () => {
  emits('close');
};

// Lifecycle
onMounted(() => {
  fetchConversations();
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

.message-preview {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f8f9fa;
  border-radius: 8px;
  border: 1px solid #e1e1e1;
}

.message-preview h3 {
  margin: 0 0 12px 0;
  font-size: 16px;
  color: #333;
  font-weight: 600;
}

.preview-message {
  background: white;
  border-radius: 8px;
  padding: 12px;
}

.preview-sender {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.preview-photo {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  margin-right: 8px;
  object-fit: cover;
}

.preview-username {
  font-size: 12px;
  font-weight: bold;
  color: #666;
}

.preview-content {
  display: flex;
  flex-direction: column;
}

.preview-image {
  max-width: 100px;
  max-height: 100px;
  border-radius: 6px;
  margin-bottom: 8px;
  object-fit: cover;
}

.preview-text {
  margin: 0;
  font-size: 14px;
  color: #333;
  line-height: 1.3;
}

.no-content {
  font-size: 14px;
  color: #666;
  font-style: italic;
}

.search-section {
  margin-bottom: 20px;
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

.conversations-section h3 {
  margin: 0 0 15px 0;
  font-size: 18px;
  color: #333;
}

.conversations-list {
  max-height: 300px;
  overflow-y: auto;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 8px;
}

.conversation-item {
  display: flex;
  align-items: center;
  padding: 12px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid #f0f0f0;
}

.conversation-item:last-child {
  border-bottom: none;
}

.conversation-item:hover {
  background-color: #f8f9fa;
}

.conversation-item.selected {
  background-color: #e3f2fd;
  border-color: #007aff;
}

.conversation-item.current-conversation {
  opacity: 0.5;
  cursor: not-allowed;
}

.conversation-item.current-conversation:hover {
  background-color: transparent;
}

.conversation-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
  border: 2px solid #ddd;
}

.conversation-item.selected .conversation-photo {
  border-color: #007aff;
}

.conversation-details {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.conversation-name {
  font-size: 16px;
  color: #333;
  font-weight: 500;
  margin-bottom: 2px;
}

.current-badge {
  font-size: 12px;
  color: #ff9500;
  font-weight: 500;
}

.last-message {
  font-size: 12px;
  color: #666;
}

.no-messages {
  font-size: 12px;
  color: #999;
  font-style: italic;
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

.modal-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #eee;
}

.left-info {
  flex: 1;
}

.selection-info {
  font-size: 14px;
  color: #666;
  font-weight: 500;
}

.right-actions {
  display: flex;
  gap: 12px;
}

.cancel-button, .forward-button {
  padding: 12px 24px;
  border-radius: 8px;
  font-size: 16px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  min-width: 100px;
}

.cancel-button {
  background-color: #f5f5f5;
  color: #333;
}

.cancel-button:hover {
  background-color: #e0e0e0;
}

.forward-button {
  background-color: #007aff;
  color: white;
}

.forward-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.forward-button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: #ccc;
}

/* Responsive design */
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    margin: 20px;
    max-height: 90vh;
  }
  
  .modal-actions {
    flex-direction: column;
    align-items: stretch;
  }
  
  .left-info {
    text-align: center;
    margin-bottom: 10px;
  }
  
  .right-actions {
    width: 100%;
  }
  
  .cancel-button, .forward-button {
    width: 100%;
  }
  
  .conversations-list {
    max-height: 250px;
  }
}
</style>