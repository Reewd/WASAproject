<template>
	<div class="modal-overlay" @click="closeModal">
		<div class="modal-content" @click.stop>
			<div class="modal-header">
				<h2>Forward Message</h2>
				<button @click="closeModal" class="close-button">✕</button>
			</div>

			<div class="modal-body">
				<!-- Message Preview -->
				<div class="message-preview">
					<h3>Forwarding:</h3>
					<div class="preview-message">
						<div class="preview-sender">
							<img
								:src="
									getImageUrl(message.sentBy.photo?.path) ||
									'/assets/icons/user-default.png'
								"
								alt="Sender Photo"
								class="preview-photo"
							/>
							<span class="preview-username">{{
								message.sentBy.username
							}}</span>
						</div>
						<div class="preview-content">
							<img
								v-if="message.photo"
								:src="getImageUrl(message.photo.path)"
								alt="Message Photo"
								class="preview-image"
							/>
							<p v-if="message.text" class="preview-text">
								{{ message.text }}
							</p>
							<span
								v-if="!message.text && !message.photo"
								class="no-content"
								>Empty message</span
							>
						</div>
					</div>
				</div>

				<!-- Add tab navigation -->
				<div class="tab-navigation">
					<button
						@click="activeTab = 'conversations'"
						:class="[
							'tab-button',
							{ active: activeTab === 'conversations' },
						]"
					>
						Existing conversations
					</button>
					<button
						@click="activeTab = 'users'"
						:class="[
							'tab-button',
							{ active: activeTab === 'users' },
						]"
					>
						New private conversation
					</button>
				</div>

				<!-- Conversations Tab Content -->
				<div
					v-if="activeTab === 'conversations'"
					class="conversations-section"
				>
					<h3>
						Select Conversation ({{ filteredConversations.length }})
					</h3>

					<div class="conversations-list">
						<div
							v-for="conversation in filteredConversations"
							:key="conversation.conversationId"
							@click="selectConversation(conversation)"
							:class="{
								'conversation-item': true,
								selected:
									selectedConversation?.conversationId ===
									conversation.conversationId,
								'current-conversation':
									conversation.conversationId ===
									currentConversationId,
							}"
						>
							<img
								:src="getConversationPhotoUrl(conversation)"
								:alt="getConversationName(conversation)"
								class="conversation-photo"
							/>
							<div class="conversation-details">
								<span class="conversation-name">{{
									getConversationName(conversation)
								}}</span>
								<span
									v-if="
										conversation.conversationId ===
										currentConversationId
									"
									class="current-badge"
								>
									Current conversation
								</span>
							</div>
							<span
								v-if="
									selectedConversation?.conversationId ===
									conversation.conversationId
								"
								class="selected-badge"
								>✓</span
							>
						</div>
					</div>
				</div>

				<!-- Users Tab Content -->
				<div v-if="activeTab === 'users'" class="users-section">
					<UserSelection
						:excludeUsers="getUsersWithExistingConversations()"
						:allowMultiple="false"
						@update:selectedUsers="handleSelectedUsersUpdate"
					/>
				</div>

				<!-- Action Buttons -->
				<div class="modal-actions">
					<div class="left-info">
						<span
							v-if="
								activeTab === 'conversations' &&
								selectedConversation
							"
							class="selection-info"
						>
							Forward to:
							{{ getConversationName(selectedConversation) }}
						</span>
						<span
							v-else-if="activeTab === 'users' && selectedUser"
							class="selection-info"
						>
							Forward to: {{ selectedUser.username }} (New
							conversation)
						</span>
						<span v-else class="selection-info">
							Select a recipient to forward to
						</span>
					</div>

					<div class="right-actions">
						<button @click="closeModal" class="cancel-button">
							Cancel
						</button>
						<button
							@click="forwardMessage"
							class="forward-button"
							:disabled="isForwardButtonDisabled"
						>
							Forward
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
import { useUser } from "../composables/useUser.js";
import { useImageUrl } from "../composables/useImageUrl.js";
import UserSelection from "../components/UserSelection.vue"; // Add this import


