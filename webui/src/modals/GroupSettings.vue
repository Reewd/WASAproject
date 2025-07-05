<template>
	<div class="modal-overlay" @click="closeModal">
		<div class="modal-content" @click.stop>
			<div class="modal-header">
				<h2>Group Settings</h2>
				<button @click="closeModal" class="close-button">✕</button>
			</div>

			<div class="modal-body">
				<!-- Group Photo Section -->
				<div class="photo-section">
					<div class="photo-container">
						<img
							:src="groupPhotoUrl"
							:alt="isGroup ? 'Group Photo' : 'Chat Photo'"
							class="group-photo"
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
					<img :src="photoPreviewUrl" alt="New group photo" />
					<button @click="removePhoto" class="remove-photo">✕</button>
				</div>

				<!-- Group Name Section -->
				<div class="name-section">
					<label for="groupName">{{
						isGroup ? "Group Name:" : "Chat Name:"
					}}</label>
					<input
						id="groupName"
						type="text"
						v-model="newGroupName"
						@input="validateGroupName"
						:placeholder="chat?.name || 'Enter group name'"
						:disabled="isUpdating"
					/>
					<div class="text-danger" v-if="groupNameError">{{ groupNameError }}</div>
				</div>

				<!-- Participants Section -->
				<div class="participants-section">
					<h3>Participants ({{ participants.length }})</h3>

					<!-- Current Participants List -->
					<div class="participants-list">
						<div
							v-for="participant in participants"
							:key="participant.userId"
							class="participant-item"
						>
							<img
								:src="getParticipantPhotoUrl(participant)"
								:alt="`${participant.username}'s photo`"
								class="participant-photo"
							/>
							<span class="participant-name">{{
								participant.username
							}}</span>
							<span
								v-if="participant.userId === getCurrentUserId"
								class="current-user-badge"
								>You</span
							>
						</div>
					</div>

					<!-- Add Participants Section -->
					<div v-if="isGroup" class="add-participants">
						<h4>Add Participants:</h4>

						<UserSelection
							:excludeUsers="participants"
							@update:selectedUsers="handleSelectedUsersUpdate"
						/>

						<!-- Show pending participants count -->
						<div
							v-if="pendingParticipants.length > 0"
							class="pending-count"
						>
							<p>
								{{ pendingParticipants.length }} new
								participant{{
									pendingParticipants.length !== 1 ? "s" : ""
								}}
								will be added
							</p>
						</div>
					</div>
				</div>

				<!-- Action Buttons -->
				<div class="modal-actions">
					<button
						v-if="isGroup"
						@click="leaveGroup"
						class="leave-button"
						:disabled="isUpdating"
					>
						{{ isUpdating ? "Leaving..." : "Leave Group" }}
					</button>

					<div class="right-actions">
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
								(!hasNameChanged &&
									!hasSelectedPhoto &&
									pendingParticipants.length === 0) ||
								(hasNameChanged && !isGroupNameValid)
							"
						>
							{{ isUpdating ? "Saving..." : "Save Changes" }}
						</button>
					</div>
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
import UserSelection from '../components/UserSelection.vue';
import groupDefaultIcon from "/assets/icons/group-default.png";
import userDefaultIcon from "/assets/icons/user-default.png";

const { getCurrentUserId } = useAuth();
const { getImageUrl } = useImageUrl();
const { useGroupNameValidation } = useValidation();
const { 
    selectedPhoto, 
    photoPreviewUrl, 
    hasSelectedPhoto,
    handlePhotoSelection, 
    removePhoto, 
    uploadSelectedPhoto,
    isUploading: isUploadingPhoto
} = usePhotoUpload();

const props = defineProps({
	chat: {
		type: Object,
		required: true,
	},
	isGroup: {
		type: Boolean,
		default: false,
	},
});

const emits = defineEmits(["close", "updated", "left"]);

// Reactive data
const photoInput = ref(null);
const isUpdating = ref(false);
const pendingParticipants = ref([]);

// Use the validation composable
const { groupName: newGroupName, groupNameError, validateGroupName, isGroupNameValid } = useGroupNameValidation();

