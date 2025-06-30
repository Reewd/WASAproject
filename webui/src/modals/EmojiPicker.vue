<template>
  <div 
    v-if="isVisible" 
    class="emoji-picker-overlay"
    @click="$emit('close')"
  >
    <div 
      class="emoji-picker-container"
      :style="containerStyle"
      @click.stop
    >
      <div class="emoji-picker-header">
        <h3>Pick an emoji</h3>
        <button @click="$emit('close')" class="close-button">✕</button>
      </div>
      
      <div class="emoji-picker-content">
        <!-- Search input -->
        <div class="search-section">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search emojis..."
            class="emoji-search"
            @input="filterEmojis"
          />
        </div>

        <!-- Category tabs -->
        <div class="category-tabs">
          <button
            v-for="category in categories"
            :key="category.name"
            @click="selectedCategory = category.name"
            :class="['category-tab', { active: selectedCategory === category.name }]"
            :title="category.label"
          >
            {{ category.icon }}
          </button>
        </div>

        <!-- Emoji grid -->
        <div class="emoji-grid">
          <button
            v-for="emoji in filteredEmojis"
            :key="emoji"
            @click="selectEmoji(emoji)"
            class="emoji-button"
            :title="getEmojiName(emoji)"
          >
            {{ emoji }}
          </button>
        </div>

        <!-- No results message -->
        <div v-if="filteredEmojis.length === 0" class="no-results">
          <p>No emojis found</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  },
  position: {
    type: Object,
    default: () => ({ x: 0, y: 0 })
  }
});

const emits = defineEmits(['close', 'selectEmoji']);

const searchQuery = ref('');
const selectedCategory = ref('smileys');