const { getUserId, getUsername } = useUser();
const { getImageUrl } = useImageUrl();

const props = defineProps({
	message: {
		type: Object,
		required: true,
	},
	currentConversationId: {
		type: Number,
		required: true,
	},
});

const emits = defineEmits(["close", "forwarded"]);

// Reactive data
const searchQuery = ref("");
const conversations = ref([]);
const selectedConversation = ref(null);
const isForwarding = ref(false);
const activeTab = ref("conversations");
const users = ref([]);
const selectedUser = ref(null);

// Computed properties
const currentUsername = computed(() => getUsername());

const filteredConversations = computed(() => {
	if (!searchQuery.value) return conversations.value;

	return conversations.value.filter((conversation) => {
		const name = getConversationName(conversation).toLowerCase();
		return name.includes(searchQuery.value.toLowerCase());
	});
});

const isForwardButtonDisabled = computed(() => {
	if (activeTab.value === "conversations") {
		return (
			!selectedConversation.value ||
			selectedConversation.value.conversationId ===
				props.currentConversationId
		);
	} else {
		return !selectedUser.value;
	}
});

const getUsersWithExistingConversations = () => {
    const usersWithConversations = [];
    
    // Get all usernames that already have conversations with current user
    conversations.value.forEach((conversation) => {
        if (!conversation.isGroup) {
            const otherParticipant = conversation.participants?.find(
                (participant) => participant.username !== currentUsername.value
            );
            if (otherParticipant) {
                usersWithConversations.push(otherParticipant);
            }
        }
    });
    
    return usersWithConversations;
};

const handleSelectedUsersUpdate = (selectedUsersList) => {
  // In single selection mode, just take the first (and only) user or null
  selectedUser.value = selectedUsersList.length > 0 ? selectedUsersList[0] : null;
};

// Methods
const fetchConversations = async () => {
	try {
		const userId = getUserId(); // Retrieve the userId using the composable

		if (!userId) {
			console.error("User ID not found");
			return;
		}

		const response = await axios.get("/conversations", {
			headers: {
				Authorization: userId,
			},
		});
		conversations.value = response.data.conversations;
		console.log("Conversations fetched:", conversations.value);
	} catch (error) {
		console.error("Error fetching conversations:", error);
	}
};

const getConversationName = (conversation) => {
	if (conversation.isGroup) {
		return conversation.name;
	} else {
		const otherParticipant = conversation.participants?.find(
			(participant) => participant.username !== currentUsername.value
		);
		return otherParticipant ? otherParticipant.username : conversation.name;
	}
};

const getConversationPhotoUrl = (conversation) => {
	if (conversation.isGroup) {
		if (conversation.photo?.path) {
			return getImageUrl(conversation.photo.path);
		}
		return "/assets/icons/group-default.png";
	} else {
		const otherParticipant = conversation.participants?.find(
			(participant) => participant.username !== currentUsername.value
		);
		if (otherParticipant?.photo?.path) {
			return getImageUrl(otherParticipant.photo.path);
		}
		return "/assets/icons/user-default.png";
	}
};

const selectConversation = (conversation) => {
	// Don't allow selecting the current conversation
	if (conversation.conversationId === props.currentConversationId) {
		return;
	}

	if (
		selectedConversation.value?.conversationId ===
		conversation.conversationId
	) {
		selectedConversation.value = null; // Deselect if clicking the same conversation
	} else {
		selectedConversation.value = conversation;
	}
};

const fetchUsers = async () => {
	try {
		const response = await axios.get("/users");
		users.value = Array.isArray(response.data.users)
			? response.data.users
			: [];
		console.log("Users fetched:", users.value);
	} catch (error) {
		console.error("Error fetching users:", error);
		users.value = [];
	} // Always fetch conversations to ensure proper filtering
	await fetchConversations();
};

