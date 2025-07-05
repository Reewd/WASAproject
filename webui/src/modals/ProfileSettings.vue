<template>
	<div class="modal-overlay" @click="closeModal">
		<div class="modal-content" @click.stop>
			<div class="modal-header">
				<h2>Profile Settings</h2>
				<button @click="closeModal" class="close-button">✕</button>
			</div>

			<div class="modal-body">
				<!-- Profile Picture Section -->
				<div class="profile-section">
					<div class="profile-picture-container">
						<img
							:src="profilePictureUrl"
							alt="Profile Picture"
							class="profile-picture"
						/>
						<button
							@click="triggerPhotoUpload"
							class="change-photo-button"
						>
							📷
						</button>
					</div>
					<input
						ref="photoInput"
						type="file"
						accept="image/*"
						@change="handlePhotoSelection"
						style="display: none"
					/>
				</div>

				<!-- Photo preview (if new photo selected) -->
				<div v-if="hasSelectedPhoto" class="photo-preview">
					<img :src="photoPreviewUrl" alt="New profile photo" />
					<button @click="removePhoto" class="remove-photo">✕</button>
				</div>

				<!-- Username Section -->
				<div class="username-section">
					<label for="username">Update username:</label>
					<input
						id="username"
						type="text"
						v-model="newUsername"
						@input="validateUsername"
						:placeholder="
							currentUser?.name
								? `${currentUser.name}`
								: 'Enter a new username'
						"
						:disabled="isUpdating"
					/>
					<div class="text-danger" v-if="usernameError">{{ usernameError }}</div>
				</div>

				<!-- Action Buttons -->
				<div class="modal-actions">
					<button
						@click="closeModal"
						class="cancel-button"
						:disabled="isUpdating"
					>
						Cancel
					</button>
					<button
						@click="saveChanges"
						class="save-button"
						:disabled="
							isUpdating ||
							(!hasUsernameChanged && !hasSelectedPhoto) ||
							!isUsernameValid
						"
					>
						{{ isUpdating ? "Saving..." : "Save Changes" }}
					</button>
				</div>
			</div>
		</div>
	</div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import axios from "../services/axios.js";
import { useAuth } from "../composables/useAuth.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import { useValidation } from "../composables/useValidation.js";
import { usePhotoUpload } from "../composables/usePhotoUpload.js";
import userDefaultIcon from "/assets/icons/user-default.png";

const { updateUser, getCurrentUserId, getCurrentUsername, getCurrentUserPhoto } = useAuth();
const { getImageUrl } = useImageUrl();
const { 
    selectedPhoto, 
    photoPreviewUrl, 
    hasSelectedPhoto,
    handlePhotoSelection, 
    removePhoto, 
    uploadSelectedPhoto,
    isUploading: isUploadingPhoto
} = usePhotoUpload();

const emits = defineEmits(["close", "updated"]);

const currentUser = ref({
    username: getCurrentUsername(),
    userId: getCurrentUserId(),
    photo: getCurrentUserPhoto(),
});
const photoInput = ref(null);
const isUpdating = ref(false);

// Use the validation composable
const { useUsernameValidation } = useValidation();
const { username: newUsername, usernameError, validateUsername, isUsernameValid } = 
  useUsernameValidation(currentUser.value?.username || '');
// Computed properties
const profilePictureUrl = computed(() => {
    return currentUser.value.photo?.path ? getImageUrl(currentUser.value.photo.path) : userDefaultIcon;
});

const hasUsernameChanged = computed(() => {
    return newUsername.value !== currentUser.value.username;
});

// Methods
const triggerPhotoUpload = () => {
	photoInput.value?.click();
};

const updateUsername = async () => {
    if (!hasUsernameChanged.value) return;

    // Validate before sending request
    validateUsername();
    if (usernameError.value) {
        throw new Error(usernameError.value);
    }

    try {
        await axios.put(
            "/me/username",
            {
                username: newUsername.value,
            },
            {
                headers: {
                    "Content-Type": "application/json",
                    Authorization: getCurrentUserId(),
                },
            }
        );

        // Update reactive state (this will update everywhere)
        updateUser({ username: newUsername.value });
    } catch (error) {
        console.error("Error updating username:", error);
        throw error;
    }
};


