<template>
  <div v-if="isVisible" class="context-menu-overlay" @click="closeMenu">
    <div 
      class="context-menu" 
      :style="menuPosition"
      @click.stop
    >
      <div class="context-menu-item" @click="handleReply">
        <span class="menu-icon">↩️</span>
        <span class="menu-text">Reply</span>
      </div>
      
      <div class="context-menu-item" @click="handleForward">
        <span class="menu-icon">➡️</span>
        <span class="menu-text">Forward</span>
      </div>
      
      <div class="context-menu-divider"></div>
      
      <div class="context-menu-item danger" @click="handleDelete">
        <span class="menu-icon">🗑️</span>
        <span class="menu-text">Delete</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false,
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 }),
  },
  message: {
    type: Object,
    default: null,
  },
});

const emits = defineEmits([
  'close',
  'reply',
  'forward',
  'delete'
]);

// Calculate menu position
const menuPosition = computed(() => ({
  left: `${props.position.x - 180}px`,
  top: `${props.position.y}px`,
}));

// Menu action handlers (no logic, just emit events)
const handleReply = () => {
  emits('reply', props.message);
  closeMenu();
};

const handleForward = () => {
  emits('forward', props.message);
  closeMenu();
};

const handleDelete = () => {
  emits('delete', props.message);
  closeMenu();
};

const closeMenu = () => {
  emits('close');
};
</script>

<style scoped>
.context-menu-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 1000;
  background-color: transparent;
}

.context-menu {
  position: absolute;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15);
  border: 1px solid #e1e1e1;
  min-width: 180px;
  overflow: hidden;
  z-index: 1001;
  animation: contextMenuAppear 0.15s ease-out;
  transform-origin: top left;
}

@keyframes contextMenuAppear {
  from {
    opacity: 0;
    transform: scale(0.65) translateY(-5px);
  }
  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.context-menu-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.15s ease;
  border-bottom: 1px solid #f0f0f0;
  user-select: none;
}

.context-menu-item:last-child {
  border-bottom: none;
}

.context-menu-item:hover {
  background-color: #f8f9fa;
}

.context-menu-item.danger:hover {
  background-color: #ffebee;
}

.context-menu-item.danger .menu-text {
  color: #d32f2f;
}

.context-menu-item.danger .menu-icon {
  filter: hue-rotate(0deg) saturate(1.2);
}

.menu-icon {
  font-size: 16px;
  margin-right: 12px;
  width: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.menu-text {
  font-size: 14px;
  font-weight: 500;
  color: #333;
  flex: 1;
}

.context-menu-divider {
  height: 1px;
  background-color: #e1e1e1;
  margin: 4px 0;
}

</style>