// Modified forwardMessage method
const forwardMessage = async () => {
	// Check for invalid selection
	if (
		(activeTab.value === "conversations" &&
			(!selectedConversation.value ||
				selectedConversation.value.conversationId ===
					props.currentConversationId)) ||
		(activeTab.value === "users" && !selectedUser.value)
	) {
		return;
	}

	isForwarding.value = true;

	try {
		let targetConversationId;

		// If forwarding to a user without existing conversation, create one first
		if (activeTab.value === "users" && selectedUser.value) {
			// Create a new conversation with the selected user
			const createResponse = await axios.post(
				"/conversations",
				{
					participants: [selectedUser.value.username],
					isGroup: false,
				},
				{
					headers: {
						"Content-Type": "application/json",
						Authorization: getUserId(),
					},
				}
			);

			// Get the newly created conversation ID
			targetConversationId = createResponse.data.conversationId;
			console.log("New conversation created:", targetConversationId);
		} else {
			// Use existing conversation
			targetConversationId = selectedConversation.value.conversationId;
		}

		// Forward the message to the target conversation
		await axios.post(
			`/conversations/${targetConversationId}/forwarded_messages`,
			{
				forwardTo: targetConversationId,
				messageId: props.message.messageId,
			},
			{
				headers: {
					"Content-Type": "application/json",
					Authorization: getUserId(),
				},
			}
		);

		console.log("Message forwarded successfully");
		emits("forwarded", {
			message: props.message,
			forwardedTo:
				activeTab.value === "users"
					? {
							conversationId: targetConversationId,
							name: selectedUser.value.username,
					  }
					: selectedConversation.value,
		});
		closeModal();
	} catch (error) {
		console.error("Error forwarding message:", error);
		alert("Failed to forward message. Please try again.");
	} finally {
		isForwarding.value = false;
	}
};

// Modified closeModal to reset selections
const closeModal = () => {
	selectedConversation.value = null;
	selectedUser.value = null;
	activeTab.value = "conversations";
	emits("close");
};

// Lifecycle
onMounted(() => {
	fetchConversations();
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
	max-height: 85vh;
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

.message-preview {
	margin-bottom: 20px;
	padding: 15px;
	background-color: #f8f9fa;
	border-radius: 8px;
	border: 1px solid #e1e1e1;
}

.message-preview h3 {
	margin: 0 0 12px 0;
	font-size: 16px;
	color: #333;
	font-weight: 600;
}

.preview-message {
	background: white;
	border-radius: 8px;
	padding: 12px;
}

.preview-sender {
	display: flex;
	align-items: center;
	margin-bottom: 8px;
}

.preview-photo {
	width: 20px;
	height: 20px;
	border-radius: 50%;
	margin-right: 8px;
	object-fit: cover;
}

.preview-username {
	font-size: 12px;
	font-weight: bold;
	color: #666;
}

.preview-content {
	display: flex;
	flex-direction: column;
}

.preview-image {
	max-width: 100px;
	max-height: 100px;
	border-radius: 6px;
	margin-bottom: 8px;
	object-fit: cover;
}

.preview-text {
	margin: 0;
	font-size: 14px;
	color: #333;
	line-height: 1.3;
}

.no-content {
	font-size: 14px;
	color: #666;
	font-style: italic;
}

.search-section {
	margin-bottom: 20px;
}

.search-input {
	width: 100%;
	padding: 12px;
	border: 1px solid #ddd;
	border-radius: 8px;
	font-size: 16px;
	box-sizing: border-box;
}

.search-input:focus {
	outline: none;
	border-color: #007aff;
	box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.1);
}

.conversations-section h3 {
	margin: 0 0 15px 0;
	font-size: 18px;
	color: #333;
}

.conversations-list {
	max-height: 300px;
	overflow-y: auto;
	border: 1px solid #eee;
	border-radius: 8px;
	padding: 8px;
}

.conversation-item {
	display: flex;
	align-items: center;
	padding: 12px 8px;
	border-radius: 6px;
	cursor: pointer;
	transition: background-color 0.2s;
	border-bottom: 1px solid #f0f0f0;
}

.conversation-item:last-child {
	border-bottom: none;
}

.conversation-item:hover {
	background-color: #f8f9fa;
}

.conversation-item.selected {
	background-color: #e3f2fd;
	border-color: #007aff;
}

.conversation-item.current-conversation {
	opacity: 0.5;
	cursor: not-allowed;
}

.conversation-item.current-conversation:hover {
	background-color: transparent;
}

.conversation-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 12px;
	object-fit: cover;
	border: 2px solid #ddd;
}