// Computed properties
const participants = computed(() => props.chat?.participants || []);

const groupPhotoUrl = computed(() => {
	if (props.chat?.photo?.path) {
		return getImageUrl(props.chat.photo.path);
	}
	return props.isGroup
		? groupDefaultIcon
		: userDefaultIcon;
});

const hasNameChanged = computed(() => {
	return newGroupName.value !== (props.chat?.name || "");
});

// Methods
const handleSelectedUsersUpdate = (selectedUsers) => {
  pendingParticipants.value = selectedUsers;
};

const initializeForm = () => {
	newGroupName.value = props.chat?.name || "";
};

const getParticipantPhotoUrl = (participant) => {
	if (participant.photo?.path) {
		return getImageUrl(participant.photo.path);
	}
	return userDefaultIcon;
};

const triggerPhotoUpload = () => {
	photoInput.value?.click();
};

const updateGroupName = async () => {
	if (!hasNameChanged.value) return;

	// Validate before sending request
	validateGroupName();
	if (groupNameError.value) {
		throw new Error(groupNameError.value);
	}

	try {
		await axios.put(
			`/conversations/${props.chat.conversationId}/name`,
			{
				name: newGroupName.value,
			},
			{
				headers: {
					"Content-Type": "application/json",
					Authorization: getCurrentUserId(),
				},
			}
		);
	} catch (error) {
		console.error("Error updating group name:", error);
		throw error;
	}
};

