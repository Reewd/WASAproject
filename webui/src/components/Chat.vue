<template>
  <div class="chat-container" v-if="conversation">
    <h2>{{ conversation.name }}</h2>
    <div class="chat-messages">
    </div>
    <ChatInput />
  </div>
  <div v-else class="no-conversation">
    <p>Please select a conversation to start chatting.</p>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import ChatInput from './ChatInput.vue'; // Import the ChatInput component

// Accept the selected conversation as a prop
defineProps({
  conversation: {
    type: Object,
    required: false, // It can be null if no conversation is selected
  },
});
</script>

<style scoped>
.chat-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 20px;
  border-left: 1px solid #ccc;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  margin-bottom: 10px;
  border: 1px solid #ddd;
  padding: 10px;
}

.no-conversation {
  flex: 1;
  display: flex;
  justify-content: center;
  align-items: center;
  color: #888;
}
</style>