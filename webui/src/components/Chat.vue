<template>
  <div class="chat-container" v-if="conversationPreview">
    <ChatHeader 
      :chat="chat" 
      :conversationPreview="conversationPreview"
      @groupUpdated="handleConversationModified"
      @leftGroup="handleLeftGroup"
    />
    
    <div class="chat-messages" ref="messagesContainer">
      <Message 
        v-for="message in chat?.messages || []" 
        :key="message.messageId"
        :message="message"
        :isGroupConversation="chat?.isGroup || false"
        :replyToMessage="getReplyToMessage(message.replyTo)"
        :conversationId="conversationPreview.conversationId"
        @reply="setReplyToMessage"
        @reactionRemoved="handleMessageModified"
        @messageDeleted="handleMessageModified"
        @messageForwarded="handleConversationModified"
        @openEmojiPicker="handleOpenEmojiPicker"
      />
    </div>
    
    <ChatInput 
      :conversationId="conversationPreview.conversationId" 
      :replyToMessage="replyingTo"
      @cancelReply="clearReplyState"
      @messageSent="fetchChat(conversationPreview.conversationId)"
    />
    
    <!-- Emoji Picker Modal -->
    <EmojiPicker
      :isVisible="showEmojiPicker"
      :position="emojiPickerPosition"
      @close="closeEmojiPicker"
      @selectEmoji="handleEmojiSelect"
    />
  </div>
  <div v-else class="no-conversation">
    <p>Please select a conversation to start chatting.</p>
  </div>
</template>

<script setup>
import { ref, watch, nextTick, onMounted, onUnmounted } from 'vue';
import axios from '../services/axios.js';
import { useAuth } from "../composables/useAuth.js";
import ChatInput from './ChatInput.vue';
import ChatHeader from './ChatHeader.vue';
import Message from './Message.vue';
import EmojiPicker from '../modals/EmojiPicker.vue';

const { getCurrentUserId } = useAuth();

const props = defineProps({
  conversationPreview: {
    type: Object,
    required: false,
  },
});

const emit = defineEmits(['groupUpdated', 'leftGroup']);

const chat = ref(null);
const replyingTo = ref(null);
const messagesContainer = ref(null);
const isLoading = ref(false);
const pollingInterval = ref(null);
const POLLING_DELAY = 2000;

// Emoji picker state
const showEmojiPicker = ref(false);
const emojiPickerPosition = ref({ x: 0, y: 0 });
const currentMessageId = ref(null);
const currentConversationId = ref(null);
// Handle group updates

const handleLeftGroup = () => {
  stopPolling();
  emit('leftGroup');
};

const handleConversationModified = () => {
  if (props.conversationPreview?.conversationId) {
    fetchChat(props.conversationPreview.conversationId);
    emit('groupUpdated');
  }
};
// Fetch chat data when conversation is selected
const fetchChat = async (conversationId) => {
  if (!conversationId) {
    chat.value = null;
    return;
  }

  isLoading.value = true;
  
  try {
    const response = await axios.get(`/conversations/${conversationId}`, {
      headers: {
        Authorization: getCurrentUserId(),
      },
    });
    
    // Only update and scroll if we have new messages or first load
    const isFirstLoad = !chat.value;
    const hasNewMessages = chat.value && (
      // Safely check if messages exists and compare lengths
      (chat.value?.messages?.length || 0) !== (response.data?.messages?.length || 0) || 
      // Safely stringify messages with fallbacks to empty arrays
      JSON.stringify(chat.value?.messages || []) !== JSON.stringify(response.data?.messages || [])
    );
    
    chat.value = response.data;
    
    // Ensure chat.value.messages exists (initialize as empty array if undefined)
    if (!chat.value.messages) {
      chat.value.messages = [];
    }
    
    // Scroll to bottom after messages are rendered if there are new messages or first load
    if (isFirstLoad || hasNewMessages) {
      await nextTick();
      scrollToBottom();
    }
    
  } catch (error) {
    console.error('Error fetching chat:', error);
    // Don't reset chat on error to maintain existing messages
  } finally {
    isLoading.value = false;
  }
};

// Start polling for new messages
const startPolling = () => {
  // Clear any existing interval first
  stopPolling();
  
  // Start new polling interval
  if (props.conversationPreview?.conversationId) {
    pollingInterval.value = setInterval(() => {
      fetchChat(props.conversationPreview.conversationId);
    }, POLLING_DELAY);
  }
};

// Stop polling
const stopPolling = () => {
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value);
    pollingInterval.value = null;
  }
};

// Handle reaction updates
const handleMessageModified = () => {
  // Refresh chat data to get updated reactions
  if (props.conversationPreview?.conversationId) {
    fetchChat(props.conversationPreview.conversationId);
  }
};

// Helper functions
const getReplyToMessage = (replyToId) => {
  if (!replyToId || !chat.value?.messages) return null;
  return chat.value.messages.find(msg => msg.messageId === replyToId) || null;
};

const setReplyToMessage = (message) => {
  replyingTo.value = message;
};

const clearReplyState = () => {
  replyingTo.value = null;
};

// Emoji picker handlers
const handleOpenEmojiPicker = (data) => {
  currentMessageId.value = data.messageId;
  currentConversationId.value = data.conversationId;
  emojiPickerPosition.value = data.position;
  showEmojiPicker.value = true;
};

const closeEmojiPicker = () => {
  showEmojiPicker.value = false;
  currentMessageId.value = null;
  currentConversationId.value = null;
};

const handleEmojiSelect = async (emoji) => {
  if (!currentMessageId.value || !currentConversationId.value) return;
  
  try {
    await axios.post(
      `/conversations/${currentConversationId.value}/messages/${currentMessageId.value}/reactions`,
      { content: emoji },
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: getCurrentUserId(),
        },
      }
    );
    
    console.log('Reaction added:', emoji);
    handleMessageModified();
    
  } catch (error) {
    console.error('Error adding reaction:', error);
    alert('Failed to add reaction. Please try again.');
  }
  
  closeEmojiPicker();
};

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
  }
};

// Watch for conversation changes
watch(() => props.conversationPreview?.conversationId, (newId, oldId) => {
  // Stop previous polling
  stopPolling();
  
  if (newId) {
    fetchChat(newId);
    // Start polling for the new conversation
    startPolling();
  } else {
    chat.value = null;
  }
  // Clear reply state when switching conversations
  replyingTo.value = null;
});

// Fetch chat on component mount if conversation is already selected
onMounted(() => {
  if (props.conversationPreview?.conversationId) {
    fetchChat(props.conversationPreview.conversationId);
    startPolling();
  }
});

// Clean up polling when component is unmounted
onUnmounted(() => {
  stopPolling();
});
</script>

<style scoped>
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  border-left: 1px solid #ccc;
}

.chat-header {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
  background-color: #f8f9fa;
}

.chat-photo {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
}

.chat-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  scroll-behavior: smooth;
  background-color: #fff;
}

.no-conversation {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #888;
  background-color: #f8f9fa;
}

/* Loading state */
.chat-messages:empty::before {
  content: 'Messages will appear here';
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  color: #666;
  font-style: italic;
}
</style>