const updateGroupPhoto = async (photoData) => {
	try {
		await axios.put(
			`/conversations/${props.chat.conversationId}/photo`,
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
	} catch (error) {
		console.error("Error updating group photo:", error);
		throw error;
	}
};

const addParticipantsToGroup = async () => {
  if (pendingParticipants.value.length === 0) return;

  try {
    await axios.post(`/conversations/${props.chat.conversationId}/participants`, {
      participants: pendingParticipants.value.map(user => user.username) // Extract usernames for backend
    }, {
      headers: {
        'Content-Type': 'application/json',
        Authorization: getCurrentUserId(),
      },
    });
    
    pendingParticipants.value = [];
  } catch (error) {
    console.error('Error adding participants:', error);
    throw error;
  }
};

const leaveGroup = async () => {
	isUpdating.value = true;

	try {
		await axios.delete(
			`/conversations/${props.chat.conversationId}/participants`,
			{
				headers: {
					"Content-Type": "application/json",
					Authorization: getCurrentUserId(),
				},
				data: {
					conversationId: props.chat.conversationId,
				},
			}
		);

		emits("left");
		closeModal();
	} catch (error) {
		console.error("Error leaving group:", error);
		alert("Failed to leave group. Please try again.");
	} finally {
		isUpdating.value = false;
	}
};

const saveChanges = async () => {
	isUpdating.value = true;

	try {
		// Upload new photo if selected
		if (hasSelectedPhoto.value) {
			const photoData = await uploadSelectedPhoto();
			await updateGroupPhoto(photoData);
		}

		// Update group name if changed
		await updateGroupName();

		// Add new participants if any
		await addParticipantsToGroup();

		// Emit update event
		emits("updated");

		closeModal();
	} catch (error) {
		console.error("Error saving changes:", error);
		alert("Failed to update group settings. Please try again.");
	} finally {
		isUpdating.value = false;
	}
};

const closeModal = () => {
	emits("close");
};

// Lifecycle
onMounted(() => {
	initializeForm();
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
	max-width: 600px;
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

.photo-section {
	display: flex;
	flex-direction: column;
	align-items: center;
	margin-bottom: 20px;
}

.photo-container {
	position: relative;
	margin-bottom: 10px;
}

.group-photo {
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #ddd;
}

.change-photo-button {
	position: absolute;
	bottom: 5px;
	right: 5px;
	width: 28px;
	height: 28px;
	border-radius: 50%;
	background-color: #007aff;
	color: white;
	border: none;
	cursor: pointer;
	font-size: 14px;
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
	width: 100px;
	height: 100px;
	border-radius: 50%;
	object-fit: cover;
	border: 3px solid #007aff;
}

.remove-photo {
	position: absolute;
	top: 5px;
	right: calc(50% - 55px);
	background-color: rgba(0, 0, 0, 0.7);
	color: white;
	border: none;
	border-radius: 50%;
	width: 20px;
	height: 20px;
	cursor: pointer;
	font-size: 10px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.remove-photo:hover {
	background-color: rgba(0, 0, 0, 0.9);
}

.name-section {
	margin-bottom: 25px;
}

.name-section label {
	display: block;
	margin-bottom: 8px;
	font-weight: bold;
	color: #333;
}

.name-section input {
	width: 100%;
	padding: 12px;
	border: 1px solid #ddd;
	border-radius: 8px;
	font-size: 16px;
	box-sizing: border-box;
}

.name-section input:focus {
	outline: none;
	border-color: #007aff;
	box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.participants-section {
	margin-bottom: 25px;
}

.participants-section h3 {
	margin: 0 0 15px 0;
	font-size: 18px;
	color: #333;
}

.participants-list {
	max-height: 200px;
	overflow-y: auto;
	border: 1px solid #eee;
	border-radius: 8px;
	padding: 10px;
	margin-bottom: 15px;
}

.participant-item {
	display: flex;
	align-items: center;
	padding: 8px 0;
	border-bottom: 1px solid #f0f0f0;
}

.participant-item:last-child {
	border-bottom: none;
}

.participant-photo {
	width: 32px;
	height: 32px;
	border-radius: 50%;
	margin-right: 10px;
	object-fit: cover;
}

.participant-name {
	flex: 1;
	font-size: 14px;
	color: #333;
}

.current-user-badge {
	background-color: #007aff;
	color: white;
	padding: 2px 8px;
	border-radius: 12px;
	font-size: 12px;
	font-weight: 500;
}

.add-participants label {
	display: block;
	margin-bottom: 8px;
	font-weight: bold;
	color: #333;
}

.add-input-container {
	display: flex;
	gap: 8px;
	margin-bottom: 10px;
}

.add-input-container input {
	flex: 1;
	padding: 10px;
	border: 1px solid #ddd;
	border-radius: 6px;
	font-size: 14px;
}

.add-input-container input:focus {
	outline: none;
	border-color: #007aff;
}

.add-button {
	padding: 10px 16px;
	background-color: #28a745;
	color: white;
	border: none;
	border-radius: 6px;
	cursor: pointer;
	font-size: 14px;
	font-weight: 500;
}

.add-button:hover:not(:disabled) {
	background-color: #218838;
}

.add-button:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.pending-participants {
	margin-top: 10px;
}

.pending-participants p {
	font-size: 12px;
	color: #666;
	margin: 0 0 5px 0;
}

.pending-list {
	display: flex;
	flex-wrap: wrap;
	gap: 5px;
}

.pending-participant {
	display: inline-flex;
	align-items: center;
	background-color: #e3f2fd;
	padding: 4px 8px;
	border-radius: 16px;
	font-size: 12px;
	color: #1976d2;
}

.remove-pending {
	background: none;
	border: none;
	color: #1976d2;
	cursor: pointer;
	margin-left: 4px;
	font-size: 10px;
}

.remove-pending:hover {
	color: #d32f2f;
}

.modal-actions {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 12px;
}

.right-actions {
	display: flex;
	gap: 12px;
}

.leave-button,
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

.leave-button {
	background-color: #dc3545;
	color: white;
}

.leave-button:hover:not(:disabled) {
	background-color: #c82333;
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

.leave-button:disabled,
.cancel-button:disabled,
.save-button:disabled {
	opacity: 0.6;
	cursor: not-allowed;
}

.pending-count {
  margin-top: 10px;
  font-size: 14px;
  color: #28a745;
  font-weight: 500;
}

.add-participants h4 {
  margin-top: 15px;
  margin-bottom: 10px;
  font-size: 16px;
  color: #333;
}

.text-danger {
	color: #dc3545;
	font-size: 14px;
	margin-top: 5px;
}
</style>