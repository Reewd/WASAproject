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

const selectedConversation = ref(null);
const conversationBarRef = ref(null);

const handleGroupUpdated = () => {
  console.log('Group updated, refreshing conversations list');
  conversationBarRef.value?.fetchConversations();
};

const handleLeftGroup = () => {
  console.log('User left group, deselecting conversation');
  selectedConversation.value = null;
  conversationBarRef.value?.fetchConversations();
};
</script>

<style>
.home-view {
    display: flex;
    height: 100vh;
}
</style>