<template>
  <div 
    v-if="reactions.length > 0" 
    class="reaction-bubble"
    @click="removeReaction"
  >
    <span 
      v-for="reaction in reactions" 
      :key="reaction.reactionId"
      class="reaction-item"
    >
      {{ reaction.content }} <span class="reaction-username">{{ reaction.sentBy.username }}</span>
    </span>
  </div>
</template>

<script setup>
const props = defineProps({
  reactions: {
    type: Array,
    default: () => [],
  },
});

const emits = defineEmits(['removeReaction']);

// Show reaction details when clicked
const removeReaction = () => {
  emits('removeReaction', {
    reactions: props.reactions
  });
};
</script>

<style scoped>
.reaction-bubble {
  display: inline-flex;
  align-items: center;
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid #e1e1e1;
  border-radius: 16px;
  padding: 4px 8px;
  margin-top: 4px;
  cursor: pointer;
  transition: all 0.2s ease;
  backdrop-filter: blur(10px);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  max-width: 200px;
  overflow: hidden;
}

.reaction-bubble:hover {
  background-color: rgba(255, 255, 255, 0.95);
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.reaction-item {
  font-size: 14px;
  margin-right: 2px;
  line-height: 1;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}

.reaction-item:last-child {
  margin-right: 0;
}

.reaction-username {
  font-size: 12px;
  color: #666;
  margin-left: 4px;
  font-weight: 500;
}
</style>