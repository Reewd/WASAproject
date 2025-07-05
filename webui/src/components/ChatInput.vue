<template>
  <div class="chat-input ">
    <!-- Reply preview -->
    <div v-if="replyToMessage" class="reply-preview">
      <div class="reply-info">
        <strong>Replying to {{ replyToMessage.sentBy.username }}</strong>
        <div class="reply-content">
          <img 
            v-if="replyToMessage.photo" 
            :src="getImageUrl(replyToMessage.photo.path)" 
            alt="Reply photo" 
            class="reply-photo"
          />
          <span>{{ replyToMessage.text || (!replyToMessage.text && replyToMessage.photo ? '📷 Photo' : '') }}</span>
        </div>
      </div>
      <button @click="cancelReply" class="cancel-reply">✕</button>
    </div>

    <!-- Photo preview -->
    <div v-if="hasSelectedPhoto" class="photo-preview">
      <img :src="photoPreviewUrl" alt="Selected photo" />
      <button @click="removePhoto" class="remove-photo">✕</button>
    </div>

    <!-- Input area -->
    <div class="input-container">
      <!-- Photo upload button -->
      <button @click="triggerPhotoUpload" class="photo-button" title="Attach photo">
        📎
      </button>
      <input
        ref="photoInput"
        type="file"
        accept="image/*"
        @change="handlePhotoSelection"
        style="display: none"
      />

      <!-- Text input -->
      <input
        ref="textInput"
        type="text"
        v-model="message"
        placeholder="Type a message..."
        @keyup.enter="sendMessage"
        :disabled="isUploading"
      />

      <!-- Send button -->
      <button 
        @click="sendMessage" 
        :disabled="isUploading || (!message.trim() && !hasSelectedPhoto)"
        class="send-button"
      >
        {{ isUploading ? 'Uploading...' : 'Send' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue';
import axios from '../services/axios.js';
import { useAuth } from '../composables/useAuth.js';
import { useImageUrl } from '../composables/useImageUrl.js';
import { usePhotoUpload } from '../composables/usePhotoUpload.js';

const { user } = useAuth();
const { getImageUrl } = useImageUrl();
const { 
    selectedPhoto, 
    photoPreviewUrl, 
    hasSelectedPhoto,
    handlePhotoSelection, 
    removePhoto, 
    uploadPhoto,
    isUploading
} = usePhotoUpload();

const props = defineProps({
  conversationId: {
    type: [String, Number],
    required: true,
  },
  replyToMessage: {
    type: Object,
    default: null,
  },
});

const emits = defineEmits(['messageSent', 'cancelReply']);

const message = ref('');
const photoInput = ref(null);
const textInput = ref(null);

const triggerPhotoUpload = () => {
  photoInput.value?.click();
};

const sendMessage = async () => {
  if (!message.value.trim() && !hasSelectedPhoto.value) {
    return;
  }

  try {
    let photoData = null;

    if (hasSelectedPhoto.value) {
      photoData = await uploadPhoto(selectedPhoto.value);
    }

    const messageRequest = {
      replyTo: props.replyToMessage?.messageId || null,
      text: message.value.trim() || null,
      photo: photoData || null,
    };

    const response = await axios.post(
      `/conversations/${props.conversationId}/messages`,
      messageRequest,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: user.value.userId,
        },
      }
    );

    message.value = '';
    removePhoto();
    cancelReply();

    emits('messageSent', response.data);
  } catch (error) {
    console.error('Error sending message:', error);
    alert('Failed to send message. Please try again.');
  } finally {
    await nextTick(); // Ensure DOM updates are applied
    textInput.value?.focus(); // Refocus the input field
  }
};

const cancelReply = () => {
  emits('cancelReply');
};
</script>

<style scoped>
.chat-input {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background-color: #f8f9fa;
  border-top: 1px solid #dee2e6;
  width: 100%;
  overflow-x: hidden;
  box-sizing: border-box;
}

.reply-preview {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #e3f2fd;
  border-left: 3px solid #2196f3;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 14px;
  width: 100%;
  overflow: hidden;
  box-sizing: border-box;
}

.reply-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  width: calc(100% - 30px);
}

.reply-info strong {
  font-size: 12px;
  color: #1976d2;
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 100%;
}

.reply-content {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
  overflow: hidden;
  width: 100%;
}

.reply-content span {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  width: 0;
  flex: 1;
  min-width: 0;
}

.reply-photo {
  width: 40px;
  height: 40px;
  border-radius: 4px;
  object-fit: cover;
  flex-shrink: 0;
}

.cancel-reply {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  color: #666;
  padding: 4px;
  flex-shrink: 0;
  width: 24px;
  height: 24px;
  margin-left: 8px;
}

.cancel-reply:hover {
  color: #333;
}

.photo-preview {
  position: relative;
  display: inline-block;
  max-width: 200px;
}

.photo-preview img {
  max-width: 100%;
  max-height: 150px;
  border-radius: 8px;
  object-fit: cover;
}

.remove-photo {
  position: absolute;
  top: 4px;
  right: 4px;
  background-color: rgba(0, 0, 0, 0.7);
  color: white;
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-photo:hover {
  background-color: rgba(0, 0, 0, 0.9);
}

.input-container {
  display: flex;
  gap: 8px;
  align-items: center;
}

.photo-button {
  background-color: #ffffff;
  border: 1px solid #ced4da;
  border-radius: 20px;
  width: 36px;
  height: 36px;
  cursor: pointer;
  font-size: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.photo-button:hover {
  background-color: #5a6268;
}

.input-container input {
  flex: 1;
  padding: 10px 14px;
  border: 1px solid #ced4da;
  border-radius: 20px;
  outline: none;
  font-size: 14px;
}

.input-container input:focus {
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.input-container input:disabled {
  background-color: #f8f9fa;
  cursor: not-allowed;
}

.send-button {
  padding: 8px 16px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background-color 0.2s;
  min-width: 60px;
}

.send-button:hover:not(:disabled) {
  background-color: #0056b3;
}

.send-button:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
}
</style>