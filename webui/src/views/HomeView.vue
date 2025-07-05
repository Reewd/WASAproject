<template>
    <div class="home-view d-flex">
        <Sidebar />
        <ConversationBar 
            ref="conversationBarRef"
            @selectConversation="selectedConversation = $event" 
        />
        <Chat
            :conversationPreview="selectedConversation"
            @groupUpdated="handleGroupUpdated"
            @leftGroup="handleLeftGroup"
        />
    </div>
</template>

<script setup>
import Sidebar from "@/components/Sidebar.vue";
import ConversationBar from "@/components/ConversationBar.vue";
import Chat from "@/components/Chat.vue";

import { ref } from "vue";

const selectedConversation = ref(null); // State to store the selected conversation
const conversationBarRef = ref(null); // Reference to the ConversationBar component

// Handler for groupUpdated event from Chat
const handleGroupUpdated = () => {
  console.log('Group updated, refreshing conversations list');
  // Call the fetchConversations method on the ConversationBar component
  conversationBarRef.value?.fetchConversations();
};

const handleLeftGroup = () => {
  console.log('User left group, deselecting conversation');
  selectedConversation.value = null; // Deselect the current conversation
  // Call the fetchConversations method to refresh the list
  conversationBarRef.value?.fetchConversations();
};
</script>

<style>
.home-view {
    display: flex;
    height: 100vh;
}
</style>