<template>
  <div 
    :class="['message-container', { 'own-message': isOwnMessage }]"
    @contextmenu.prevent="showContextMenu"
    @click="hideContextMenu"
  >
    <!-- Existing message content -->
    <div v-if="isGroupConversation && !isOwnMessage" class="message-header">
      <img 
        :src="getImageUrl(message.sentBy.photo?.path) || '/assets/icons/user-default.png'" 
        alt="User Photo" 
        class="profile-picture"
      />
      <span class="username">{{ message.sentBy.username }}</span>
    </div>

    <div :class="['message-bubble', { 'own-bubble': isOwnMessage }]">
      <!-- Reply to message -->
      <div v-if="message.replyTo" class="reply-container">
        <div class="reply-header">
          <strong>{{ replyToMessage?.sentBy.username || 'Unknown User' }}</strong>
        </div>
        <div class="reply-content">
          <span v-if="replyToMessage?.text">{{ replyToMessage.text }}</span>
          <span v-else-if="replyToMessage?.photo" class="attachment-indicator">📷 Photo</span>
          <span v-else>Message not available</span>
        </div>
      </div>

      <!-- Message content -->
      <div class="message-content">
        <!-- Photo -->
        <img 
          v-if="message.photo" 
          :src="getImageUrl(message.photo.path)" 
          alt="Message Photo" 
          class="message-photo"
        />
        
        <!-- Text -->
        <p v-if="message.text" class="message-text">{{ message.text }}</p>
      </div>

      <!-- Message metadata -->
      <div class="message-metadata">
        <span class="timestamp">{{ formatTimestamp(message.timestamp) }}</span>
        <span v-if="isOwnMessage" :class="['status', message.status]">
          {{ getStatusIcon(message.status) }}
        </span>
      </div>
    </div>

    <!-- Context Menu -->
    <MessageContextMenu
      :isVisible="showMenu"
      :position="menuPosition"
      :message="message"
      @close="hideContextMenu"
      @reply="handleReply"
      @copy="handleCopy"
      @forward="handleForward"
      @delete="handleDelete"
    />
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useUser } from '../composables/useUser.js';
import { useImageUrl } from '../composables/useImageUrl.js';
import MessageContextMenu from './MessageContextMenu.vue';

const { getImageUrl } = useImageUrl();
const { getUsername } = useUser();
const currentUsername = getUsername();

const props = defineProps({
  message: {
    type: Object,
    required: true,
  },
  isGroupConversation: {
    type: Boolean,
    default: false,
  },
  replyToMessage: {
    type: Object,
    default: null,
  },
});

const emits = defineEmits(['reply']);

// Context menu state
const showMenu = ref(false);
const menuPosition = ref({ x: 0, y: 0 });

// Check if this message is sent by the current user
const isOwnMessage = computed(() => {
  return props.message.sentBy.username === currentUsername;
});

// Context menu handlers
const showContextMenu = (event) => {
  event.preventDefault();
  menuPosition.value = {
    x: event.clientX,
    y: event.clientY
  };
  showMenu.value = true;
};

const hideContextMenu = () => {
  showMenu.value = false;
};

const handleReply = (message) => {
  console.log('Reply to message:', message);
  emits('reply', message);
};

const handleCopy = (message) => {
  console.log('Copy message:', message);
  // Copy logic would go here
};

const handleForward = (message) => {
  console.log('Forward message:', message);
  // Forward logic would go here
};

const handleDelete = (message) => {
  console.log('Delete message:', message);
  // Delete logic would go here
};

// Format timestamp
const formatTimestamp = (timestamp) => {
  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
};

// Get status icon
const getStatusIcon = (status) => {
  switch (status) {
    case 'sent':
      return '✓';
    case 'delivered':
      return '✓✓';
    case 'read':
      return '✓✓';
    default:
      return '';
  }
};
</script>

<style scoped>
.message-container {
  display: flex;
  flex-direction: column;
  margin: 8px 0;
  align-items: flex-start;
}

.message-container.own-message {
  align-items: flex-end;
}

.message-header {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
  margin-left: 8px;
}

.profile-picture {
  width: 24px;
  height: 24px;
  border-radius: 50%;
  margin-right: 8px;
}

.username {
  font-size: 12px;
  font-weight: bold;
  color: #666;
}

.message-bubble {
  max-width: 70%;
  background-color: #e5e5ea;
  border-radius: 18px;
  padding: 12px 16px;
  position: relative;
}

.message-bubble.own-bubble {
  background-color: #007aff;
  color: white;
}

.reply-container {
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 8px;
  padding: 8px;
  margin-bottom: 8px;
  border-left: 3px solid #007aff;
}

.own-bubble .reply-container {
  background-color: rgba(255, 255, 255, 0.2);
  border-left-color: white;
}

.reply-header {
  font-size: 12px;
  margin-bottom: 2px;
}

.reply-content {
  font-size: 14px;
  opacity: 0.8;
}

.message-content {
  display: flex;
  flex-direction: column;
}

.message-photo {
  max-width: 200px;
  max-height: 200px;
  border-radius: 12px;
  margin-bottom: 8px;
  object-fit: cover;
}

.message-text {
  margin: 0;
  font-size: 16px;
  line-height: 1.3;
  word-wrap: break-word;
}

.message-metadata {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 4px;
  font-size: 11px;
  opacity: 0.7;
}

.timestamp {
  margin-right: 8px;
}

.status {
  font-size: 12px;
}

.status.delivered {
  color: #34c759;
}

.status.read {
  color: #34c759;
}

.attachment-indicator {
  font-style: italic;
}

/* Responsive design */
@media (max-width: 768px) {
  .message-bubble {
    max-width: 85%;
  }
  
  .message-photo {
    max-width: 150px;
    max-height: 150px;
  }
}
.message-container {
  position: relative; /* Add this for context menu positioning */
}

/* Add subtle hover effect to indicate interactive message */
.message-bubble {
  transition: transform 0.1s ease;
}

.message-container:hover .message-bubble {
  transform: scale(1.02);
}
</style>