const updateProfilePhoto = async (photoData) => {
	try {
		await axios.put(
			"/me/photo",
			{
				photo: photoData,
			},
			{
				headers: {
					"Content-Type": "application/json",
					Authorization: getCurrentUserId(),
				},
			}
		);

		// Update local storage
        updateUser({ photo: photoData });
	} catch (error) {
		console.error("Error updating profile photo:", error);
		throw error;
	}
};

const saveChanges = async () => {
    isUpdating.value = true;

    try {
        // Upload new photo if selected
        if (hasSelectedPhoto.value) {
            const photoData = await uploadSelectedPhoto();
            await updateProfilePhoto(photoData);
        }

        // Update username if changed
        await updateUsername();

        // Emit update event
        emits("updated");

        closeModal();
    } catch (error) {
        console.error("Error saving changes:", error);
        
        // Check if it's a 409 conflict error for username
        if (error.response && error.response.status === 409) {
            alert("Username taken");
        } else {
            alert("Failed to update profile. Please try again.");
        }
    } finally {
        isUpdating.value = false;
    }
};

const closeModal = () => {
	emits("close");
};

// Lifecycle
onMounted(() => {
	newUsername.value = currentUser.value.username;
});
</script>

<style scoped>
.modal-overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	background-color: rgba(0, 0, 0, 0.5);
	display: flex;
	justify-content: center;
	align-items: center;
	z-index: 1000;
}

.modal-content {
	background: white;
	border-radius: 12px;
	width: 90%;
	max-width: 500px;
	max-height: 80vh;
	overflow-y: auto;
}

.modal-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 20px;
	border-bottom: 1px solid #eee;
}

.modal-header h2 {
	margin: 0;
	font-size: 24px;
	color: #333;
}

.close-button {
	background: none;
	border: none;
	font-size: 24px;
	cursor: pointer;
	color: #666;
	width: 32px;
	height: 32px;
	border-radius: 50%;
	display: flex;
	align-items: center;
	justify-content: center;
}

.close-button:hover {
	background-color: #f0f0f0;
}

.modal-body {
	padding: 20px;
}

.profile-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-bottom: 20px;
}

.profile-picture-container {
	position: relative;
	margin-bottom: 10px;
}

.profile-picture {
	width: 120px;
	height: 120px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #ddd;
}

.change-photo-button {
	position: absolute;
	bottom: 5px;
	right: 5px;
	width: 32px;
	height: 32px;
	border-radius: 50%;
	background-color: #007aff;
	color: white;
	border: none;
	cursor: pointer;
	font-size: 16px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.change-photo-button:hover {
	background-color: #0056b3;
}

.photo-preview {
	position: relative;
	display: flex;
	justify-content: center;
	margin-bottom: 20px;
}

.photo-preview img {
	width: 120px;
	height: 120px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #007aff;
}

.remove-photo {
	position: absolute;
	top: 5px;
	right: calc(50% - 65px);
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

.username-section {
	margin-bottom: 30px;
}

.username-section label {
	display: block;
	margin-bottom: 8px;
	font-weight: bold;
	color: #333;
}

.username-section input {
	width: 100%;
	padding: 12px;
	border: 1px solid #ddd;
	border-radius: 8px;
	font-size: 16px;
	box-sizing: border-box;
}

.username-section input:focus {
	outline: none;
	border-color: #007aff;
	box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.username-section input:disabled {
	background-color: #f5f5f5;
	cursor: not-allowed;
}

.text-danger {
	color: #dc3545;
	font-size: 14px;
	margin-top: 5px;
}

.modal-actions {
	display: flex;
	gap: 12px;
	justify-content: flex-end;
}

.cancel-button,
.save-button {
	padding: 12px 24px;
	border-radius: 8px;
	font-size: 16px;
	font-weight: 500;
	cursor: pointer;
	border: none;
	min-width: 100px;
}

.cancel-button {
	background-color: #f5f5f5;
	color: #333;
}

.cancel-button:hover:not(:disabled) {
	background-color: #e0e0e0;
}

.save-button {
	background-color: #007aff;
	color: white;
}

.save-button:hover:not(:disabled) {
	background-color: #0056b3;
}

.cancel-button:disabled,
.save-button:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

</style>
