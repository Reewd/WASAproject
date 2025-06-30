<template>
  <div 
    v-if="aggregatedReactions.length > 0" 
    class="reaction-bubble"
    @click="removeReaction"
  >
    <span 
      v-for="reaction in aggregatedReactions" 
      :key="reaction.emoji"
      class="reaction-item"
    >
      {{ reaction.emoji }}
      <span v-if="reaction.count > 1" class="emoji-count">{{ reaction.count }}</span>
    </span>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  reactions: {
    type: Array,
    default: () => [],
  },
});

const emits = defineEmits(['removeReaction']);

// Aggregate reactions by emoji and count them
const aggregatedReactions = computed(() => {
  const reactionMap = new Map();
  
  props.reactions.forEach(reaction => {
    const emoji = reaction.content;
    if (reactionMap.has(emoji)) {
      reactionMap.set(emoji, {
        emoji,
        count: reactionMap.get(emoji).count + 1,
        users: [...reactionMap.get(emoji).users, reaction.sentBy]
      });
    } else {
      reactionMap.set(emoji, {
        emoji,
        count: 1,
        users: [reaction.sentBy]
      });
    }
  });
  
  // Convert to array and sort by count (most popular first)
  return Array.from(reactionMap.values()).sort((a, b) => b.count - a.count);
});

// Show reaction details when clicked
const removeReaction = () => {
  emits('removeReaction', {
    reactions: props.reactions,
    aggregated: aggregatedReactions.value
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
  margin-right: 6px;
  line-height: 1;
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
}

.reaction-item:last-child {
  margin-right: 0;
}

.emoji-count {
  font-size: 11px;
  font-weight: 600;
  color: #666;
  margin-left: 2px;
}
</style>