// Emoji categories with their emojis
const categories = ref([
  {
    name: 'smileys',
    label: 'Smileys & Emotion',
    icon: '😊',
    emojis: [
      '😀', '😃', '😄', '😁', '😆', '😅', '😂', '🤣', '😊', '😇',
      '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘', '😗', '😙', '😚',
      '😋', '😛', '😝', '😜', '🤪', '🤨', '🧐', '🤓', '😎', '🤩',
      '🥳', '😏', '😒', '😞', '😔', '😟', '😕', '🙁', '☹️', '😣',
      '😖', '😫', '😩', '🥺', '😢', '😭', '😤', '😠', '😡', '🤬',
      '🤯', '😳', '🥵', '🥶', '😱', '😨', '😰', '😥', '😓', '🤗',
      '🤔', '🤭', '🤫', '🤥', '😶', '😐', '😑', '😬', '🙄', '😯',
      '😦', '😧', '😮', '😲', '🥱', '😴', '🤤', '😪', '😵', '🤐'
    ]
  },
  {
    name: 'gestures',
    label: 'People & Body',
    icon: '👍',
    emojis: [
      '👍', '👎', '👌', '✌️', '🤞', '🤟', '🤘', '🤙', '👈', '👉',
      '👆', '🖕', '👇', '☝️', '👋', '🤚', '🖐️', '✋', '🖖', '👏',
      '🙌', '👐', '🤲', '🤝', '🙏', '✍️', '💅', '🤳', '💪', '🦾',
      '🦵', '🦿', '🦶', '👂', '🦻', '👃', '🧠', '🫀', '🫁', '🦷',
      '🦴', '👀', '👁️', '👅', '👄', '💋', '🩸'
    ]
  },
  {
    name: 'animals',
    label: 'Animals & Nature',
    icon: '🐶',
    emojis: [
      '🐶', '🐱', '🐭', '🐹', '🐰', '🦊', '🐻', '🐼', '🐨', '🐯',
      '🦁', '🐮', '🐷', '🐽', '🐸', '🐵', '🙈', '🙉', '🙊', '🐒',
      '🐔', '🐧', '🐦', '🐤', '🐣', '🐥', '🦆', '🦅', '🦉', '🦇',
      '🐺', '🐗', '🐴', '🦄', '🐝', '🐛', '🦋', '🐌', '🐞', '🐜',
      '🦟', '🦗', '🕷️', '🕸️', '🦂', '🐢', '🐍', '🦎', '🦖', '🦕',
      '🐙', '🦑', '🦐', '🦞', '🦀', '🐡', '🐠', '🐟', '🐬', '🐳',
      '🐋', '🦈', '🐊', '🐅', '🐆', '🦓', '🦍', '🦧', '🐘', '🦛',
      '🦏', '🐪', '🐫', '🦒', '🦘', '🐃', '🐂', '🐄', '🐎', '🐖'
    ]
  },
  {
    name: 'food',
    label: 'Food & Drink',
    icon: '🍕',
    emojis: [
      '🍎', '🍐', '🍊', '🍋', '🍌', '🍉', '🍇', '🍓', '🫐', '🍈',
      '🍒', '🍑', '🥭', '🍍', '🥥', '🥝', '🍅', '🍆', '🥑', '🥦',
      '🥬', '🥒', '🌶️', '🫑', '🌽', '🥕', '🫒', '🧄', '🧅', '🥔',
      '🍠', '🥐', '🥯', '🍞', '🥖', '🥨', '🧀', '🥚', '🍳', '🧈',
      '🥞', '🧇', '🥓', '🥩', '🍗', '🍖', '🦴', '🌭', '🍔', '🍟',
      '🍕', '🫓', '🥪', '🥙', '🧆', '🌮', '🌯', '🫔', '🥗', '🥘',
      '🫕', '🍝', '🍜', '🍲', '🍛', '🍣', '🍱', '🥟', '🦪', '🍤',
      '🍙', '🍚', '🍘', '🍥', '🥠', '🥮', '🍢', '🍡', '🍧', '🍨'
    ]
  },
  {
    name: 'activities',
    label: 'Activities',
    icon: '⚽',
    emojis: [
      '⚽', '🏀', '🏈', '⚾', '🥎', '🎾', '🏐', '🏉', '🥏', '🎱',
      '🪀', '🏓', '🏸', '🏒', '🏑', '🥍', '🏏', '🪃', '🥅', '⛳',
      '🪁', '🏹', '🎣', '🤿', '🥊', '🥋', '🎽', '🛹', '🛷', '⛸️',
      '🥌', '🎿', '⛷️', '🏂', '🪂', '🏋️', '🤼', '🤸', '⛹️', '🤺',
      '🏇', '🧘', '🏄', '🏊', '🤽', '🚣', '🧗', '🚵', '🚴', '🏆',
      '🥇', '🥈', '🥉', '🏅', '🎖️', '🏵️', '🎗️', '🎫', '🎟️', '🎪',
      '🤹', '🎭', '🩰', '🎨', '🎬', '🎤', '🎧', '🎼', '🎵', '🎶'
    ]
  },
  {
    name: 'objects',
    label: 'Objects',
    icon: '🎁',
    emojis: [
      '🎁', '🎀', '🎊', '🎉', '🎈', '🎂', '🍰', '🧁', '🍭', '🍬',
      '🍫', '🍩', '🍪', '🎄', '🎃', '🎆', '🎇', '🧨', '✨', '🎖️',
      '🏆', '⚽', '🏀', '🏈', '⚾', '🥎', '🎾', '🏐', '🏉', '🥏',
      '📱', '💻', '🖥️', '🖨️', '⌨️', '🖱️', '🖲️', '💽', '💾', '💿',
      '📀', '🧮', '🎥', '🎞️', '📽️', '🎬', '📺', '📻', '🎙️', '🎚️',
      '🎛️', '🧭', '⏱️', '⏲️', '⏰', '🕰️', '⌛', '⏳', '📡', '🔋'
    ]
  }
]);

