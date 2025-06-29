<template>
  <div class="home-view d-flex">
    <Sidebar />
    <ConversationBar @selectConversation="selectedConversation = $event" />
    <Chat :conversation="selectedConversation" />
  </div>
</template>

<script setup>
import Sidebar from '@/components/Sidebar.vue';
import ConversationBar from '@/components/ConversationBar.vue';
import Chat from '@/components/Chat.vue';

import { ref } from 'vue';

const selectedConversation = ref(null); // State to store the selected conversation
</script>

<style>
.home-view {
  display: flex;
  height: 100vh;
}
</style>