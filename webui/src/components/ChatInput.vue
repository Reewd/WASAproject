<template>
  <div class="chat-input ">
    <!-- Reply preview -->
    <div v-if="replyToMessage" class="reply-preview">
      <div class="reply-info">
        <strong>Replying to {{ replyToMessage.sentBy.username }}</strong>
        <span>{{ replyToMessage.text || '📷 Photo' }}</span>
      </div>
      <button @click="cancelReply" class="cancel-reply">✕</button>
    </div>

    <!-- Photo preview -->
    <div v-if="selectedPhoto" class="photo-preview">
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
        :disabled="isUploading || (!message.trim() && !selectedPhoto)"
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
import { useUser } from '../composables/useUser.js';

const { getUserId } = useUser();

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
const selectedPhoto = ref(null);
const photoInput = ref(null);
const textInput = ref(null);
const isUploading = ref(false);

// Create preview URL for selected photo
const photoPreviewUrl = computed(() => {
  return selectedPhoto.value ? URL.createObjectURL(selectedPhoto.value) : null;
});

const triggerPhotoUpload = () => {
  photoInput.value?.click();
};

const handlePhotoSelection = (event) => {
  const file = event.target.files[0];
  if (file) {
    // Validate file type
    if (!file.type.startsWith('image/')) {
      alert('Please select a valid image file');
      return;
    }
    
    // Validate file size (e.g., max 10MB)
    if (file.size > 10 * 1024 * 1024) {
      alert('File size must be less than 10MB');
      return;
    }
    
    selectedPhoto.value = file;
  }
  // Reset input value to allow selecting the same file again
  event.target.value = '';
};

const removePhoto = () => {
  selectedPhoto.value = null;
  if (photoPreviewUrl.value) {
    URL.revokeObjectURL(photoPreviewUrl.value);
  }
};

const uploadPhoto = async (photoFile) => {
  const formData = new FormData();
  formData.append('image', photoFile);

  try {
    const response = await axios.post('/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
        Authorization: getUserId(),
      },
    });
    console.log('Photo uploaded successfully:', response.data);
    return response.data; // Assuming this returns a Photo object
  } catch (error) {
    console.error('Error uploading photo:', error);
    throw error;
  }
};

const sendMessage = async () => {
  if (!message.value.trim() && !selectedPhoto.value) {
    return; // Prevent sending empty messages
  }

  isUploading.value = true;

  try {
    let photoData = null;

    // Upload photo if selected
    if (selectedPhoto.value) {
      photoData = await uploadPhoto(selectedPhoto.value);
    }

    // Prepare message request
    const messageRequest = {
      replyTo: props.replyToMessage?.messageId || null,
      text: message.value.trim() || null,
      photo: photoData || null,
    };
    console.log('Message request:', messageRequest);

    // Send message
    const response = await axios.post(
      `/conversations/${props.conversationId}/messages`,
      messageRequest,
      {
        headers: {
          'Content-Type': 'application/json',
          Authorization: getUserId(),
        },
      }
    );

    // Clear input and photo
    message.value = '';
    removePhoto();
    cancelReply();

    // Emit success event
    emits('messageSent', response.data);

    console.log('Message sent successfully:', response.data);
  } catch (error) {
    console.error('Error sending message:', error);
    alert('Failed to send message. Please try again.');
  } finally {
    isUploading.value = false;
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
}

.reply-info {
  display: flex;
  flex-direction: column;
}

.reply-info strong {
  font-size: 12px;
  color: #1976d2;
  margin-bottom: 2px;
}

.cancel-reply {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 16px;
  color: #666;
  padding: 4px;
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