<template>
  <div class="conversation-preview">
    <img 
      :src="conversationPhotoUrl" 
      alt="Conversation Photo" 
      class="conversation-photo"
    />
    <div class="conversation-details">
      <h3 class="conversation-name">
        {{ displayName }}
      </h3>
      <p class="last-message" v-if="conversation.lastMessage">
        <strong>{{ conversation.lastMessage.sentBy.username }}:</strong> 
        {{ conversation.lastMessage.text || '📷 Photo' }}
      </p>
      <p class="last-message" v-else>No messages yet</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useAuth } from '../composables/useAuth.js';
import { useImageUrl } from '../composables/useImageUrl.js';
import groupDefaultIcon from "/assets/icons/group-default.png";
import userDefaultIcon from "/assets/icons/user-default.png";

const { getCurrentUsername, getCurrentUserId } = useAuth();
const { getImageUrl } = useImageUrl();

const props = defineProps({
  conversation: {
    type: Object,
    required: true,
  },
});

// Find the other participant in private conversations
const otherParticipant = computed(() => {
  if (props.conversation.isGroup) return null;
  
  return props.conversation.participants?.find(
    participant => participant.userId !== getCurrentUserId()
  ) || null;
});

// Compute the display name based on whether it's a group or private conversation
const displayName = computed(() => {
  if (props.conversation.isGroup) {
    return props.conversation.name;
  } else {
    return otherParticipant.value ? otherParticipant.value.username : props.conversation.name;
  }
});

// Compute the photo URL based on conversation type
const conversationPhotoUrl = computed(() => {
  if (props.conversation.isGroup) {
    // Group conversation - use group photo or default group icon
    if (props.conversation.photo?.path) {
      return getImageUrl(props.conversation.photo.path);
    }
    return groupDefaultIcon;
  } else {
    // Private conversation - use other participant's photo or default user icon
    if (otherParticipant.value?.photo?.path) {
      return getImageUrl(otherParticipant.value.photo.path);
    }
    return userDefaultIcon;
  }
});
</script>

<style scoped>
.conversation-preview {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 5px;
}

.conversation-preview:hover {
  background-color: #f0f0f0;
}

.conversation-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.conversation-preview.selected {
  background-color: #333;
  color: #fff;
}

.conversation-details {
  flex: 1;
}

.conversation-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
}

.last-message {
  font-size: 14px;
  color: #666;
  margin: 0;
}

.conversation-preview.selected .last-message {
  color: #ccc;
}
</style>