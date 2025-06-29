<template>
  <div class="user-preview">
    <img 
      :src="getImageUrl(user.photo?.path) || '/assets/icons/user-default.png'" 
      alt="User Photo" 
      class="user-photo"
    />
    <div class="user-details">
      <h3 class="user-name">{{ user.username }}</h3>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue';
import { useImageUrl } from '../composables/useImageUrl.js';

const { getImageUrl } = useImageUrl();
// Accept a single user object as a prop
defineProps({
  user: {
    type: Object,
    required: true,
  },
});
</script>

<style scoped>
.user-preview {
  display: flex;
  align-items: center;
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 5px;
}

.user-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 10px;
}

.user-details {
  flex: 1;
}

.user-name {
  font-size: 16px;
  font-weight: bold;
  margin: 0;
}
</style>