<template>
  <div class="chat-container" v-if="conversationPreview">
    <ChatHeader 
      :chat="chat" 
      :conversationPreview="conversationPreview"
      @groupUpdated="handleGroupUpdated"
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
        @reactionRemoved="handleReactionUpdate"
        @openEmojiPicker="handleOpenEmojiPicker"
      />
    </div>
    
    <ChatInput 
      :conversationId="conversationPreview.conversationId" 
      :replyToMessage="replyingTo"
      @cancelReply="clearReplyState"
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
import { ref, watch, nextTick, onMounted } from 'vue';
import axios from '../services/axios.js';
import { useUser } from '../composables/useUser.js';
import ChatInput from './ChatInput.vue';
import ChatHeader from './ChatHeader.vue';
import Message from './Message.vue';
import EmojiPicker from '../modals/EmojiPicker.vue';

const { getUserId } = useUser();

const props = defineProps({
  conversationPreview: {
    type: Object,
    required: false,
  },
});

const emit = defineEmits(['groupUpdated']);

const chat = ref(null);
const replyingTo = ref(null);
const messagesContainer = ref(null);
const isLoading = ref(false);

// Emoji picker state
const showEmojiPicker = ref(false);
const emojiPickerPosition = ref({ x: 0, y: 0 });
const currentMessageId = ref(null);
const currentConversationId = ref(null);
// Handle group updates
const handleGroupUpdated = () => {
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
        Authorization: getUserId(),
      },
    });
    
    chat.value = response.data;
    console.log('Chat data fetched:', chat.value);
    
    // Scroll to bottom after messages are rendered
    await nextTick();
    scrollToBottom();
    
  } catch (error) {
    console.error('Error fetching chat:', error);
    chat.value = null;
  } finally {
    isLoading.value = false;
  }
};

// Handle reaction updates
const handleReactionUpdate = () => {
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
          Authorization: getUserId(),
        },
      }
    );
    
    console.log('Reaction added:', emoji);
    handleReactionUpdate();
    
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
watch(() => props.conversationPreview?.conversationId, (newId) => {
  if (newId) {
    fetchChat(newId);
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
  }
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

/* Responsive design */
@media (max-width: 768px) {
  .chat-header {
    padding: 12px 16px;
  }
  
  .chat-messages {
    padding: 12px;
  }
  
  .chat-photo {
    width: 32px;
    height: 32px;
  }
  
  .chat-header h2 {
    font-size: 16px;
  }
}
</style>