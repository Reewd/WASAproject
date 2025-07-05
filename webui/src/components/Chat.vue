<template>
  <!-- Main container for the chat interface, displayed only if a conversation is selected -->
  <div class="chat-container" v-if="conversationPreview">
    <!-- Chat header component, displays chat details and handles group-related events -->
    <ChatHeader 
      :chat="chat" 
      :conversationPreview="conversationPreview"
      @groupUpdated="handleConversationModified"
      @leftGroup="handleLeftGroup"
    />
    
    <!-- Container for chat messages, scrollable and bound to a ref for programmatic scrolling -->
    <div class="chat-messages" ref="messagesContainer">
      <!-- Message component, rendered for each message in the chat -->
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
    
    <!-- Chat input component for sending messages and handling reply state -->
    <ChatInput 
      :conversationId="conversationPreview.conversationId" 
      :replyToMessage="replyingTo"
      @cancelReply="clearReplyState"
      @messageSent="fetchChat(conversationPreview.conversationId)"
    />
    
    <!-- Emoji picker modal, displayed when the emoji picker is opened -->
    <EmojiPicker
      :isVisible="showEmojiPicker"
      :position="emojiPickerPosition"
      @close="closeEmojiPicker"
      @selectEmoji="handleEmojiSelect"
    />
  </div>
  
  <!-- Placeholder message displayed when no conversation is selected -->
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

const { user } = useAuth();

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

const showEmojiPicker = ref(false);
const emojiPickerPosition = ref({ x: 0, y: 0 });
const currentMessageId = ref(null);
const currentConversationId = ref(null);

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

const fetchChat = async (conversationId) => {
  if (!conversationId) {
    chat.value = null;
    return;
  }

  isLoading.value = true;
  
  try {
    const response = await axios.get(`/conversations/${conversationId}`, {
      headers: {
        Authorization: user.value.userId,
      },
    });
    
    const isFirstLoad = !chat.value;
    const hasNewMessages = chat.value && (
      (chat.value?.messages?.length || 0) !== (response.data?.messages?.length || 0) || 
      JSON.stringify(chat.value?.messages || []) !== JSON.stringify(response.data?.messages || [])
    );
    
    chat.value = response.data;
    
    if (!chat.value.messages) {
      chat.value.messages = [];
    }
    
    if (isFirstLoad || hasNewMessages) {
      await nextTick();
      scrollToBottom();
    }
    
  } catch (error) {
    console.error('Error fetching chat:', error);
  } finally {
    isLoading.value = false;
  }
};

const startPolling = () => {
  stopPolling();
  
  if (props.conversationPreview?.conversationId) {
    pollingInterval.value = setInterval(() => {
      fetchChat(props.conversationPreview.conversationId);
    }, POLLING_DELAY);
  }
};

const stopPolling = () => {
  if (pollingInterval.value) {
    clearInterval(pollingInterval.value);
    pollingInterval.value = null;
  }
};

const handleMessageModified = () => {
  if (props.conversationPreview?.conversationId) {
    fetchChat(props.conversationPreview.conversationId);
  }
};

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
          Authorization: user.value.userId,
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

watch(() => props.conversationPreview?.conversationId, (newId, oldId) => {
  stopPolling();
  
  if (newId) {
    fetchChat(newId);
    startPolling();
  } else {
    chat.value = null;
  }
  replyingTo.value = null;
});

onMounted(() => {
  if (props.conversationPreview?.conversationId) {
    fetchChat(props.conversationPreview.conversationId);
    startPolling();
  }
});

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