.conversation-item.selected .conversation-photo {
	border-color: #007aff;
}

.conversation-details {
	flex: 1;
	display: flex;
	flex-direction: column;
}

.conversation-name {
	font-size: 16px;
	color: #333;
	font-weight: 500;
	margin-bottom: 2px;
}

.current-badge {
	font-size: 12px;
	color: #ff9500;
	font-weight: 500;
}

/* Unused styles removed */

.selected-badge {
	background-color: #007aff;
	color: white;
	padding: 4px 8px;
	border-radius: 50%;
	font-size: 12px;
	font-weight: bold;
	width: 20px;
	height: 20px;
	display: flex;
	align-items: center;
	justify-content: center;
}

.modal-actions {
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 12px;
	margin-top: 20px;
	padding-top: 20px;
	border-top: 1px solid #eee;
}

.left-info {
	flex: 1;
}

.selection-info {
	font-size: 14px;
	color: #666;
	font-weight: 500;
}

.right-actions {
	display: flex;
	gap: 12px;
}

.cancel-button,
.forward-button {
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

.cancel-button:hover {
	background-color: #e0e0e0;
}

.forward-button {
	background-color: #007aff;
	color: white;
}

.forward-button:hover:not(:disabled) {
	background-color: #0056b3;
}

.forward-button:disabled {
	opacity: 0.6;
	cursor: not-allowed;
	background-color: #ccc;
}

.tab-navigation {
	display: flex;
	margin-bottom: 20px;
	border-bottom: 1px solid #eee;
}

.tab-button {
	flex: 1;
	padding: 12px;
	background: none;
	border: none;
	border-bottom: 2px solid transparent;
	font-size: 16px;
	font-weight: 500;
	cursor: pointer;
	color: #666;
	transition: all 0.2s;
}

.tab-button.active {
	color: #007aff;
	border-bottom-color: #007aff;
}

.users-section h3 {
	margin: 0 0 15px 0;
	font-size: 18px;
	color: #333;
}

.users-list {
	max-height: 300px;
	overflow-y: auto;
	border: 1px solid #eee;
	border-radius: 8px;
	padding: 8px;
}

.user-item {
	display: flex;
	align-items: center;
	padding: 12px 8px;
	border-radius: 6px;
	cursor: pointer;
	transition: background-color 0.2s;
	border-bottom: 1px solid #f0f0f0;
}

.user-item:last-child {
	border-bottom: none;
}

.user-item:hover {
	background-color: #f8f9fa;
}

.user-item.selected {
	background-color: #e3f2fd;
	border-color: #007aff;
}

.user-photo {
	width: 40px;
	height: 40px;
	border-radius: 50%;
	margin-right: 12px;
	object-fit: cover;
	border: 2px solid #ddd;
}

.user-item.selected .user-photo {
	border-color: #007aff;
}

.user-details {
	flex: 1;
}

.user-name {
	font-size: 16px;
	color: #333;
	font-weight: 500;
}

.empty-list-message {
	padding: 20px;
	text-align: center;
	color: #666;
	font-style: italic;
}
</style>