// Get all emojis for the selected category or search results
const filteredEmojis = computed(() => {
  const currentCategory = categories.value.find(cat => cat.name === selectedCategory.value);
  let emojis = currentCategory ? currentCategory.emojis : [];
  
  if (searchQuery.value.trim()) {
    // Simple search - in a real app you might want more sophisticated search
    const query = searchQuery.value.toLowerCase();
    emojis = categories.value
      .flatMap(cat => cat.emojis)
      .filter(emoji => {
        const name = getEmojiName(emoji).toLowerCase();
        return name.includes(query);
      });
  }
  
  return emojis;
});

// Container positioning styles
const containerStyle = computed(() => {
  const { x, y } = props.position;
  const viewportWidth = window.innerWidth || document.documentElement.clientWidth;
  const viewportHeight = window.innerHeight || document.documentElement.clientHeight;
  
  // Adjust position to ensure modal stays within viewport
  const adjustedX = Math.min(Math.max(x - 140, 20), viewportWidth - 300);
  const adjustedY = Math.min(Math.max(y - 340, 20), viewportHeight - 340);
  
  return {
    position: 'absolute',
    left: `${adjustedX}px`,
    top: `${adjustedY}px`,
    zIndex: 10000
  };
});

// Simple emoji name mapping (you could expand this)
const emojiNames = {
  '😀': 'grinning face',
  '😃': 'grinning face with big eyes',
  '😄': 'grinning face with smiling eyes',
  '😁': 'beaming face with smiling eyes',
  '😆': 'grinning squinting face',
  '😅': 'grinning face with sweat',
  '😂': 'face with tears of joy',
  '🤣': 'rolling on the floor laughing',
  '😊': 'smiling face with smiling eyes',
  '😇': 'smiling face with halo',
  '🙂': 'slightly smiling face',
  '🙃': 'upside down face',
  '😉': 'winking face',
  '😌': 'relieved face',
  '😍': 'smiling face with heart eyes',
  '🥰': 'smiling face with hearts',
  '😘': 'face blowing a kiss',
  '😗': 'kissing face',
  '😙': 'kissing face with smiling eyes',
  '😚': 'kissing face with closed eyes',
  '😋': 'face savoring food',
  '😛': 'face with tongue',
  '😝': 'squinting face with tongue',
  '😜': 'winking face with tongue',
  '🤪': 'zany face',
  '🤨': 'face with raised eyebrow',
  '🧐': 'face with monocle',
  '🤓': 'nerd face',
  '😎': 'smiling face with sunglasses',
  '🤩': 'star struck',
  '🥳': 'partying face',
  '😏': 'smirking face',
  '😒': 'unamused face',
  '😞': 'disappointed face',
  '😔': 'pensive face',
  '😟': 'worried face',
  '😕': 'confused face',
  '🙁': 'slightly frowning face',
  '☹️': 'frowning face',
  '😣': 'persevering face',
  '😖': 'confounded face',
  '😫': 'tired face',
  '😩': 'weary face',
  '🥺': 'pleading face',
  '😢': 'crying face',
  '😭': 'loudly crying face',
  '😤': 'face with steam from nose',
  '😠': 'angry face',
  '😡': 'pouting face',
  '🤬': 'face with symbols on mouth',
  '🤯': 'exploding head',
  '😳': 'flushed face',
  '🥵': 'hot face',
  '🥶': 'cold face',
  '😱': 'face screaming in fear',
  '😨': 'fearful face',
  '😰': 'anxious face with sweat',
  '😥': 'sad but relieved face',
  '😓': 'downcast face with sweat',
  '🤗': 'hugging face',
  '🤔': 'thinking face',
  '🤭': 'face with hand over mouth',
  '🤫': 'shushing face',
  '🤥': 'lying face',
  '😶': 'face without mouth',
  '😐': 'neutral face',
  '😑': 'expressionless face',
  '😬': 'grimacing face',
  '🙄': 'face with rolling eyes',
  '😯': 'hushed face',
  '😦': 'frowning face with open mouth',
  '😧': 'anguished face',
  '😮': 'face with open mouth',
  '😲': 'astonished face',
  '🥱': 'yawning face',
  '😴': 'sleeping face',
  '🤤': 'drooling face',
  '😪': 'sleepy face',
  '😵': 'dizzy face',
  '🤐': 'zipper mouth face',
  '👍': 'thumbs up',
  '👎': 'thumbs down',
  '👌': 'ok hand',
  '✌️': 'victory hand',
  '🤞': 'crossed fingers',
  '🤟': 'love you gesture',
  '🤘': 'sign of the horns',
  '🤙': 'call me hand',
  '👈': 'backhand index pointing left',
  '👉': 'backhand index pointing right',
  '👆': 'backhand index pointing up',
  '🖕': 'middle finger',
  '👇': 'backhand index pointing down',
  '☝️': 'index pointing up',
  '👋': 'waving hand',
  '🤚': 'raised back of hand',
  '🖐️': 'hand with fingers splayed',
  '✋': 'raised hand',
  '🖖': 'vulcan salute',
  '👏': 'clapping hands',
  '🙌': 'raising hands',
  '👐': 'open hands',
  '🤲': 'palms up together',
  '🤝': 'handshake',
  '🙏': 'folded hands'
};

