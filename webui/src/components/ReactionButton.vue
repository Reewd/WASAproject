<template>
  <div 
    class="reaction-button-container"
    :class="{ 'own-message': isOwnMessage }"
    @mouseenter="showButton = true"
    @mouseleave="showButton = false"
  >
    <button
      @click="handleButtonClick"
      class="reaction-button"
      :class="{ 'visible': showButton }"
      :title="'Add reaction'"
    >
      😊
    </button>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  isOwnMessage: {
    type: Boolean,
    default: false,
  },
  conversationId: {
    type: Number,
    required: true,
  },
});

const emits = defineEmits(['openEmojiPicker']);

const showButton = ref(false);

const handleButtonClick = (event) => {
  event.stopPropagation();
  const rect = event.target.getBoundingClientRect();
  const position = {
    x: rect.left + rect.width / 2,
    y: rect.top
  };
  
  emits('openEmojiPicker', position);
};
</script>

<style scoped>
.reaction-button-container {
  position: absolute;
  top: 50%;
  right: -40px;
  transform: translateY(-50%);
  z-index: 10;
}

.reaction-button-container.own-message {
  right: auto;
  left: -40px;
}

.reaction-button {
  background: white;
  border: 1px solid #e1e1e1;
  border-radius: 50%;
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  cursor: pointer;
  opacity: 0;
  transition: all 0.2s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.reaction-button.visible {
  opacity: 1;
}

.reaction-button:hover {
  background-color: #f8f9fa;
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.reaction-button:active {
  transform: scale(0.95);
}

.message-container:hover .reaction-button {
  opacity: 1;
}
</style>