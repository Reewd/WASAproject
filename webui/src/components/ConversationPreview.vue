<template>
  <div class="conversation-preview">
    <img 
      :src="conversation.photo?.path || getDefaultPhoto(conversation.isGroup)" 
      alt="Conversation Photo" 
      class="conversation-photo"
    />
    <div class="conversation-details">
      <h3 class="conversation-name">
        {{ displayName }}
      </h3>
      <p class="last-message" v-if="conversation.lastMessage">
        <strong>{{ conversation.lastMessage.sentBy.username }}:</strong> 
        {{ conversation.lastMessage.text || 'Attachment' }}
      </p>
      <p class="last-message" v-else>No messages yet</p>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useUser } from '../composables/useUser.js';

const { getUsername } = useUser();
const currentUsername = getUsername();

const props = defineProps({
  conversation: {
    type: Object,
    required: true,
  },
});

// Compute the display name based on whether it's a group or private conversation
const displayName = computed(() => {
  if (props.conversation.isGroup) {
    return props.conversation.name;
  } else {
    // For private conversations, find the other participant's name
    const otherParticipant = props.conversation.participants?.find(
      participant => participant.username !== currentUsername
    );
    return otherParticipant ? otherParticipant.username : props.conversation.name;
  }
});

// Function to return default photo paths
const getDefaultPhoto = (isGroup) => {
  return isGroup 
    ? '/assets/icons/group-default.png' 
    : '/assets/icons/user-default.png';
};
</script>

<style scoped>
.conversation-preview {
  display: flex;
  align-items: center;
  padding: 10px;
  border-radius: 5px;
}
.conversation-preview:hover {
  background-color: #f0f0f0; /* Light gray background on hover */
}

.conversation-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
}

.conversation-preview.selected {
  background-color: #333; /* Dark gray background for selected conversation */
  color: #fff; /* White text for better contrast */
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
</style>