const getEmojiName = (emoji) => {
  return emojiNames[emoji] || emoji;
};

const filterEmojis = () => {
  // This reactive computed will handle the filtering
};

const selectEmoji = (emoji) => {
  emits('selectEmoji', emoji);
  emits('close');
};

onMounted(() => {
  // Set initial category
  selectedCategory.value = 'smileys';
});
</script>

<style scoped>
.emoji-picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  height: 100%;
  background: transparent;
  z-index: 9999;
}

.emoji-picker-container {
  background: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 280px;
  max-height: 320px;
  display: flex;
  flex-direction: column;
  border: 1px solid #e1e1e1;
}

.emoji-picker-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  border-bottom: 1px solid #e1e1e1;
  background: #f8f9fa;
  border-radius: 8px 8px 0 0;
}

.emoji-picker-header h3 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #333;
}

.close-button {
  background: none;
  border: none;
  font-size: 16px;
  cursor: pointer;
  color: #666;
  padding: 4px;
  border-radius: 4px;
  transition: background-color 0.2s;
}

.close-button:hover {
  background-color: #e9ecef;
}

.emoji-picker-content {
  display: flex;
  flex-direction: column;
  flex: 1;
  overflow: hidden;
}

.search-section {
  padding: 8px 12px;
  border-bottom: 1px solid #e1e1e1;
}

.emoji-search {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 13px;
  outline: none;
  transition: border-color 0.2s;
}

.emoji-search:focus {
  border-color: #007bff;
}

.category-tabs {
  display: flex;
  padding: 6px 8px;
  border-bottom: 1px solid #e1e1e1;
  background: #f8f9fa;
  overflow-x: auto;
}

.category-tab {
  background: none;
  border: none;
  padding: 6px 8px;
  cursor: pointer;
  font-size: 14px;
  border-radius: 4px;
  margin-right: 2px;
  transition: background-color 0.2s;
  flex-shrink: 0;
}

.category-tab:hover {
  background-color: #e9ecef;
}

.category-tab.active {
  background-color: #007bff;
  color: white;
}

.emoji-grid {
  padding: 8px;
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 3px;
  overflow-y: auto;
  flex: 1;
  max-height: 180px;
}

.emoji-button {
  background: none;
  border: none;
  padding: 8px;
  cursor: pointer;
  font-size: 18px;
  border-radius: 6px;
  transition: background-color 0.2s, transform 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 36px;
}

.emoji-button:hover {
  background-color: #f0f0f0;
  transform: scale(1.1);
}

.emoji-button:active {
  transform: scale(0.95);
}

.no-results {
  padding: 20px;
  text-align: center;
  color: #666;
}

.no-results p {
  margin: 0;
  font-size: 14px;
}


</style>